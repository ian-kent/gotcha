package main

import (
	"log"
	"os"
	"path/filepath"
)

var version = "0.01"

func main() {
	log.Printf("gotcha %s!\n", version)

	if len(os.Args) >= 2 {
		cmd := os.Args[1]
		args := os.Args[2:]
		switch cmd {
		case "new":
			if len(args) == 0 {
				log.Fatalf("Missing application name, e.g. 'gotcha new MyApp'")
			}
			new(args[0])
		default:
			log.Fatalf("Unrecognised command: %s\n", cmd)
		}
	}
}

func new(name string) {
	log.Printf("Creating application: '%s'\n", name)

	err := os.Mkdir(name, 0777)
	if err != nil {
		log.Fatalf("Error creating application directory: %s", err)
	}

	f, err := os.Create(filepath.FromSlash(name + "/main.go"))
	if err != nil {
		log.Fatalf("Error creating main.go: %s", err)
	}

	_, err = f.WriteString(`package main

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
	w.WriteText("Hello world!")
}
`)

}
