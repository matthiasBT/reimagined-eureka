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
	logger logging.ILogger
	db     *sqlx.DB
}

func NewSQLiteStorage(logger logging.ILogger, path string) (*SQLiteStorage, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Warningln("Database doesn't exist. Creating a new database")
	}
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %v", err)
	}
	storage := &SQLiteStorage{logger: logger, db: db}
	if err := storage.Init(); err != nil {
		return nil, fmt.Errorf("failed to init the database: %v", err)
	}
	return storage, nil
}

func (s *SQLiteStorage) Init() error {
	return migrations.Migrate(s.db)
}

func (s *SQLiteStorage) Shutdown() {
	if err := s.db.Close(); err != nil {
		s.logger.Failureln("Failed to shut down the database: %v", err)
	}
}
