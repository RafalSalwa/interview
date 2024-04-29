package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

func NewUniversalRedisClient(ctx context.Context, cfg *Config) (*redis.Client, error) {
	universalClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	if err := universalClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return universalClient, nil
}
