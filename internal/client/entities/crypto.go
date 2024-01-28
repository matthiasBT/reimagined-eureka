package entities

// TODO: make some methods private

type ICryptoProvider interface {
	VerifyPassword(user *User, password string) error
	HashPassword(user *User, password string) error
}
