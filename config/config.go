package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	EthNodeURL        string
	EthRequestTimeout time.Duration
}

func New() *Config {
	ethRequestConfig := getEnv("ETH_REQUEST_TIMEOUT", "2000")
	ethRequestMs, err := strconv.ParseInt(ethRequestConfig, 10, 64)
	if err != nil {
		ethRequestMs = 2000
	}

	return &Config{
		EthNodeURL:        getEnv("ETH_NODE_URL", "https://cloudflare-eth.com"),
		EthRequestTimeout: time.Duration(ethRequestMs) * time.Millisecond,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
