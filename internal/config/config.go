package config

type Config struct {
	ServerHost      string
	ServerPort      uint16
	BaseShorterHost string
	BaseShorterPort uint16
}

func NewDefaultConfig() *Config {
	return &Config{"localhost", 8080, "http://localhost", 8080}
}
