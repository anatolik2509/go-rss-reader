package config

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewMigrationDb(config DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName, "disable")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error while opening pgx connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("connection to db failed: %w", err)
	}
	return db, nil
}

func NewPgxPool(config DatabaseConfig) (pool *pgxpool.Pool, err error) {
	ctx := context.Background()
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName, "disable")
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}
	pool, err = pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}
	return
}
