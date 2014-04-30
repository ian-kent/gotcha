package main

import (
	"log"
	"os"
	"os/exec"
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
		case "install":
			if len(args) > 0 {
				log.Fatalf("No additional arguments required for install: %s", args)
			}
			if _, err := os.Stat("Makefile"); os.IsNotExist(err) {
				log.Fatalf("Current directory doesn't appear to be a Gotcha application")
			}
			out, err := exec.Command("make").Output()
			if err != nil {
				log.Fatalf("Error installing application: %s", err)
			}
			log.Printf("Install successful: %s", out)
		default:
			log.Fatalf("Unrecognised command: %s\n", cmd)
		}
	}
}

func new(name string) {
	log.Printf("Creating application: '%s'\n", name)

	// TODO clean this up
	
	createDir(name)
	writeAsset("assets/demo_app/main.go", name+"/main.go")
	writeAsset("assets/demo_app/Makefile", name+"/Makefile")
	writeAsset("assets/demo_app/README.md", name+"/README.md")

	createDir(name + "/assets/templates")
	writeAsset("assets/demo_app/assets/templates/index.html", name+"/assets/templates/index.html")
	writeAsset("assets/demo_app/assets/templates/error.html", name+"/assets/templates/error.html")
	writeAsset("assets/demo_app/assets/templates/notfound.html", name+"/assets/templates/notfound.html")
	createDir(name + "/assets/images")
	writeAsset("assets/demo_app/assets/images/logo-ish.png", name+"/assets/images/logo-ish.png")
	createDir(name + "/assets/css")
	writeAsset("assets/demo_app/assets/css/default.css", name+"/assets/css/default.css")
}

func createDir(dir string) {
	log.Printf("Creating directory %s", dir)

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatalf("Error creating directory %s: %s", dir, err)
	}
}
func writeAsset(input string, output string) {
	log.Printf("Writing asset %s to %s", input, output)

	f, err := os.Create(filepath.FromSlash(output))
	if err != nil {
		log.Fatalf("Error creating %s: %s", output, err)
	}

	bytes, err := Asset(input)
	if err != nil {
		log.Fatalf("Error loading asset %s: %s", input, err)
	}

	_, err = f.Write(bytes)
	if err != nil {
		log.Fatalf("Error writing output %s: %s", output, err)
	}
}
