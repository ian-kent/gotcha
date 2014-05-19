package Config

import (
	"errors"
	"flag"
	"github.com/ian-kent/go-log/log"
	"github.com/ian-kent/gotcha/events"
	"os"
)

type Config struct {
	Listen      string
	AssetLoader func(string) ([]byte, error)
	Cache       map[string]interface{}
	Events      *events.Emitter
}

func Create(assetLoader func(string) ([]byte, error)) *Config {
	if assetLoader == nil {
		log.Warn("No asset loader provided; content loading will fail")
		assetLoader = func(asset string) ([]byte, error) {
			log.Warn("Attempted to load asset but no asset loader provided: %s", asset)
			return nil, errors.New("Not found")
		}
	}

	config := &Config{
		Listen:      ":7050",
		AssetLoader: assetLoader,
		Cache:       make(map[string]interface{}),
		Events:      &events.Emitter{},
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
