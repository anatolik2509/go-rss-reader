package db

import (
	"context"
	"errors"
	"fmt"
	"rss-reader-backend/internal/service/account"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CreateAccountQuery = "INSERT INTO account(Login, password_hash) VALUES ($1, $2) RETURNING id;"
	FindAccountQuery   = "SELECT id, Login, password_hash FROM account WHERE Login = $1;"
)

type PgAccountRepository struct {
	pool *pgxpool.Pool
}

func NewPgAccountRepository(pool *pgxpool.Pool) account.Repository {
	return &PgAccountRepository{pool: pool}
}

func (r *PgAccountRepository) Save(ctx context.Context, acc account.Model) (id uint64, err error) {
	rows, err := r.pool.Query(ctx, CreateAccountQuery, acc.Login, acc.PasswordHash)
	if err != nil {
		return 0, fmt.Errorf("saving account to db: %w", err)
	}
	id, err = pgx.CollectExactlyOneRow[uint64](rows, pgx.RowTo)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return 0, fmt.Errorf("saving account to db: %w", account.ErrSuchLoginAlreadyExists)
	}
	if err != nil {
		return 0, fmt.Errorf("saving account to db: %w", err)
	}
	return id, nil
}

func (r *PgAccountRepository) Find(ctx context.Context, login string) (account.Model, error) {
	rows, err := r.pool.Query(ctx, FindAccountQuery, login)
	if err != nil {
		return account.Model{}, fmt.Errorf("searching account in db: %w", err)
	}
	var model account.Model
	model, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[account.Model])
	if errors.Is(err, pgx.ErrNoRows) {
		return account.Model{}, fmt.Errorf("searching account in db: %w", account.ErrNoSuchAccount)
	}
	if err != nil {
		return account.Model{}, fmt.Errorf("searching account in db: %w", err)
	}
	return model, nil
}
