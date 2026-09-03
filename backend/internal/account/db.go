package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountDto struct {
	Id           uint64
	Login        string
	PasswordHash string
}

type AccountRepository interface {
	Save(context.Context, AccountDto) (id uint64, err error)
	Find(ctx context.Context, login string) (AccountDto, error)
}

const (
	CreateAccountQuery = "INSERT INTO account(login, password_hash) VALUES ($1, $2) RETURNING id;"
	FindAccountQuery   = "SELECT id, login, password_hash FROM account WHERE login = $1;"
)

type PgAccountRepository struct {
	pool *pgxpool.Pool
}

func NewPgAccountRepository(pool *pgxpool.Pool) AccountRepository {
	return &PgAccountRepository{pool: pool}
}

func (r *PgAccountRepository) Save(ctx context.Context, acc AccountDto) (id uint64, err error) {
	rows, err := r.pool.Query(ctx, CreateAccountQuery, acc.Login, acc.PasswordHash)
	if err != nil {
		return 0, fmt.Errorf("saving account to db: %w", err)
	}
	id, err = pgx.CollectExactlyOneRow[uint64](rows, pgx.RowTo)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return 0, fmt.Errorf("saving account to db: %w", ErrSuchLoginAlreadyExists)
	}
	if err != nil {
		return 0, fmt.Errorf("saving account to db: %w", err)
	}
	return id, nil
}

func (r *PgAccountRepository) Find(ctx context.Context, login string) (AccountDto, error) {
	rows, err := r.pool.Query(ctx, FindAccountQuery, login)
	if err != nil {
		return AccountDto{}, fmt.Errorf("searching account in db: %w", err)
	}
	var account AccountDto
	account, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[AccountDto])
	if errors.Is(err, pgx.ErrNoRows) {
		return AccountDto{}, fmt.Errorf("searching account in db: %w", ErrNoSuchAccount)
	}
	if err != nil {
		return AccountDto{}, fmt.Errorf("searching account in db: %w", err)
	}
	return account, nil
}
