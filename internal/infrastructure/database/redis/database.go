package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type DatabaseConfig struct {
	Address      string
	Port         int
	User         string
	Password     string
	UserPassword string
	DB           int
}

type Database struct {
	*redis.Client
}

func NewDatabase(ctx context.Context, cfg DatabaseConfig) (Database, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		Username: cfg.User,
		Password: cfg.UserPassword,
		DB:       cfg.DB,
	}

	client := redis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		return Database{}, err
	}

	if err := redisotel.InstrumentTracing(client); err != nil {
		return Database{}, err
	}

	return Database{
		Client: client,
	}, nil
}
