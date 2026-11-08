package config

import "os"

const (
	serverAddressEnvName   string = "SERVER_ADDRESS"
	baseURLEnvName         string = "BASE_URL"
	fileStoragePathEnvName string = "FILE_STORAGE_PATH"
)

func ParseEnv(c *Config) {
	if value, ok := os.LookupEnv(serverAddressEnvName); ok {
		c.ServerAddr = value
	}
	if value, ok := os.LookupEnv(baseURLEnvName); ok {
		c.BaseShorterAddr = value
	}
	if value, ok := os.LookupEnv(fileStoragePathEnvName); ok {
		c.StorageFilePath = value
	}
}
