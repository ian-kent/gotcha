package Config

import (
	"flag"
	"os"
	"github.com/ian-kent/gotcha/events"
)

type Config struct {
	Listen      string
	AssetLoader func(string) ([]byte, error)
	Cache       map[string]interface{}
	Events		chan *events.Event
}

func Create(assetLoader func(string) ([]byte, error)) *Config {
	config := &Config{
		Listen:      ":7050",
		AssetLoader: assetLoader,
		Cache:       make(map[string]interface{}),
		Events:      make(chan *events.Event),
	}

	config.env()
	config.flags()

	return config
}

func (config *Config) env() {
	if listen := os.Getenv("GOTCHA_LISTEN"); listen != "" {
		config.Listen = listen
	}
}

func (config *Config) flags() {
	var listen string
	flag.StringVar(&listen, "listen", "", "Interface to listen on, e.g. '0.0.0.0:7050' or ':7050'")
	flag.Parse()

	if listen != "" {
		config.Listen = listen
	}
}
