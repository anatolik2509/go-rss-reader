package bootstrap

import (
	"rss-reader-backend/internal/db"
	"rss-reader-backend/internal/rest"
	"rss-reader-backend/internal/rest/middleware"
	"rss-reader-backend/internal/service/account"
	"rss-reader-backend/internal/session"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type Sessions struct {
	sessionStore      scs.Store
	scsSessionManager *scs.SessionManager
	sessionManager    session.Manager
	accountManager    account.Manager
	authMiddleware    *middleware.AuthMiddleware
	sessionMiddleware *middleware.SessionMiddleware
	router            chi.Router
}

func MustConfigureSessions(app *App) *Sessions {
	sessions := &Sessions{}
	accountRepository := db.NewPgAccountRepository(app.Db)
	passwordHasher := &account.BcryptPasswordHasher{}
	sessions.sessionStore = session.NewScsRedisStore(app.Redis)
	sessions.scsSessionManager = session.NewScsManager(sessions.sessionStore)
	sessions.sessionManager = session.NewManager(sessions.scsSessionManager)
	sessions.accountManager = account.NewPasswordAccountManager(accountRepository, passwordHasher)
	sessions.authMiddleware = middleware.NewAuthMiddleware(sessions.sessionManager)
	sessions.sessionMiddleware = middleware.NewSessionMiddleware(sessions.scsSessionManager)
	sessions.router = rest.MustConfigureRouter(sessions.accountManager, sessions.sessionManager)
	return sessions
}
