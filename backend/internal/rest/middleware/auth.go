package middleware

import (
	"log/slog"
	"net/http"
	"rss-reader-backend/internal/session"
)

type AuthMiddleware struct {
	sessionManager session.Manager
}

func NewAuthMiddleware(sessionManager session.Manager) *AuthMiddleware {
	return &AuthMiddleware{sessionManager}
}

func (middleware *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok, err := middleware.sessionManager.CheckSession(r, w)
		if err != nil {
			slog.ErrorContext(r.Context(), "Session check failed", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
