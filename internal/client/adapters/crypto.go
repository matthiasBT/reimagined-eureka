package adapters

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	clientEntities "reimagined_eureka/internal/client/entities"
)

type CryptoProvider struct{}

func NewCryptoProvider() *CryptoProvider {
	return &CryptoProvider{}
}

func (c *CryptoProvider) VerifyPassword(user *clientEntities.User, password string) error {
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return fmt.Errorf("failed to check password hash: %v", err)
	}
	return nil
}

func (c *CryptoProvider) HashPassword(user *clientEntities.User, password string) error {
	pwdHash, err := c.hashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = pwdHash
	return nil
}

func (c *CryptoProvider) hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	return hashedPassword, nil
}
