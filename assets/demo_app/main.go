package main

import (
	gotcha "github.com/ian-kent/gotcha/app"
	"github.com/ian-kent/gotcha/events"
	"github.com/ian-kent/gotcha/http"
	"github.com/ian-kent/gotcha/assets/demo_app/wrappers"
	"github.com/ian-kent/go-log/log"
	"github.com/ian-kent/go-log/layout"
	"net/url"
	"strconv"
	nethttp "net/http"
)

func main() {
	// Create our Gotcha application
	var app = gotcha.Create(Asset)

	// Set the logger output pattern
	log.Logger().Appender().SetLayout(layout.Pattern("[%d] [%p] %m"))

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

	// Listen to some events
	app.On(events.BeforeHandler, func(session *http.Session, next func()) {
		n := 0
		c, ok := session.Request.Cookies["test"]
		if ok {
			n, _ = strconv.Atoi(c.Value)
		}
		session.Stash["test"] = n
		log.Printf("Got BeforeHandler event! n = %d", n)
		next()
	})
	app.On(events.AfterHandler, func(session *http.Session, next func()) {
		n := session.Stash["test"].(int) + 1
		session.Response.Cookies.Set(&nethttp.Cookie{
			Name: "test",
			Value: strconv.Itoa(n),
		})
		log.Println("Got AfterHandler event!")
		next()
	})
	app.On(events.AfterResponse, func(session *http.Session, next func()) {
		log.Println("Got AfterResponse event!")
		next()
	})

	// Start our application
	log.Println("Starting application")
	app.Start()

	c := make(chan int)
	<- c
}

func example(session *http.Session) {
	// Stash a value and render a template
	session.Stash["Title"] = "Welcome to Gotcha"
	session.Render("index.html")
}

func example2(session *http.Session) {
	// Action composition, pass the first action another action
	wrappers.Foo(session, func(session *http.Session) {
		session.Response.WriteText(session.Stash["foo"].(string))
	})
}

func example3(session *http.Session) {
	session.Redirect(&url.URL{Path:"/foo"})
}

func err(session *http.Session) {
	panic("Arrggghh")
}

