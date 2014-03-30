package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jshrake/pollr/lib"
)

type RestPoll struct {
	Id       string          `json:"id"`
	Question string          `json:"question"`
	Choices  []*pollr.Choice `json:"choices"`
}

func NewRestPoll(poll *pollr.Poll) ([]byte, error) {
	p := &RestPoll{}
	p.Id = fmt.Sprintf("/polls/%s", poll.Id.Hex())
	p.Question = poll.Question
	p.Choices = poll.Choices
	return json.Marshal(p)
}

func JsonPollResponse(w http.ResponseWriter, poll *pollr.Poll) error {
	response, err := NewRestPoll(poll)
	if err != nil {
		return err
	}
	// Respond to the client with the poll data
	// and the location of the new poll
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	return nil
}
