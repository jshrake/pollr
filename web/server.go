package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
)

var addr = flag.String("addr", "localhost:8080", "server address")

func main() {
	flag.Parse()
	m := martini.Classic()
	m.Use(gzip.All())
	m.Use(martini.Static("public/src"))
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/src/index.html")
	})
	log.Fatal(http.ListenAndServe(*addr, m))
}
