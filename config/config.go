package Config

type Config struct {
	Listen string
	AssetLoader func(string)([]byte,error)
	Cache map[string]interface{}
}

func Create(assetLoader func(string)([]byte,error)) *Config {
	return &Config{
		Listen: ":7050",
		AssetLoader: assetLoader,
		Cache: make(map[string]interface{}),
	}
}
