package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/jshrake/pollr/lib"
	"github.com/martini-contrib/sessions"
)

// GetPollHandler attempts to find the requested poll and responds with a json representation
func GetPollHandler(w http.ResponseWriter, r *http.Request, s sessions.Session, params martini.Params, ctx pollr.ApplicationContext) {
	// If this is the first time we've seen this user
	// create a new voter entry for them in the table and
	// send them their voter id
	if s.Get("voterId") == nil {
		voter := pollr.NewVoter()
		ctx.CreateVoter(voter)
		s.Options(sessions.Options{
			Domain:   r.Host,
			Path:     "/",
			MaxAge:   0,
			HttpOnly: false,
			Secure:   false,
		})
		s.Set("voterId", voter.Id.Hex())
	}
	id := params["id"]
	poll, err := ctx.FindPoll(id)
	if err != nil {
		log.Print(err)
		http.Error(w, fmt.Sprintf("Poll %s not found", id), 404)
		return
	}
	JsonPollResponse(w, poll)
}
