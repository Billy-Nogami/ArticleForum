package main

import (
	"ArticleForum/internal/config"
	"ArticleForum/internal/graph"
	"ArticleForum/internal/storage"
	"ArticleForum/internal/storage/memory"
	"ArticleForum/internal/storage/postgres"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	cfg := config.Load()

	var store storage.Storage
	var err error

	switch cfg.StorageType {
	case "postgres":
		if cfg.PostgresDSN == "" {
			log.Fatal("PostgreSQL DSN is required when using postgres storage")
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
