package main

import(
	"log"
	gotcha "github.com/ian-kent/Go-Gotcha/app"
	"github.com/ian-kent/Go-Gotcha/config"
	"github.com/ian-kent/Go-Gotcha/http"
	"github.com/ian-kent/Go-Gotcha/router"
)

func main() {
	var config = Config.Create()
	config.Listen = ":7050";

	var app = gotcha.Create(config)

	router := app.Router
	router.Get("/", example)

	log.Println("Starting application")
	app.Start()

	<-app.Ch
}

func example(w *http.Response, r *http.Request, route *Router.Route) {
	w.Write([]byte("Hello world!"))
}
