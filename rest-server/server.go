package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/jshrake/pollr/lib"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/sessions"
)

var configFile = flag.String("config", "dev-config.json", "config file")

type App struct {
	m   *martini.ClassicMartini
	ctx pollr.ApplicationContext
}

func NewApp(context pollr.ApplicationContext) *App {
	app := &App{
		m:   martini.Classic(),
		ctx: context,
	}
	app.m.MapTo(app.ctx, (*pollr.ApplicationContext)(nil))
	return app
}

func (app *App) Run(address string) {
	defer app.ctx.CleanUp()
	log.Fatal(http.ListenAndServe(address, app.m))
}

func (app *App) SetUpMiddleware(config *pollr.Config) {
	// Setup CORS
	app.m.Use(cors.Allow(&cors.Options{
		AllowCredentials: true,
		AllowOrigins:     []string{config.WebAddress},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
	}))

	// Session handling
	app.m.Use(sessions.Sessions("pollr_session", sessions.NewCookieStore([]byte(config.AppSecret))))
}

func (app *App) SetUpRoutes(config *pollr.Config) {
	app.m.Post("/polls", CreatePollHandler)
	app.m.Get("/polls/:id", GetPollHandler)
	app.m.Put("/polls/:id", PollVoteHandler)
	//TODO: For whatever reason, the CORS middleware doesn't write the following headers
	// on the preflight response
	app.m.Options("**", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.WebAddress)
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	})
}

func main() {
	flag.Parse()
	config := pollr.NewConfig(*configFile)
	context := pollr.NewContext(config.Database.Address, config.Redis.Address)
	app := NewApp(context)
	app.SetUpMiddleware(config)
	app.SetUpRoutes(config)
	app.Run(config.RestAddress)
}
