package app

import (
	"github.com/ian-kent/gotcha/config"
	"github.com/ian-kent/gotcha/router"
	"github.com/ian-kent/gotcha/http"
	"github.com/ian-kent/go-log/log"
	nethttp "net/http"
)

type App struct {
	Config *Config.Config
	Router *Router.Router
	Server *nethttp.Server
}

func Create(assetLoader func(string) ([]byte, error)) *App {
	config := Config.Create(assetLoader)

	app := &App{
		Config: config,
		Router: Router.Create(config),
	}
	return app
}

func (app *App) Start() {
	app.Server = &nethttp.Server{
		Addr:    app.Config.Listen,
		Handler: app.Router,
	}
	log.Printf("Starting application on %s", app.Config.Listen)
	go func() {
		err := app.Server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error binding to %s: %s", app.Config.Listen, err)
		}
	}()
}

func (app *App) On(event int, handler func(*http.Session, func())) {
	app.Config.Events.On(event, func(s interface{}, next func()) {
		handler(s.(*http.Session), next)
	})
}
