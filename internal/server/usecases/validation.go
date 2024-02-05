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

func validateCredentials(w http.ResponseWriter, r *http.Request) *common.CredentialsReq {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Supply data as JSON"))
		return nil
	}
	var creds common.CredentialsReq
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to read request body"))
		return nil
	}
	if err := json.Unmarshal(body, &creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse credentials create/update request"))
		return nil
	}
	if creds.Login == "" ||
		creds.Meta == "" ||
		creds.Value == nil ||
		creds.Value.Ciphertext == nil ||
		creds.Value.Salt == nil ||
		creds.Value.Nonce == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("All mandatory fields must be supplied"))
	}
	return &creds
}

func validateNote(w http.ResponseWriter, r *http.Request) *common.NoteReq {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Supply data as JSON"))
		return nil
	}
	var note common.NoteReq
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to read request body"))
		return nil
	}
	if err := json.Unmarshal(body, &note); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse note request"))
		return nil
	}
	if note.Meta == "" ||
		note.Value == nil ||
		note.Value.Ciphertext == nil ||
		note.Value.Salt == nil ||
		note.Value.Nonce == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("All mandatory fields must be supplied"))
	}
	return &note
}
