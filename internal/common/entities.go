package common

type Credentials struct {
	Login    string            `json:"login"`
	Password string            `json:"password"`
	Entropy  *EncryptionResult `json:"entropy"`
}

type EncryptionResult struct {
	Result, Salt, Nonce []byte
}
