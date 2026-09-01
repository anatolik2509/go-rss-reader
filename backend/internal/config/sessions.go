package config

import (
	"rss-reader-backend/internal/account"
	"rss-reader-backend/internal/session"

	"github.com/go-chi/chi/v5"
)

type Sessions struct {
	sessionManager    session.Manager
	accountManager    account.AccountManager
	authMiddleware    *session.AuthMiddleware
	sessionMiddleware *session.SessionMiddleware
	router            chi.Router
}

func MustConfigureSessions(app *App) *Sessions {
	sessions := &Sessions{}
	accountRepository := account.NewPgAccountRepository(app.db)
	passwordHasher := &account.BcryptPasswordHasher{}
	sessions.accountManager = account.NewPasswordAccountManager(accountRepository, passwordHasher)
	sessions.sessionManager, sessions.sessionMiddleware, sessions.authMiddleware = session.New(app.redis, app.logger)
	sessions.router = account.MustConfigureRouter(app.db, app.redis, sessions.accountManager, sessions.sessionManager, app.logger)
	return sessions
}
