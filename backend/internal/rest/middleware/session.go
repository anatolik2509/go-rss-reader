package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type SessionMiddleware struct {
	sessionManager *scs.SessionManager
}

func NewSessionMiddleware(sessionManager *scs.SessionManager) *SessionMiddleware {
	return &SessionMiddleware{sessionManager}
}

func (sm *SessionMiddleware) CreateSessions(next http.Handler) http.Handler {
	return sm.sessionManager.LoadAndSave(next)
}
