package job

import (
	"context"
	"log/slog"
	"rss-reader-backend/internal/service/article"
	"rss-reader-backend/internal/service/source"
	"time"
)

type RssLoader struct {
	sourceService  source.Service
	articleService article.Service
}

func NewRssLoader() *RssLoader {
	return &RssLoader{}
}

func (r *RssLoader) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := r.LoadArticles(ctx)
			if err != nil {
				slog.WarnContext(ctx, "Failed to load articles", slog.String("error", err.Error()))
			}
		}
	}
}

func (r *RssLoader) LoadArticles(ctx context.Context) error {
	sources, err := r.sourceService.GetSources(ctx)
	if err != nil {
		return err
	}
	for _, rssSource := range sources {
		err = r.articleService.LoadArticles(ctx, rssSource)
		if err != nil {
			return err
		}
	}
	return nil
}
