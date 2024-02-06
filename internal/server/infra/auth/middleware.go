package auth

import (
	"context"
	"net/http"
	"time"

	"reimagined_eureka/internal/common"
	"reimagined_eureka/internal/server/entities"
	"reimagined_eureka/internal/server/infra/logging"
)

func Middleware(logger logging.ILogger, userRepo entities.UserRepo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		checkAuthFn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && (r.URL.Path == "/api/user/register" || r.URL.Path == "/api/user/login") {
				logger.Infoln("No auth check necessary")
				next.ServeHTTP(w, r)
				return
			}
			cookie, err := r.Cookie(common.SessionCookieName)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Missing session cookie"))
				return
			}
			session, err := userRepo.FindSession(r.Context(), cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to find a session"))
				return
			}
			if session == nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Session not found"))
				return
			}
			if time.Now().After(session.ExpiresAt) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Session has expired"))
				return
			}
			ctx := context.WithValue(r.Context(), entities.ContextKey{Key: "user_id"}, session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(checkAuthFn)
	}
}
