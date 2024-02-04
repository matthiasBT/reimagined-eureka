package entities

import "reimagined_eureka/internal/common"

type IStorage interface {
	Init() error
	Shutdown()
	Tx() (ITx, error)
	ReadUser(login string) (*User, error) // TODO: split into separate repos?
	SaveUser(user *User, entropy *common.Entropy) error
	ReadCredentials(login, what string) ([]*Credential, error)
}

type ITx interface {
	Commit() error
	Rollback() error
}

type User struct {
	ID           int    `db:"id"`
	Login        string `db:"login"`
	PasswordHash []byte `db:"pwd_hash"`
	// TODO: split into a separate entity
	EntropyHash      []byte `db:"entropy_hash"`
	EntropyEncrypted []byte `db:"entropy_encrypted"`
	EntropySalt      []byte `db:"entropy_salt"`
	EntropyNonce     []byte `db:"entropy_nonce"`
}

type Credential struct {
	ID                int    `db:"id"`
	Purpose           string `db:"purpose"`
	Login             string `db:"login"`
	EncryptedPassword []byte `db:"encrypted_password"`
	Nonce             []byte `db:"nonce"`
	Salt              []byte `db:"salt"`
}
