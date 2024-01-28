package entities

import "reimagined_eureka/internal/common"

type IStorage interface {
	Init() error
	Shutdown()
	Tx() (ITx, error)
	ReadUser(login string) (*User, error)
	SaveUser(user *User, entropy *common.EncryptionResult) error
}

type ITx interface {
	Commit() error
	Rollback() error
}

type User struct {
	Login        string `db:"login"`
	PasswordHash []byte `db:"pwd_hash"`
}
