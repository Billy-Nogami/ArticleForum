package main

import (
	"ArticleForum/internal/config"
	"ArticleForum/internal/graph"
	"ArticleForum/internal/storage"
	"ArticleForum/internal/storage/memory"
	"ArticleForum/internal/storage/postgres"
	"ArticleForum/pkg/migrations"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	var store storage.Storage
	var err error

	switch cfg.StorageType {
	case "postgres":
		if cfg.PostgresDSN == "" {

			cfg.PostgresDSN = buildPostgresDSN()
		}

		if err := waitForPostgres(cfg.PostgresDSN, 10, 2*time.Second); err != nil {
			log.Fatalf("PostgreSQL is not available: %v", err)
		}

		if err := migrations.RunMigrations(cfg.PostgresDSN); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}

		store, err = postgres.NewPostgresStorage(cfg.PostgresDSN)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		log.Println("Using PostgreSQL storage")
	default:
		store = memory.NewMemoryStorage()
		log.Println("Using in-memory storage")
	}

	resolver := graph.NewResolver(store)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

func buildPostgresDSN() string {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "12345678")
	dbname := getEnv("POSTGRES_DB", "articleforum")
	sslmode := getEnv("POSTGRES_SSLMODE", "disable")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func waitForPostgres(dsn string, maxAttempts int, waitInterval time.Duration) error {
	var db *sql.DB
	var err error

	for i := 0; i < maxAttempts; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Failed to open database (attempt %d/%d): %v", i+1, maxAttempts, err)
			time.Sleep(waitInterval)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = db.PingContext(ctx); err != nil {
			log.Printf("Database not ready (attempt %d/%d): %v", i+1, maxAttempts, err)
			db.Close()
			time.Sleep(waitInterval)
			continue
		}

		db.Close()
		log.Println("Database is ready!")
		return nil
	}

	return fmt.Errorf("database not available after %d attempts: %v", maxAttempts, err)
}
