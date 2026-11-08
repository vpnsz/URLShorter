package config

import "os"

const (
	serverAddressEnvName   string = "SERVER_ADDRESS"
	baseURLEnvName         string = "BASE_URL"
	fileStoragePathEnvName string = "FILE_STORAGE_PATH"
	databaseAddrEnvName    string = "DATABASE_DSN"
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
	if value, ok := os.LookupEnv(databaseAddrEnvName); ok {
		c.DatabaseAddr = value
	}
}
