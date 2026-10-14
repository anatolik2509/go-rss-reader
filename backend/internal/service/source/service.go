package source

import (
	"context"
	"net/url"
)

type RssSource struct {
	Id   uint64
	Name string
	Url  *url.URL
}

type Store interface {
	AddSource(ctx context.Context, source RssSource) (id uint64, err error)
	GetSources(ctx context.Context) ([]RssSource, error)
}

type Service struct {
	store Store
}

func (s *Service) GetSources(ctx context.Context) ([]RssSource, error) {
	return s.store.GetSources(ctx)
}

func (s *Service) AddSource(ctx context.Context, source RssSource) (RssSource, error) {
	id, err := s.store.AddSource(ctx, source)
	if err != nil {
		return RssSource{}, err
	}
	return RssSource{
		Id:   id,
		Name: source.Name,
		Url:  source.Url,
	}, nil
}

func NewService(store Store) *Service {
	return &Service{store: store}
}
