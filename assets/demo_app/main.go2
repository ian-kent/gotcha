package main

import(
	"log"
	gotcha "github.com/ian-kent/gotcha/app"
	"github.com/ian-kent/gotcha/http"
	"github.com/ian-kent/gotcha/router"
)

func main() {
	var app = gotcha.Create(Asset)

	r := app.Router
	r.Get("/", example)
	r.Get("/foo", example2)

	log.Println("Starting application")
	app.Start()

	<-app.Ch
}

func example(session *http.Session, route *Router.Route) {
	session.Stash["Title"] = "Welcome to Gotcha"
	session.Render("index.html")
}

func foo(session *http.Session, route *Router.Route, f Router.HandlerFunc) {
	session.Stash["foo"] = "bar"
	f(session, route)
}

func example2(session *http.Session, route *Router.Route) {
	foo(session, route, func(session *http.Session, route *Router.Route){
		session.Response.WriteText(session.Stash["foo"].(string))
	});
}
