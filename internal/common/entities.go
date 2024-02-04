package common

type UserCredentials struct {
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Entropy  *Entropy `json:"entropy"`
}

type EncryptionResult struct {
	Ciphertext, Salt, Nonce []byte
}

type Entropy struct {
	*EncryptionResult
	Hash []byte
}
