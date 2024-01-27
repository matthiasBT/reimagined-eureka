package usecases

import (
	"encoding/json"
	"io"
	"net/http"

	"reimagined_eureka/internal/server/entities"
)

const MinLoginLength = 5
const MinPasswordLength = 5

func validateUserAuthReq(w http.ResponseWriter, r *http.Request) *entities.UserAuthRequest {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Supply data as JSON"))
		return nil
	}
	var userReq entities.UserAuthRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to read request body"))
		return nil
	}
	if err := json.Unmarshal(body, &userReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse user create request"))
		return nil
	}
	if len(userReq.Login) < MinLoginLength || len(userReq.Password) < MinPasswordLength {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Login or password is too short"))
		return nil
	}
	return &userReq
}
