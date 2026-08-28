package config

import "os"

const (
	serverAddressEnvName   string = "SERVER_ADDRESS"
	baseUrlEnvName         string = "BASE_URL"
	fileStoragePathEvnName string = "FILE_STORAGE_PATH"
)

func ParseEnv(c *Config) {
	if value, ok := os.LookupEnv(serverAddressEnvName); ok {
		c.ServerAddr = value
	}
	if value, ok := os.LookupEnv(baseUrlEnvName); ok {
		c.BaseShorterAddr = value
	}
	if value, ok := os.LookupEnv(fileStoragePathEvnName); ok {
		c.StorageFilePath = value
	}
}
