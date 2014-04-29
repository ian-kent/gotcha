package Config

type Config struct {
	Listen string
}

func Create() *Config {
	return &Config{
		Listen: ":7050",
	}
}
