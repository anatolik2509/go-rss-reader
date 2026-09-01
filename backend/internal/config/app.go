package config

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	migrationsFS embed.FS
	config       *Config
	db           *pgxpool.Pool
	redis        *redis.Client
	logger       *slog.Logger
	sessions     *Sessions
	router       chi.Router
}

func MustCreateNewApp(migrationsFS embed.FS) *App {
	app := &App{}
	app.migrationsFS = migrationsFS
	var err error
	app.config, err = LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	doMigrations(app)
	app.db, err = NewPgxPool(app.config.Database)
	if err != nil {
		panic(err)
	}
	app.redis, err = NewRedisClient(app.config.Redis)
	if err != nil {
		panic(err)
	}
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(app.logger)
	app.sessions = MustConfigureSessions(app)
	app.router = MustConfigureRootRouter(app)
	return app
}

func (app *App) StartApp() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	fmt.Println("Starting server...")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.router,
	}
	defer srv.Shutdown(ctx)
	fmt.Println("Listening on port 8080")
	go func() {
		app.logger.Error(srv.ListenAndServe().Error())
	}()
	<-ctx.Done()
	fmt.Println("")
	fmt.Println("Goodbye")
}

func doMigrations(app *App) {
	migrator := MustGetNewMigrator(app.migrationsFS)
	db, err := NewMigrationDb(app.config.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	migrator.ApplyMigrations(db)
}
