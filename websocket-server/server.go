package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/jshrake/pollr/lib"
)

var configFile = flag.String("config", "dev-config.json", "config file")

func main() {
	flag.Parse()
	config := pollr.NewConfig(*configFile)
	m := martini.Classic()
	context := pollr.NewContext(config.Database.Address, config.Redis.Address)
	defer context.CleanUp()
	hubManager := NewPollHubManager(context)
	m.Map(hubManager)
	m.Get("/polls/:id", AddConnectionToPollHub)
	log.Fatal(http.ListenAndServe(config.WsAddress, m))
}
