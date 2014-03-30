package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 256
)

// conn represents a websocket connection and a channel for communicating
// with a PollHub
type conn struct {
	ws   *websocket.Conn
	send chan uint64
}

func NewConn(ws *websocket.Conn) *conn {
	return &conn{
		ws:   ws,
		send: make(chan uint64, 256),
	}
}

// ReadPump receives messages from the websocket connection
// Currently, we only use websockets in a send-only mode
// TODO(jshrake): update this fucntion to handle websocket closing
func (c *conn) ReadPump() {
	defer func() {
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		// Message
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("Closing websocket connection due to: %v", err)
			break
		}
	}
}

// WritePump receives messages from the PollHub and
// writtes the message to this websocket connection
func (c *conn) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, []byte(strconv.FormatUint(msg, 10))); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *conn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}
