package adapters

import clientEntities "reimagined_eureka/internal/client/entities"

type CryptoProvider struct{}

func NewCryptoProvider() *CryptoProvider {
	return &CryptoProvider{}
}

func (c *CryptoProvider) HashPassword(password string) ([]byte, error) {
	panic("implement me")
}

func (c *CryptoProvider) GenerateSalt() ([]byte, error) {
	panic("implement me")
}

func (c *CryptoProvider) VerifyPassword(user *clientEntities.User, password string) (bool, error) {
	panic("implement me")
}

func (c *CryptoProvider) PrepareUserForSave(user *clientEntities.User) error {
	user.PasswordHash = []byte{}
	user.PasswordSalt = []byte{}
	return nil
}
