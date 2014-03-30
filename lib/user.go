package pollr

import "labix.org/v2/mgo/bson"

type Voter struct {
	Id    bson.ObjectId   `json:"id" bson:"_id"`
	Votes map[string]uint `json:"votes bson:"votes"`
}

func NewVoter() *Voter {
	return &Voter{
		Id:    bson.NewObjectId(),
		Votes: make(map[string]uint),
	}
}
