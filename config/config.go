package config

import "os"

type Config struct {
	EthNodeURL string
}

func New() *Config {
	return &Config{
		EthNodeURL: getEnv("ETH_NODE_URL", "https://cloudflare-eth.com"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
