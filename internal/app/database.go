package app

import (
	"context"

	"github.com/rozhnof/auth-service/internal/infrastructure/database/postgres"
	"github.com/rozhnof/auth-service/internal/infrastructure/database/redis"
	"github.com/rozhnof/auth-service/internal/pkg/config"
)

func NewPostgresDatabase(ctx context.Context, cfg config.Postgres) (postgres.Database, error) {
	postgresConfig := postgres.DatabaseConfig{
		Address:  cfg.Address,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
		SSL:      cfg.SSL,
	}

	return postgres.NewDatabase(ctx, postgresConfig)
}

func NewRedisDatabase(ctx context.Context, cfg config.Redis) (redis.Database, error) {
	redisConfig := redis.DatabaseConfig{
		Address:  cfg.Address,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	return redis.NewDatabase(ctx, redisConfig)
}
