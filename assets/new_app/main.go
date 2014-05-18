package main

import (
	gotcha "github.com/ian-kent/gotcha/app"
	"github.com/ian-kent/gotcha/http"
)

func main() {
	// Create our Gotcha application
	var app = gotcha.Create(Asset)

	// Get the router
	r := app.Router

	// Create some routes
	r.Get("/", welcome)

	// Serve static content (but really use a CDN)
	r.Get("/images/(?P<file>.*)", r.Static("assets/images/{{file}}"))
	r.Get("/css/(?P<file>.*)", r.Static("assets/css/{{file}}"))

	// Start our application
	app.Start()

	<-make(chan int)
}

func welcome(session *http.Session) {
	// Stash a value and render a template
	session.Stash["Title"] = "Welcome to Gotcha"
	session.Render("index.html")
}
