package article

import (
	"context"
	"net/url"
	"rss-reader-backend/internal/service/source"
	"time"
)

type Article struct {
	Id          uint64
	Title       string
	Description string
	Link        string
	PubTime     time.Time
	SourceId    uint64
}

type Service struct {
	store  Store
	client RssClient
}

func NewService(store Store, client RssClient) *Service {
	return &Service{store: store, client: client}
}

func (s *Service) GetArticlePage(ctx context.Context, page int, size int) ([]Article, error) {
	return s.store.GetArticlePage(ctx, page, size)
}

func (s *Service) LoadArticles(ctx context.Context, source source.RssSource) error {
	articles, err := s.client.GetArticles(ctx, source.Url)
	if err != nil {
		return err
	}
	return s.store.SaveArticles(ctx, articles)
}

type Store interface {
	GetArticlePage(ctx context.Context, page int, size int) ([]Article, error)
	SaveArticles(ctx context.Context, articles []Article) error
}

type RssClient interface {
	GetArticles(ctx context.Context, url *url.URL) ([]Article, error)
}
