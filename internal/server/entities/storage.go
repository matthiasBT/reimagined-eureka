package entities

import (
	"context"
	"errors"

	"reimagined_eureka/internal/common"
)

var (
	ErrLoginAlreadyTaken = errors.New("login already taken")
)

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
	FindUser(ctx context.Context, request *common.Credentials) (*User, error)
	CreateSession(ctx context.Context, tx Tx, user *User, token string) (*Session, error)
	FindSession(ctx context.Context, token string) (*Session, error)
}
