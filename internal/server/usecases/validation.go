package usecases

import (
	"encoding/json"
	"io"
	"net/http"

	"reimagined_eureka/internal/common"
)

func validateUserAuthReq(w http.ResponseWriter, r *http.Request, entropyRequired bool) *common.UserCredentials {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Supply data as JSON"))
		return nil
	}
	var creds common.UserCredentials
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to read request body"))
		return nil
	}
	if err := json.Unmarshal(body, &creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse user create request"))
		return nil
	}
	if len(creds.Login) < common.MinLoginLength || len(creds.Password) < common.MinPasswordLength {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("userID or password is too short"))
		return nil
	}
	if entropyRequired && (creds.Entropy == nil ||
		creds.Entropy.Hash == nil ||
		creds.Entropy.Ciphertext == nil ||
		creds.Entropy.Salt == nil ||
		creds.Entropy.Nonce == nil) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("All entropy fields must be supplied"))
		return nil
	}
	return &creds
}
