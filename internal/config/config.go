package config

type Config struct {
	ServerAddr      string
	BaseShorterAddr string
	DatabaseAddr    string
	StorageFilePath string
}

func NewDefaultConfig() *Config {
	return &Config{"localhost:8080", "http://localhost:8080", "", "./default_storage.txt"}
}
