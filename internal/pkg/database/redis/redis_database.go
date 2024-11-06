package redisdb

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Address      string
	Port         int
	User         string
	Password     string
	UserPassword string
	DB           int
}

type Redis struct {
	*redis.Client
}

func NewRedis(ctx context.Context, cfg RedisConfig) (*Redis, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		Username: cfg.User,
		Password: cfg.UserPassword,
		DB:       cfg.DB,
	}

	client := redis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	if err := redisotel.InstrumentTracing(client); err != nil {
		return nil, err
	}

	return &Redis{
		Client: client,
	}, nil
}
