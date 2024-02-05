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

type CredentialsReq struct {
	ServerID    *int
	Login, Meta string
	Value       *EncryptionResult
}

type Credential struct {
	ID                int    `db:"id"`
	UserID            int    `db:"user_id"`
	Meta              string `db:"meta"`
	Login             string `db:"login"`
	EncryptedPassword []byte `db:"encrypted_password"`
	Salt              []byte `db:"salt"`
	Nonce             []byte `db:"nonce"`
}

type NoteReq struct {
	ServerID *int
	Meta     string
	Value    *EncryptionResult
}

type Note struct {
	ID               int    `db:"id"`
	UserID           int    `db:"user_id"`
	Meta             string `db:"meta"`
	EncryptedContent []byte `db:"encrypted_content"`
	Salt             []byte `db:"salt"`
	Nonce            []byte `db:"nonce"`
}

type FileReq struct {
	ServerID *int
	Meta     string
	Value    *EncryptionResult
}

type File struct {
	ID               int    `db:"id"`
	UserID           int    `db:"user_id"`
	Meta             string `db:"meta"`
	EncryptedContent []byte `db:"encrypted_content"`
	Salt             []byte `db:"salt"`
	Nonce            []byte `db:"nonce"`
}
