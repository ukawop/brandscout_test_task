package config

type Config struct {
	ServerAddress     string
	MaxHandlerWorkers int
}

func NewConfig(maxWorkers int) *Config {
	return &Config{
		ServerAddress:     ":8080",
		MaxHandlerWorkers: maxWorkers,
	}
}
