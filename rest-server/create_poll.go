package main

import (
	"encoding/json"
	"net/http"

	"github.com/jshrake/pollr/lib"
)

type CreatePollEvent struct {
	Question string   `json:"question"`
	Choices  []string `json:"choices"`
}

/* CreatePollHandler handles requests to create a new poll. It expects a
json request of the following format:
{

	"question": "what's your favorite color?",
	"choices": ["red", "blue", "green"]
}
It responds with the complete json poll representation.
*/
func CreatePollHandler(w http.ResponseWriter, r *http.Request, ctx pollr.ApplicationContext) {
	createPollEvent := &CreatePollEvent{}
	if err := json.NewDecoder(r.Body).Decode(&createPollEvent); err != nil {
		http.Error(w, "Bad poll data", 404)
		return
	}
	poll := pollr.NewPoll(createPollEvent.Question)
	for _, choice := range createPollEvent.Choices {
		poll.AddChoice(pollr.NewChoice(choice))
	}
	if dbErr := ctx.CreatePoll(poll); dbErr != nil {
		http.Error(w, "Database error creating poll", 404)
		return
	}
	JsonPollResponse(w, poll)
}
