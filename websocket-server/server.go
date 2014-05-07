package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/jshrake/pollr/lib"
)

var (
	hostUri  = flag.String("addr", "localhost:12345", "host address")
	dbUri    = flag.String("db", "localhost:27017", "database address")
	redisUri = flag.String("redis", "localhost:6379", "redis address")
)

func main() {
	flag.Parse()
	m := martini.Classic()
	context := pollr.NewContext(*dbUri, *redisUri)
	defer context.CleanUp()
	hubManager := NewPollHubManager(context)
	m.Map(hubManager)
	m.Get("/polls/:id", AddConnectionToPollHub)
	log.Fatal(http.ListenAndServe(*hostUri, m))
}
