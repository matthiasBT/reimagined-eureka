package entities

import (
	"context"
	"errors"

	"reimagined_eureka/internal/common"
)

var (
	ErrLoginAlreadyTaken = errors.New("login already taken")
)

const DefaultVersion = 1

type Tx interface {
	Commit() error
	Rollback() error
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) error
}

type Storage interface {
	Tx(ctx context.Context) (Tx, error)
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	GetContext(ctx context.Context, dest any, query string, args ...any) error
}

type UserRepo interface {
	CreateUser(ctx context.Context, tx Tx, login string, pwdhash []byte, entropy *common.Entropy) (
		*User, error,
	)
	FindUser(ctx context.Context, request *common.UserCredentials) (*User, error)
	CreateSession(ctx context.Context, tx Tx, user *User, token string) (*Session, error)
	FindSession(ctx context.Context, token string) (*Session, error)
}

type CredentialsRepo interface {
	Write(ctx context.Context, tx Tx, userID int, data *common.CredentialsReq) (int, error)
	Read(ctx context.Context, tx Tx, userID int, rowID int) (*common.CredentialsReq, error)
	// ReadVersion
}

type NotesRepo interface {
	Write(ctx context.Context, tx Tx, userID int, data *common.NoteReq) (int, error)
	Read(ctx context.Context, tx Tx, userID int, rowID int) (*common.NoteReq, error)
	// ReadVersion
}

type FilesRepo interface {
	Write(ctx context.Context, tx Tx, userID int, data *common.FileReq) (int, error)
	Read(ctx context.Context, tx Tx, userID int, rowID int) (*common.FileReq, error)
	// ReadVersion
}
