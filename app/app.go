package app

import (
	"github.com/ian-kent/gotcha/config"
	"github.com/ian-kent/gotcha/router"
	"github.com/ian-kent/gotcha/events"
	"log"
	"net/http"
)

type App struct {
	Config *Config.Config
	Router *Router.Router
	Server *http.Server
	Events chan *events.Event
}

func Create(assetLoader func(string) ([]byte, error)) *App {
	config := Config.Create(assetLoader)

	app := &App{
		Config: config,
		Router: Router.Create(config),
		Events: config.Events,
	}
	return app
}

func (app *App) Start() {
	app.Server = &http.Server{
		Addr:    app.Config.Listen,
		Handler: app.Router,
	}
	log.Printf("Starting application on %s\n", app.Config.Listen)
	go func() {
		err := app.Server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error binding to %s: %s", app.Config.Listen, err)
		}
	}()
}
