package adapters

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"

	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/client/infra/migrations"
)

type SQLiteStorage struct {
	db *sqlx.DB
}

func NewSQLiteStorage(logger logging.ILogger, path string) (*SQLiteStorage, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Warning("Database doesn't exist. Creating a new database")
	}
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %v", err)
	}
	if err := migrations.Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate the database: %v", err)
	}
	storage := &SQLiteStorage{db}
	return storage, nil
}

func (s *SQLiteStorage) Init() error {
	return nil
}
