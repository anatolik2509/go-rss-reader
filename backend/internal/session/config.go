package session

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func New(redisClient *redis.Client, logger *slog.Logger) (m Manager, sm *SessionMiddleware, am *AuthMiddleware) {
	store := newScsRedisStore(redisClient)
	scsManager := newScsManager(store)
	m = newManager(scsManager)
	sm = newSessionMiddleware(scsManager, logger)
	am = NewAuthMiddleware(m, logger)
	return
}
