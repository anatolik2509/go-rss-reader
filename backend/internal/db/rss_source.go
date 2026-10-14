package db

import (
	"context"
	"fmt"
	"net/url"
	"rss-reader-backend/internal/service/source"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dto struct {
	Id   uint64
	Name string
	Url  string
}

func (s *Dto) parseDto() (*source.RssSource, error) {
	sourceUrl, err := url.Parse(s.Url)
	if err != nil {
		return nil, fmt.Errorf("parsing source url: %w", err)
	}
	parsedSource := source.RssSource{
		Id:   s.Id,
		Name: s.Name,
		Url:  sourceUrl,
	}
	return &parsedSource, err
}

func toDto(s *source.RssSource) Dto {
	return Dto{
		Id:   s.Id,
		Name: s.Name,
		Url:  s.Url.String(),
	}
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool}
}

func (s *PostgresStore) AddSource(ctx context.Context, source source.RssSource) (id uint64, err error) {
	rows, err := s.pool.Query(ctx,
		"INSERT INTO rss_source (name, url) VALUES ($1, $2) RETURNING id", source.Name, source.Url.String(),
	)
	if err != nil {
		return 0, fmt.Errorf("error while insert new source: %w", err)
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("id of inserted row not returned: %w", err)
	}
	return
}

func (s *PostgresStore) GetSources(ctx context.Context) (sources []source.RssSource, err error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name, url FROM rss_source")
	if err != nil {
		return nil, fmt.Errorf("error while querying sources: %w", err)
	}
	defer rows.Close()
	dbSources, err := rowsToSources(rows)
	if err != nil {
		return
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while querying sources: %w", err)
	}

	sources = make([]source.RssSource, 0)
	for _, s := range dbSources {
		parsedSource, err := s.parseDto()
		if err != nil {
			return nil, err
		}
		sources = append(sources, *parsedSource)
	}

	return
}

func rowsToSources(rows pgx.Rows) ([]Dto, error) {
	sources := make([]Dto, 0)
	for rows.Next() {
		var s Dto
		if err := rows.Scan(&s.Id, &s.Name, &s.Url); err != nil {
			return nil, fmt.Errorf("row mapping failed: %w", err)
		}
		sources = append(sources, s)
	}
	return sources, nil
}
