package source

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MustConfigureRouter(pool *pgxpool.Pool, logger *slog.Logger) chi.Router {
	store := NewPostgresStore(pool)
	handler := &handler{store: store}
	router := chi.NewRouter()
	router.Get("/", handler.GetSources)
	router.Post("/", handler.AddSource)
	return router
}
