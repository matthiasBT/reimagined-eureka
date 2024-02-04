package entities

import "reimagined_eureka/internal/common"

type ICryptoProvider interface {
	VerifyPassword(user *User, password string) error
	HashPassword(user *User, password string) error
	HashSecurely(secret string) ([]byte, error)
	SetMasterKey(masterKey string)
	Encrypt(what string) (*common.EncryptionResult, error)
}
