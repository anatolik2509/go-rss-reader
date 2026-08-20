package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"rss-reader-backend/internal/config"
	"rss-reader-backend/internal/domain/source"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(fmt.Errorf("main: %w", err))
	}
	pool, err := prepareDb(cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("main: %w", err))
	}
	router := configureRouter(pool)
	fmt.Println("Starting server...")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Listening on port 8080")
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	<-ctx.Done()
	fmt.Println("")
	fmt.Println("Goodbye")
}

func configureRouter(pool *pgxpool.Pool) chi.Router {
	router := config.MustConfigureRootRouter()
	sourcesRouter := source.MustConfigureSourceRouter(pool)
	router.Mount("/api/sources", sourcesRouter)
	return router
}

func prepareDb(cfg *config.Config) (*pgxpool.Pool, error) {
	migrationConnect, err := config.NewMigrationDb(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("prepareDb: open migration connect: %w", err)
	}
	err = config.MustGetNewMigrator(MigrationsFS).ApplyMigrations(migrationConnect)
	if err != nil {
		return nil, fmt.Errorf("prepareDb: start migration: %w", err)
	}
	err = migrationConnect.Close()
	if err != nil {
		return nil, fmt.Errorf("prepareDb: close migration connect: %w", err)
	}
	pool, err := config.NewPgxPool(cfg.Database)
	if err != nil {
		log.Fatal(fmt.Errorf("prepareDb: create pgx pool: %w", err))
	}
	return pool, nil
}
