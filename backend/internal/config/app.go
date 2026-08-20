package config

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	config *Config
	db     *pgxpool.Pool
	redis  *redis.Client
	router chi.Router
}

func NewApp() (*App, error) {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		return nil, err
	}
	db, err := NewPgxPool(config.Database)
	if err != nil {
		return nil, err
	}
	redis, err := NewRedisClient(config.Redis)
	if err != nil {
		return nil, err
	}
	router := MustConfigureRootRouter()
	app := &App{
		config: config,
		db:     db,
		redis:  redis,
		router: router,
	}
	return app, nil
}
