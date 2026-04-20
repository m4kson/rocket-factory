package migrator

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db             *sql.DB
	migrationsPath string
}

func NewMigrator(db *sql.DB, migrationsPath string) *Migrator {
	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
	}
}

func (m *Migrator) Up() error {
	err := goose.Up(m.db, m.migrationsPath)
	if err != nil {
		return err
	}

	return nil
}
