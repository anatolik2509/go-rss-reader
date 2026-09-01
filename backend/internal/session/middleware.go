package session

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type AuthMiddleware struct {
	sessionManager Manager
	logger         *slog.Logger
}

func NewAuthMiddleware(sessionManager Manager, logger *slog.Logger) *AuthMiddleware {
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

type SessionMiddleware struct {
	sessionManager *scs.SessionManager
	logger         *slog.Logger
}

func newSessionMiddleware(sessionManager *scs.SessionManager, logger *slog.Logger) *SessionMiddleware {
	return &SessionMiddleware{sessionManager, logger}
}

func (sm *SessionMiddleware) CreateSessions(next http.Handler) http.Handler {
	return sm.sessionManager.LoadAndSave(next)
}
