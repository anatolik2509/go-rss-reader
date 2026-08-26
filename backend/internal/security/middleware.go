package security

import (
	"log/slog"
	"net/http"
)

type AuthMiddleware struct {
	sessionManager SessionManager
	logger *slog.Logger
}

func NewAuthMiddleware(sessionManager SessionManager, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{sessionManager, logger}
}

func (middleware *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok, err := middleware.sessionManager.CheckSession(r, w)
		if err != nil {
			middleware.logger.ErrorContext(r.Context(), "Session check failed", "error", err)
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
