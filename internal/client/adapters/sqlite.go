package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"

	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/client/infra/logging"
	"reimagined_eureka/internal/client/infra/migrations"
	"reimagined_eureka/internal/common"
)

var txOpt = sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  false,
}

type SQLiteTxx struct {
	tx *sqlx.Tx
}

func (t *SQLiteTxx) Commit() error {
	return t.tx.Commit()
}

func (t *SQLiteTxx) Rollback() error {
	return t.tx.Rollback()
}

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
	s.logger.Debugln("Closing the database. Your data will be saved")
	if err := s.db.Close(); err != nil {
		s.logger.Failureln("Failed to shut down the database: %v", err)
	}
}

func (s *SQLiteStorage) Tx() (clientEntities.ITx, error) {
	tx, err := s.db.BeginTxx(context.Background(), &txOpt)
	if err != nil {
		return nil, err
	}
	return &SQLiteTxx{tx: tx}, nil
}

func (s *SQLiteStorage) ReadUser(login string) (*clientEntities.User, error) {
	var user = clientEntities.User{}
	query := "select * from users where login = $1"
	if err := s.db.Get(&user, query, login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user %s: %v", login, err)
	}
	return &user, nil
}

func (s *SQLiteStorage) SaveUser(user *clientEntities.User, entropy *common.Entropy) (int, error) {
	query := `
		insert into users(login, pwd_hash, entropy_hash, entropy_encrypted, entropy_salt, entropy_nonce)
		values ($1, $2, $3, $4, $5, $6)
	`
	result, err := s.db.Exec(
		query,
		user.Login,
		user.PasswordHash,
		entropy.Hash,
		entropy.Ciphertext,
		entropy.Salt,
		entropy.Nonce,
	)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil // TODO: think about proper type conversion and type choices
}

func (s *SQLiteStorage) ReadCredentials(userID int) ([]*clientEntities.CredentialLocal, error) {
	var creds []*clientEntities.CredentialLocal
	query := "select * from credentials where user_id = $1"
	if err := s.db.Select(&creds, query, userID); err != nil {
		return nil, err
	}
	return creds, nil
}

func (s *SQLiteStorage) SaveCredentials(credentials *clientEntities.CredentialLocal) error {
	query := `
		insert into credentials(server_id, user_id, meta, login, encrypted_password, salt, nonce)
		values ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.db.Exec(
		query,
		credentials.ServerID,
		credentials.UserID,
		credentials.Meta,
		credentials.Login,
		credentials.EncryptedPassword,
		credentials.Salt,
		credentials.Nonce,
	)
	return err
}

func (s *SQLiteStorage) ReadNotes(userID int) ([]*clientEntities.Note, error) {
	var notes []*clientEntities.Note
	query := "select * from notes where user_id = $1"
	if err := s.db.Select(&notes, query, userID); err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *SQLiteStorage) ReadFiles(userID int) ([]*clientEntities.File, error) {
	var files []*clientEntities.File
	query := "select * from files where user_id = $1"
	if err := s.db.Select(&files, query, userID); err != nil {
		return nil, err
	}
	return files, nil
}

func (s *SQLiteStorage) ReadBankCards(userID int) ([]*clientEntities.BankCard, error) {
	var cards []*clientEntities.BankCard
	query := "select * from bank_cards where user_id = $1"
	if err := s.db.Select(&cards, query, userID); err != nil {
		return nil, err
	}
	return cards, nil
}
