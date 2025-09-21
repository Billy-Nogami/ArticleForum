package config

import (
	"flag"
	"fmt"
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildPostgresDSN() string {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "password")
	dbname := getEnv("POSTGRES_DB", "articleforum")
	sslmode := getEnv("POSTGRES_SSLMODE", "disable")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}
