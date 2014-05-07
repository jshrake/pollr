package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/jshrake/pollr/lib"
	"github.com/martini-contrib/sessions"
)

type VotePollEvent struct {
	ChoiceId uint `json:"choiceId"`
}

/* Handle requests to vote in a poll
Expects the following json request.
{
	choiceId: 0
}
Response: json response of the created poll
*/
func PollVoteHandler(w http.ResponseWriter, r *http.Request, s sessions.Session,
	params martini.Params, ctx pollr.ApplicationContext) {
	// Find the requested poll and return a 404 if not found
	voterId := s.Get("voterId")
	if voterId == nil {
		fmt.Println("Can't find voterId!")
		http.Error(w, "Cheater", 404)
		return
	}
	poll, err := ctx.FindPoll(params["id"])
	if err != nil {
		log.Printf("Poll %s not found", params["id"])
		http.Error(w, "Poll not found", 404)
		return
	}
	// Decode the json request body to find the choiceId
	votePollEvent := &VotePollEvent{}
	if err := json.NewDecoder(r.Body).Decode(&votePollEvent); err != nil {
		fmt.Printf("VOTE: ERROR DECODING VOTEPOLLEVENT %s", r.Body)
		http.Error(w, "Bad vote data", 404)
		return
	}
	if votePollEvent.ChoiceId > uint(len(poll.Choices)) {
		fmt.Println("VOTE: ERROR CHOICEID OUT OF RANGE")
		http.Error(w, "Bad vote data", 404)
		return
	}
	ctx.Vote(voterId.(string), params["id"], votePollEvent.ChoiceId)
}
