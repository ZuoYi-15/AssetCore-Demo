package redis

import (
	"context"

	"asset-core/internal/config"

	goredis "github.com/go-redis/redis/v8"
)

func New(cfg config.RedisConfig) *goredis.Client {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_ = client.Ping(context.Background()).Err()
	return client
}
