package entities

type ICryptoProvider interface {
	HashSecurely(secret string) ([]byte, error)
	CheckHash(secret string, hash []byte) error
}
