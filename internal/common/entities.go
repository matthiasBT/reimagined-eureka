package common

type Credentials struct {
	Login    string            `json:"login"`
	Password string            `json:"password"`
	Entropy  *EncryptionResult `json:"entropy"`
}

type EncryptionResult struct {
	Plaintext               string // TODO: remove this!
	Ciphertext, Salt, Nonce []byte
}

// TODO: add a separate entity for entropy
