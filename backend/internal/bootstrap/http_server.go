package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MustConfigureRootRouter(app *App) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(app.Sessions.sessionMiddleware.CreateSessions)

	router.Mount("/auth", app.Sessions.router)
	sourceRouter := MustConfigureRouter(app)
	router.Route("/api", func(r chi.Router) {
		r.Use(app.Sessions.authMiddleware.RequireAuth)
		r.Mount("/sources", sourceRouter)
	})
	return router
}

// todo добавить chi render
func MustConfigureRouter(app *App) chi.Router {
	rssSourceHandler := ConfigureRssSourceHandler(app)
	router := chi.NewRouter()
	router.Get("/", rssSourceHandler.GetSources)
	router.Post("/", rssSourceHandler.AddSource)
	return router
}
