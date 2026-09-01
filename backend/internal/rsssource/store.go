package source

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sync/atomic"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type rssSource struct {
	Id   uint64
	Name string
	Url  string
}

type Store interface {
	AddSource(ctx context.Context, source rssSource) (id uint64, err error)
	GetSources(ctx context.Context) ([]rssSource, error)
}

type InMemoryStore struct {
	counter atomic.Uint64
	store   map[uint64]rssSource
}

func NewInMemoryStore() *InMemoryStore {
	var s = InMemoryStore{store: make(map[uint64]rssSource), counter: atomic.Uint64{}}
	return &s
}

func (s *InMemoryStore) AddSource(ctx context.Context, source rssSource) (id uint64, err error) {
	id = s.counter.Add(1)
	s.store[id] = source
	return
}

func (s *InMemoryStore) GetSources(ctx context.Context) ([]rssSource, error) {
	sources := slices.Collect(maps.Values(s.store))
	return sources, nil
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool}
}

func (s *PostgresStore) AddSource(ctx context.Context, source rssSource) (id uint64, err error) {
	rows, err := s.pool.Query(ctx,
		"INSERT INTO rss_source (name, url) VALUES ($1, $2) RETURNING id", source.Name, source.Url,
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

func (s *PostgresStore) GetSources(ctx context.Context) (sources []rssSource, err error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name, url FROM rss_source")
	if err != nil {
		return nil, fmt.Errorf("error while querying sources: %w", err)
	}
	defer rows.Close()
	sources, err = rowsToSources(sources, rows)
	if err != nil {
		return
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while querying sources: %w", err)
	}
	return
}

func rowsToSources(sources []rssSource, rows pgx.Rows) ([]rssSource, error) {
	sources = make([]rssSource, 0)
	for rows.Next() {
		var s rssSource
		if err := rows.Scan(&s.Id, &s.Name, &s.Url); err != nil {
			return nil, fmt.Errorf("row mapping failed: %w", err)
		}
		sources = append(sources, s)
	}
	return sources, nil
}
