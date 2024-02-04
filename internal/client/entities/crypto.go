package entities

import "reimagined_eureka/internal/common"

type ICryptoProvider interface {
	HashPassword(user *User, password string) error
	VerifyPassword(user *User, password string) error
	Hash(data []byte) ([]byte, error)
	VerifyHash(data, target []byte) error
	SetMasterKey(masterKey string)
	Encrypt(what []byte) (*common.EncryptionResult, error)
	Decrypt(what, salt, nonce []byte) ([]byte, error)
}
