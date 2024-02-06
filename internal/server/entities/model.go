package entities

import (
	"time"
)

type ContextKey struct {
	Key string
}

type User struct {
	ID               int    `db:"id"`
	Login            string `db:"login"`
	PasswordHash     []byte `db:"password_hash"`
	EntropyHash      []byte `db:"entropy_hash"`
	EntropyEncrypted []byte `db:"entropy_encrypted"`
	EntropySalt      []byte `db:"entropy_salt"`
	EntropyNonce     []byte `db:"entropy_nonce"`
}

type Session struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}
