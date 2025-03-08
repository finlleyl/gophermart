package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress string
}

func LoadConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.RunAddress, "a", ":8080", "run address")
	flag.Parse()

	if addr := os.Getenv("RUN_ADDRESS"); addr != "" {
		cfg.RunAddress = addr
	}

	return cfg
}
