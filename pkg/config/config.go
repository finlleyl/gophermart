package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress  string
	DatabaseURI string
}

func LoadConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.RunAddress, "a", ":8080", "run address")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "database uri")
	flag.Parse()

	if addr := os.Getenv("RUN_ADDRESS"); addr != "" {
		cfg.RunAddress = addr
	}

	if db := os.Getenv("DATABASE_URI"); db != "" {
		cfg.DatabaseURI = db
	}

	return cfg
}
