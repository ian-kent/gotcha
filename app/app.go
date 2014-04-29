package app

import (
	"github.com/ian-kent/Go-Gotcha/config"
	"github.com/ian-kent/Go-Gotcha/router"
	"log"
	"net/http"
)

type App struct {
	Config *Config.Config
	Router *Router.Router
	Server *http.Server
	Ch     chan int
}

func Create(config *Config.Config) *App {
	router := &Router.Router{}
	app := &App{
		Config: config,
		Router: router,
		Server: &http.Server{
			Addr:    config.Listen,
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
