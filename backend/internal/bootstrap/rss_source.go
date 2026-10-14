package bootstrap

import (
	"rss-reader-backend/internal/db"
	"rss-reader-backend/internal/rest"
	"rss-reader-backend/internal/service/source"
)

func ConfigureRssSourceHandler(app *App) *rest.RssSourceHandler {
	store := db.NewPostgresStore(app.Db)
	sourceService := source.NewService(store)
	return rest.NewRssSourceHandler(sourceService)
}
