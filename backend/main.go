package main

import (
	"embed"
	"rss-reader-backend/internal/bootstrap"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func main() {
	bootstrap.MustCreateNewApp(MigrationsFS).StartApp()
}
