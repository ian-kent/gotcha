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

func Create(assetLoader func(string)([]byte,error)) *App {
	config := Config.Create(assetLoader)

	app := &App{
		Config: config,
		Router: Router.Create(config),
		Ch: make(chan int),
	}
	return app
}

func (app *App) Start() {
	app.Server = &http.Server{
		Addr:    app.Config.Listen,
		Handler: app.Router,
	};
	go app.Server.ListenAndServe()
	log.Printf("Listening on %s\n", app.Config.Listen)
}
