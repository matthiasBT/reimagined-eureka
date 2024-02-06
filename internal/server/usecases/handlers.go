package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"reimagined_eureka/internal/common"
	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/config"

	"github.com/go-chi/chi/v5"
)

func (c *BaseController) signUp(w http.ResponseWriter, r *http.Request) {
	userReq := validateUserAuthReq(w, r, true)
	if userReq == nil {
		return
	}
	pwdhash, err := c.crypto.HashSecurely(userReq.Password)
	if err != nil {
		return
	}
	token := generateSessionToken()
	tx, _ := c.stor.Tx(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create user"))
		return
	}
	defer tx.Commit()
	user, err := c.userRepo.CreateUser(r.Context(), tx, userReq.Login, pwdhash, userReq.Entropy)
	if err != nil {
		defer tx.Rollback()
		if errors.Is(err, entities.ErrLoginAlreadyTaken) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("login is already taken"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to create a new user"))
		}
		return
	}
	session, err := c.userRepo.CreateSession(r.Context(), tx, user, token)
	if err != nil {
		defer tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create a session"))
		return
	}
	authorize(w, session)
}

func (c *BaseController) signIn(w http.ResponseWriter, r *http.Request) {
	userReq := validateUserAuthReq(w, r, false)
	if userReq == nil {
		return
	}
	user, err := c.userRepo.FindUser(r.Context(), userReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to find the user"))
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User doesn't exist"))
		return
	}
	if err := c.crypto.CheckHash(userReq.Password, user.PasswordHash); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect password"))
		return
	}
	token := generateSessionToken()
	tx, _ := c.stor.Tx(r.Context())
	defer tx.Commit()
	if err != nil {
		defer tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create a user session"))
		return
	}
	session, err := c.userRepo.CreateSession(r.Context(), tx, user, token)
	if err != nil {
		defer tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create a user session"))
		return
	}
	entropyEncrypted := common.EncryptionResult{
		Ciphertext: user.EntropyEncrypted,
		Salt:       user.EntropySalt,
		Nonce:      user.EntropyNonce,
	}
	entropy := &common.Entropy{
		EncryptionResult: &entropyEncrypted,
		Hash:             user.EntropyHash,
	}
	entropyData, err := json.Marshal(entropy)
	if err != nil {
		defer tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to write user entropy data"))
		return
	}
	authorize(w, session)
	w.Write(entropyData)
}

func (c *BaseController) writeCredentials(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(w, r)
	if userID == nil {
		return
	}
	creds := validateCredentials(w, r)
	if creds == nil {
		return
	}
	tx, _ := c.stor.Tx(r.Context())
	defer tx.Commit()
	rowId, err := c.credsRepo.Write(r.Context(), tx, *userID, creds)
	if err != nil {
		defer tx.Rollback()
		if errors.Is(err, entities.ErrDoesntExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Errorf("failed to write credentials: %v", err).Error()
		w.Write([]byte(msg))
		return
	}
	w.Write([]byte(strconv.Itoa(rowId)))
}

func (c *BaseController) deleteCredentials(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(w, r)
	if userID == nil {
		return
	}
	rowIDRaw := chi.URLParam(r, "rowID")
	rowID, err := strconv.Atoi(rowIDRaw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Errorf("not a valid row ID: %v", err).Error()
		w.Write([]byte(msg))
		return
	}
	tx, _ := c.stor.Tx(r.Context())
	defer tx.Commit()
	if err := c.credsRepo.Delete(r.Context(), tx, *userID, rowID); err != nil {
		defer tx.Rollback()
		if errors.Is(err, entities.ErrDoesntExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Errorf("failed to delete credentials: %v", err).Error()
		w.Write([]byte(msg))
		return
	}
}

func (c *BaseController) deleteNote(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(w, r)
	if userID == nil {
		return
	}
	rowIDRaw := chi.URLParam(r, "noteID")
	rowID, err := strconv.Atoi(rowIDRaw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Errorf("not a valid row ID: %v", err).Error()
		w.Write([]byte(msg))
		return
	}
	tx, _ := c.stor.Tx(r.Context())
	defer tx.Commit()
	if err := c.notesRepo.Delete(r.Context(), tx, *userID, rowID); err != nil {
		defer tx.Rollback()
		if errors.Is(err, entities.ErrDoesntExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Errorf("failed to delete note: %v", err).Error()
		w.Write([]byte(msg))
		return
	}
}

func (c *BaseController) writeNote(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(w, r)
	if userID == nil {
		return
	}
	note := validateNote(w, r)
	if note == nil {
		return
	}
	tx, _ := c.stor.Tx(r.Context())
	defer tx.Commit()
	rowId, err := c.notesRepo.Write(r.Context(), tx, *userID, note)
	if err != nil {
		defer tx.Rollback()
		if errors.Is(err, entities.ErrDoesntExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Errorf("failed to write note: %v", err).Error()
		w.Write([]byte(msg))
		return
	}
	w.Write([]byte(strconv.Itoa(rowId)))
}

func (c *BaseController) writeFile(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(w, r)
	if userID == nil {
		return
	}
	file := validateFile(w, r)
	if file == nil {
		return
	}
	tx, _ := c.stor.Tx(r.Context())
	defer tx.Commit()
	rowId, err := c.filesRepo.Write(r.Context(), tx, *userID, file)
	if err != nil {
		defer tx.Rollback()
		if errors.Is(err, entities.ErrDoesntExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Errorf("failed to write file: %v", err).Error()
		w.Write([]byte(msg))
		return
	}
	w.Write([]byte(strconv.Itoa(rowId)))
}

func (c *BaseController) writeCard(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(w, r)
	if userID == nil {
		return
	}
	file := validateCard(w, r)
	if file == nil {
		return
	}
	tx, _ := c.stor.Tx(r.Context())
	defer tx.Commit()
	rowId, err := c.cardsRepo.Write(r.Context(), tx, *userID, file)
	if err != nil {
		defer tx.Rollback()
		if errors.Is(err, entities.ErrDoesntExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Errorf("failed to write card: %v", err).Error()
		w.Write([]byte(msg))
		return
	}
	w.Write([]byte(strconv.Itoa(rowId)))
}

func (c *BaseController) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func authorize(w http.ResponseWriter, session *entities.Session) {
	cookie := http.Cookie{
		Name:     common.SessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  time.Now().Add(config.SessionTTL),
		HttpOnly: true,  // Protect against XSS attacks
		Secure:   false, // Should be true in production to send only over HTTPS
	}
	http.SetCookie(w, &cookie)
}

func getUserID(w http.ResponseWriter, r *http.Request) *int {
	userID := r.Context().Value(entities.ContextKey{Key: "user_id"})
	if userID == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to find the user_id in the context"))
		return nil
	}
	res := userID.(int)
	return &res
}
