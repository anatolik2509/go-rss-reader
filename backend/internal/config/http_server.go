package config

import (
	"rss-reader-backend/internal/domain/source"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MustConfigureRootRouter(app *App) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(app.sessions.scsSessionManager.LoadAndSave)
	
	router.Mount("/auth", app.sessions.router)
	sourceRouter := source.MustConfigureRouter(app.db)
	router.Route("/api", func(r chi.Router) {
		r.Use(app.sessions.authMiddleware.RequireAuth)
		r.Mount("/sources", sourceRouter)
	})
	return router
}
