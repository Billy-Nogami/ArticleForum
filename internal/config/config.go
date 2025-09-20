package config

import (
	"flag"
	"os"
)

type Config struct {
	Port        string
	StorageType string
	PostgresDSN string
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.StorageType, "storage", "memory", "Storage type: memory or postgres")
	flag.StringVar(&cfg.PostgresDSN, "postgres-dsn", "", "PostgreSQL data source name")
	flag.Parse()

	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}
