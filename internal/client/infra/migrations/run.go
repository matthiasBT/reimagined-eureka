package migrations

import (
	"embed"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed sqlite/*.sql
var embedMigrations embed.FS

func Migrate(db *sqlx.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}
	if err := goose.Up(db.DB, "sqlite"); err != nil {
		return err
	}
	return nil
}
