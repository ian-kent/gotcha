package main

import (
	gotcha "github.com/ian-kent/gotcha/app"
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

	// Serve static content (but really use a CDN)
	r.Get("/images/(?P<file>.*)", r.Static("assets/images/{{file}}"))
	r.Get("/css/(?P<file>.*)", r.Static("assets/css/{{file}}"))

	// Start our application
	log.Println("Starting application")
	app.Start()

	// Wait for a signal from Gotcha to exit
	<-app.Ch
}

func example(session *http.Session, route *Router.Route) {
	// Stash a value and render a template
	session.Stash["Title"] = "Welcome to Gotcha"
	session.Render("index.html")
}

// An action to wrap other actions
func foo(session *http.Session, route *Router.Route, f Router.HandlerFunc) {
	session.Stash["foo"] = "bar"
	// Call the nested action
	f(session, route)
}

func example2(session *http.Session, route *Router.Route) {
	// Action composition, pass the first action another action
	foo(session, route, func(session *http.Session, route *Router.Route) {
		session.Response.WriteText(session.Stash["foo"].(string))
	})
}

func example3(session *http.Session, route *Router.Route) {
	session.Redirect(&url.URL{Path:"/foo"})
}
