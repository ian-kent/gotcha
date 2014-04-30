package main

import (
	gotcha "github.com/ian-kent/gotcha/app"
	"github.com/ian-kent/gotcha/events"
	"github.com/ian-kent/gotcha/http"
	"github.com/ian-kent/gotcha/router"
	"log"
	"net/url"
)

func main() {
	// Create our Gotcha application
	var app = gotcha.Create(Asset)

	// Get the router
	r := app.Router

	// Create someroutes
	r.Get("/", example)
	r.Get("/foo", example2)
	r.Get("/bar", example3)
	r.Get("/err", err)

	// Serve static content (but really use a CDN)
	r.Get("/images/(?P<file>.*)", r.Static("assets/images/{{file}}"))
	r.Get("/css/(?P<file>.*)", r.Static("assets/css/{{file}}"))

	// Start our application
	log.Println("Starting application")
	app.Start()

	for {
		select {
			case e := <- app.Events:
				switch e.Event {
					case events.AfterHandler:
						log.Println("Got AfterHandler event!")
				}	
		}
	}
}

func example(session *http.Session) {
	// Stash a value and render a template
	session.Stash["Title"] = "Welcome to Gotcha"
	session.Render("index.html")
}

// An action to wrap other actions
func foo(session *http.Session, f Router.HandlerFunc) {
	session.Stash["foo"] = "bar"
	// Call the nested action
	f(session)
}

func example2(session *http.Session) {
	// Action composition, pass the first action another action
	foo(session, func(session *http.Session) {
		session.Response.WriteText(session.Stash["foo"].(string))
	})
}

func example3(session *http.Session) {
	session.Redirect(&url.URL{Path:"/foo"})
}

func err(session *http.Session) {
	panic("Arrggghh")
}

