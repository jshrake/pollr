package pollr

import "labix.org/v2/mgo/bson"

type Choice struct {
	Text  string `json:"text" bson:"text"`
	Votes uint64 `json:"votes" bson:"votes"`
}

func NewChoice(text string) *Choice {
	return &Choice{
		Text:  text,
		Votes: 0,
	}
}

type Poll struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Question string        `json:"question" bson:"question"`
	Choices  []*Choice     `json:"choices" bson:"choices"`
}

// NewPoll returns a new poll with
func NewPoll(question string, choices ...*Choice) *Poll {
	return &Poll{
		Id:       bson.NewObjectId(),
		Question: question,
		Choices:  choices,
	}
}

// AddChoice adds a choice to the poll
func (p *Poll) AddChoice(c *Choice) {
	p.Choices = append(p.Choices, c)
}
