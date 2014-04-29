package Config

type Config struct {
	Listen string
}

func Create() *Config {
	return &Config{}
}
