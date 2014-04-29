package app

import (
	"log"
	"net/http"
	"github.com/ian-kent/Go-Gotcha/router"
	"github.com/ian-kent/Go-Gotcha/config"
)

type App struct {
	Config *Config.Config
	Router *Router.Router
	Server *http.Server
	Ch chan int
}

func Create(config *Config.Config) *App {
	router := &Router.Router{}
	app := &App{
		Config: config,
		Router: router,
		Server: &http.Server{
			Addr: config.Listen,
			Handler: router,
		},
		Ch: make(chan int),
	}
	return app
}

func (app *App) Start() {
	go app.Server.ListenAndServe()
	log.Printf("Listening on %s\n", app.Config.Listen)
}
