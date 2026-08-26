package security

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func MustConfigureRouter(pool *pgxpool.Pool, redisClient *redis.Client, accountManager AccountManager, sessionManager SessionManager, logger *slog.Logger) chi.Router {
	httpHandler := NewHttpHandler(accountManager, sessionManager, logger)
	router := chi.NewRouter()
	router.Post("/signUp", httpHandler.SignUp)
	router.Post("/signIn", httpHandler.SignIn)
	return router
}
