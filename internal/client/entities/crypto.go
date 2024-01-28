package entities

// TODO: make some methods private

type ICryptoProvider interface {
	VerifyPassword(user *User, password string) (bool, error)
	PrepareUserForSave(user *User) error

	HashPassword(password string) ([]byte, error)
	GenerateSalt() ([]byte, error)
}
