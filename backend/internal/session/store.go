package session

import (
	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
)

func NewScsRedisStore(redisClient *redis.Client) scs.Store {
	return goredisstore.New(redisClient)
}
