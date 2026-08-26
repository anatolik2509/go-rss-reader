package config

import (
	"rss-reader-backend/internal/security"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type Sessions struct {
	sessionManager security.SessionManager
	accountManager security.AccountManager
	scsSessionManager *scs.SessionManager
	authMiddleware *security.AuthMiddleware
	router chi.Router
}

func MustConfigureSessions(app *App) *Sessions {
	sessions := &Sessions{}
	accountRepository := security.NewPgAccountRepository(app.db)
	passwordHasher := &security.BcryptPasswordHasher{}
	sessions.accountManager = security.NewPasswordAccountManager(accountRepository, passwordHasher)
	scsSessionStore := security.NewScsRedisSessionStore(app.redis)
	sessions.scsSessionManager = security.NewScsSessionManager(scsSessionStore)
	sessions.sessionManager = security.NewAppSessionManager(sessions.scsSessionManager)
	sessions.authMiddleware = security.NewAuthMiddleware(sessions.sessionManager, app.logger)
	sessions.router = security.MustConfigureRouter(app.db, app.redis, sessions.accountManager, sessions.sessionManager, app.logger)
	return sessions
}
