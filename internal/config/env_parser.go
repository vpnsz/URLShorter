package config

import "os"

const (
	serverAddressEnvName string = "SERVER_ADDRESS"
	baseUrlEnvName       string = "BASE_URL"
)

func ParseEnv(c *Config) {
	if value, ok := os.LookupEnv(serverAddressEnvName); ok {
		c.ServerAddr = value
	}
	if value, ok := os.LookupEnv(baseUrlEnvName); ok {
		c.BaseShorterAddr = value
	}
}
