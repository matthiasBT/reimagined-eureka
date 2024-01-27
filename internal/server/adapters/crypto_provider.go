package adapters

import (
	"golang.org/x/crypto/bcrypt"

	"reimagined_eureka/internal/server/infra/logging"
)

type CryptoProvider struct {
	Logger logging.ILogger
}

func (cr *CryptoProvider) HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		cr.Logger.Errorf("Failed to hash password: %s", err.Error())
		return nil, err
	}
	return hashedPassword, nil
}

func (cr *CryptoProvider) CheckPassword(password string, hash []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		cr.Logger.Errorf("Password hash didn't match the password: %s", err.Error())
		return err
	}
	return nil
}
