package pollr

import (
	"fmt"
	"log"
	"time"
	"github.com/garyburd/redigo/redis"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func NewMongoSession(mongoUri string) *mgo.Session {
	// Connect to mongo
	sess, mgoErr := mgo.Dial(mongoUri)
	if mgoErr != nil {
		log.Fatal(mgoErr)
	}
	if mgoErr = sess.Ping(); mgoErr != nil {
		log.Fatal(mgoErr)
	}
	return sess
}

func NewRedisConnPool(redisUri string) *redis.Pool {
	// Connect to redis
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisUri)
			if err != nil {
				return nil, err
			}
			/*
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			*/
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type ApplicationContext interface {
	FindPoll(string) (*Poll, error)
	CreatePoll(*Poll) error
	CreateVoter(*Voter) error
	Vote(string, string, uint) error
	GetConn() redis.Conn
	CleanUp()
}

type Context struct {
	mgoSession *mgo.Session
	redisPool  *redis.Pool
}

func NewContext(mongoUri, redisUri string) *Context {
	return &Context{
		mgoSession: NewMongoSession(mongoUri),
		redisPool:  NewRedisConnPool(redisUri),
	}
}

// Call as defer c.CleanUp() in the main function
func (c *Context) CleanUp() {
	c.mgoSession.Close()
	c.redisPool.Close()
}

// FindPoll retrieves a poll
func (c *Context) FindPoll(pollId string) (*Poll, error) {
	id := bson.ObjectIdHex(pollId)
	cursor := c.mgoSession.DB("pollr").C("polls")
	result := &Poll{}
	err := cursor.FindId(id).One(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

func (c *Context) CreatePoll(poll *Poll) error {
	cursor := c.mgoSession.DB("pollr").C("polls")
	return cursor.Insert(poll)
}

func (c *Context) CreateVoter(voter *Voter) error {
	cursor := c.mgoSession.DB("pollr").C("voters")
	return cursor.Insert(voter)
}

func (c *Context) Vote(voterId, pollId string, choiceId uint) error {
	// Find the voter and check if they already voted in this poll
	voterCursor := c.mgoSession.DB("pollr").C("voters")
	voter := Voter{}
	findVoterErr := voterCursor.FindId(bson.ObjectIdHex(voterId)).One(&voter)
	if findVoterErr != nil {
		log.Printf("Couldn't find voter %v", voterId)
		return findVoterErr
	}
	if _, ok := voter.Votes[pollId]; ok {
		log.Printf("Voter %v already voted!", voterId)
		return nil
	}
	// Update the poll with the new vote, record the users vote, and publish the vote to redis
	poll := bson.ObjectIdHex(pollId)
	cursor := c.mgoSession.DB("pollr").C("polls")
	err := cursor.UpdateId(poll, bson.M{"$inc": bson.M{fmt.Sprintf("choices.%v.votes", choiceId): 1}})
	if err != nil {
		log.Printf("Couldn't cast vote for poll %s choice %v", poll.Hex(), choiceId)
		return err
	}
	voter.Votes[pollId] = choiceId
	voterCursor.UpdateId(bson.ObjectIdHex(voterId), voter)
	go c.publish(fmt.Sprintf("polls:%s", pollId), choiceId)
	return nil
}

func (c *Context) publish(channel, value interface{}) {
	conn := c.GetConn()
	defer conn.Close()
	conn.Do("Publish", channel, value)
}

func (c *Context) GetConn() redis.Conn {
	return c.redisPool.Get()
}
