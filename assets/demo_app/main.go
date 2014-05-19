package main

import (
	"fmt"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/log"
	gotcha "github.com/ian-kent/gotcha/app"
	"github.com/ian-kent/gotcha/assets/demo_app/wrappers"
	"github.com/ian-kent/gotcha/events"
	"github.com/ian-kent/gotcha/form"
	"github.com/ian-kent/gotcha/http"
	nethttp "net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {
	// Create our Gotcha application
	var app = gotcha.Create(Asset)

	// Set the logger output pattern
	log.Logger().Appender().SetLayout(layout.Pattern("[%d] [%p] %m"))
	log.Logger().SetLevel(log.Stol("TRACE"))

	// Get the router
	r := app.Router

	// Create some routes
	r.Get("/", example)
	r.Post("/", examplepost)
	r.Get("/foo", example2)
	r.Get("/bar", example3)
	r.Get("/stream", streamed)
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
			Name:  "test",
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
	app.Start()

	<-make(chan int)
}

func example(session *http.Session) {
	// Stash a value and render a template
	session.Stash["Title"] = "Welcome to Gotcha"
	session.Render("index.html")
}

type ExampleForm struct {
	Title string `minlength:1; maxlength:200`
}

func examplepost(session *http.Session) {
	m := &ExampleForm{}
	session.Stash["fh"] = form.New(session, m).Populate().Validate()
	log.Info("Got posted title: %s", m.Title)
	example(session)
}

func example2(session *http.Session) {
	// Action composition, pass the first action another action
	wrappers.Foo(session, func(session *http.Session) {
		session.Response.WriteText(session.Stash["foo"].(string))
	})
}

func example3(session *http.Session) {
	session.Redirect(&url.URL{Path: "/foo"})
}

func err(session *http.Session) {
	panic("Arrggghh")
}

func streamed(session *http.Session) {
	c := session.Response.Chunked()
	for i := 1; i <= 5; i++ {
		time.Sleep(1 * time.Second)
		b := []byte(fmt.Sprintf("Counter: %d\n", i))
		c <- b
	}
	c <- make([]byte, 0)
}
