package security

import (
	"log"
	"net/http"
)

type AuthMiddleware struct {
	sessionManager SessionManager
}

func NewAuthMiddleware(sessionManager *SessionManager) *AuthMiddleware {
	return &AuthMiddleware{*sessionManager}
}

func (middleware *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok, err := middleware.sessionManager.CheckSession(r, w)
		if err != nil {
			log.Printf("error checking session cookie %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}
