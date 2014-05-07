package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
)

var addr = flag.String("addr", "localhost:8080", "server address")

func main() {
	flag.Parse()
	m := martini.Classic()
	m.Use(martini.Static("public/src"))
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/src/index.html")
	})
	log.Fatal(http.ListenAndServe(*addr, m))
}
