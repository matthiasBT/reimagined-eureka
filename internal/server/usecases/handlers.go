package usecases

import (
	"errors"
	"net/http"
	"time"

	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/config"
)

func (c *BaseController) register(w http.ResponseWriter, r *http.Request) {
	userReq := validateUserAuthReq(w, r)
	if userReq == nil {
		return
	}
	pwdhash, err := c.crypto.HashPassword(userReq.Password)
	if err != nil {
		return
	}
	token := generateSessionToken()
	tx, err := c.stor.Tx(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create user"))
		return
	}
	defer tx.Commit()
	user, err := c.userRepo.CreateUser(r.Context(), tx, userReq.Login, pwdhash)
	if err != nil {
		defer tx.Rollback()
		if errors.Is(err, entities.ErrLoginAlreadyTaken) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Login is already taken"))
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
	userReq := validateUserAuthReq(w, r)
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
	if err := c.crypto.CheckPassword(userReq.Password, user.PasswordHash); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect password"))
		return
	}
	token := generateSessionToken()
	tx, err := c.stor.Tx(r.Context())
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
	authorize(w, session)
}

func (c *BaseController) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func authorize(w http.ResponseWriter, session *entities.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Path:     "/",
		Expires:  time.Now().Add(config.SessionTTL),
		HttpOnly: true,  // Protect against XSS attacks
		Secure:   false, // Should be true in production to send only over HTTPS
	})
}
