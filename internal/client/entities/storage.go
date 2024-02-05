package entities

import (
	"reimagined_eureka/internal/common"
)

type IStorage interface {
	Init() error
	Shutdown()
	Tx() (ITx, error)
	ReadUser(login string) (*User, error) // TODO: split into separate repos?
	SaveUser(user *User, entropy *common.Entropy) (int, error)
	ReadCredentials(userID int) ([]*CredentialLocal, error)
	SaveCredentials(credentials *CredentialLocal) error
	ReadNotes(userID int) ([]*Note, error)
	ReadFiles(userID int) ([]*File, error)
	ReadBankCards(userID int) ([]*BankCard, error)
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

type CookieEncrypted struct {
	EncryptedValue []byte `db:"value_encrypted"`
	Salt           []byte `db:"salt"`
	Nonce          []byte `db:"nonce"`
}

type CredentialLocal struct {
	common.Credential
	ServerID int `db:"server_id"`
}

type Note struct {
	ID               int    `db:"id"`
	ServerID         int    `db:"server_id"`
	UserID           int    `db:"user_id"`
	Meta             string `db:"meta"`
	EncryptedContent []byte `db:"encrypted_content"`
	Salt             []byte `db:"salt"`
	Nonce            []byte `db:"nonce"`
}

type File struct {
	ID               int    `db:"id"`
	ServerID         int    `db:"server_id"`
	UserID           int    `db:"user_id"`
	Meta             string `db:"meta"`
	EncryptedContent []byte `db:"encrypted_content"`
	Salt             []byte `db:"salt"`
	Nonce            []byte `db:"nonce"`
}

type BankCard struct {
	ID               int    `db:"id"`
	ServerID         int    `db:"server_id"`
	UserID           int    `db:"user_id"`
	Meta             string `db:"meta"`
	EncryptedContent []byte `db:"encrypted_content"`
	Salt             []byte `db:"salt"`
	Nonce            []byte `db:"nonce"`
}
