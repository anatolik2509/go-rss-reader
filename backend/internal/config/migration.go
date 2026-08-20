package config

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

const migrationsDir = "migrations"

type Migrator struct {
	srcDriver source.Driver
}

func MustGetNewMigrator(MigrationsFS fs.FS) *Migrator {
	d, err := iofs.New(MigrationsFS, migrationsDir)
	if err != nil {
		panic(err)
	}
	return &Migrator{
		srcDriver: d,
	}
}

func (m *Migrator) ApplyMigrations(db *sql.DB) (err error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	defer db.Close()
	if err != nil {
		return fmt.Errorf("unable to create db instance: %w", err)
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "psql_db", driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %w", err)
	}

	defer func() {
		closeError := errors.Join(migrator.Close())
		if err == nil {
			err = closeError
		}
	}()

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations %w", err)
	}

	return nil
}
