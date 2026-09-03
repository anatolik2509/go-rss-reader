package main

import (
	"embed"
	"rss-reader-backend/internal/config"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func main() {
	config.MustCreateNewApp(MigrationsFS).StartApp()
}
