package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"github.com/jshrake/pollr/lib"
)

type PollHubManager struct {
	sync.Mutex
	hubs map[string]*PollHub
	ctx  pollr.ApplicationContext
}

func NewPollHubManager(context pollr.ApplicationContext) *PollHubManager {
	return &PollHubManager{
		hubs: make(map[string]*PollHub),
		ctx:  context,
	}
}

func (ph *PollHubManager) GetHub(id string) *PollHub {
	ph.Lock()
	defer ph.Unlock()
	hub, ok := ph.hubs[id]
	if !ok {
		hub = NewPollHub(id, &redis.PubSubConn{ph.ctx.GetConn()})
		ph.hubs[id] = hub
	}
	return hub
}

func (ph *PollHubManager) RegisterNewConnection(pollId string, wsConn *websocket.Conn) {
	hub := ph.GetHub(pollId)
	conn := NewConn(wsConn)
	hub.Register(conn)
	defer hub.Unregister(conn)
	go conn.WritePump()
	conn.ReadPump()
}

// AddConnectionToPollHub handles upgrading requests to the websocket protocol
// and listening for updates to polls
func AddConnectionToPollHub(w http.ResponseWriter, r *http.Request, params martini.Params, hm *PollHubManager) {
	// Ensure this was a GET request and attempt to upgrade
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	wsConn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	// Registers the connection with the correct poll hub,
	// and starts the blocking read loop from
	// redis to the websocket
	hm.RegisterNewConnection(params["id"], wsConn)
}

// Represents all the websocket connections
// listening for updates for a particular poll
type PollHub struct {
	pollId      string
	psConn      *redis.PubSubConn
	vote        chan uint
	connections map[*conn]bool
	register    chan *conn
	unregister  chan *conn
	stop        chan bool
}

func NewPollHub(id string, psConn *redis.PubSubConn) *PollHub {
	psConn.Subscribe(fmt.Sprintf("polls:%s", id))
	pollHub := &PollHub{
		pollId:      id,
		psConn:      psConn,
		vote:        make(chan uint),
		connections: make(map[*conn]bool),
		register:    make(chan *conn),
		unregister:  make(chan *conn),
	}
	go pollHub.run()
	return pollHub
}

// Register a new websocket connection
// This method blocks
func (h *PollHub) Register(c *conn) {
	h.register <- c
}

// Unregister a websocket connection
// This method blocks
func (h *PollHub) Unregister(c *conn) {
	h.unregister <- c
}

// run handles registering and unregistering websocket connections,
// as well as updating all web socket connections when a new vote is
// published to redis
func (h *PollHub) run() {
	defer h.psConn.Close()
	receiveChan := receive(h.psConn)
	for {
		select {
		// Listen for any published messages
		// and broadcast them to all websocket connections
		case choiceId := <-receiveChan:
			for c := range h.connections {
				select {
				case c.send <- choiceId:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case <-h.stop:
			for c := range h.connections {
				close(c.send)
				delete(h.connections, c)
			}
			return
		}
	}
}

// receive returns a read-only channel that contains votes
// published from the passed in redis pubsub connection
func receive(pubSubConn *redis.PubSubConn) <-chan uint64 {
	c := make(chan uint64)
	go func() {
		defer close(c)
		for {
			switch n := pubSubConn.Receive().(type) {
			case redis.Message:
				choiceId, err := strconv.ParseUint(string(n.Data), 0, 64)
				if err == nil {
					c <- choiceId
				} else {
					log.Printf("Error converting %s to uint", n.Data)
				}
			case redis.PMessage:
				log.Printf("PMessage: %s %s %s\n", n.Pattern, n.Channel, n.Data)
			case redis.Subscription:
				log.Printf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
				if n.Count == 0 {
					return
				}
			case error:
				log.Printf("PollHub::receive - error receiving redis message: %v\n", n)
				return
			}
		}
	}()
	return c
}
