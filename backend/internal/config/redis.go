package config

// todo возможно надо убрать
import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisConfig RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host,
		Username: redisConfig.User,
		Password: redisConfig.Password,
		DB:       0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	return rdb, err
}
