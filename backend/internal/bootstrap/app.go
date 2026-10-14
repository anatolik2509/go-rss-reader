package bootstrap

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
	MigrationsFS embed.FS
	Config       *Config
	Db           *pgxpool.Pool
	Redis        *redis.Client
	Sessions     *Sessions
	Router       chi.Router
}

func MustCreateNewApp(migrationsFS embed.FS) *App {
	app := &App{}
	app.MigrationsFS = migrationsFS
	var err error
	app.Config, err = LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	doMigrations(app)
	app.Db, err = NewPgxPool(app.Config.Database)
	if err != nil {
		panic(err)
	}
	app.Redis, err = NewRedisClient(app.Config.Redis)
	if err != nil {
		panic(err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	app.Sessions = MustConfigureSessions(app)
	app.Router = MustConfigureRootRouter(app)
	return app
}

func (app *App) StartApp() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	fmt.Println("Starting server...")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.Router,
	}
	defer func(srv *http.Server, ctx context.Context) {
		_ = srv.Shutdown(ctx)
	}(srv, ctx)
	fmt.Println("Listening on port 8080")
	go func() {
		slog.Error(srv.ListenAndServe().Error())
	}()
	<-ctx.Done()
	fmt.Println("")
	fmt.Println("Goodbye")
}

func doMigrations(app *App) {
	migrator := MustGetNewMigrator(app.MigrationsFS)
	db, err := NewMigrationDb(app.Config.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = migrator.ApplyMigrations(db)
	if err != nil {
		panic(err)
	}
}
