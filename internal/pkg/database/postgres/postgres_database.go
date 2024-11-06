package pgxdb

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Address  string
	Port     int
	User     string
	Password string
	DB       string
	SSL      string
}

func NewDatabase(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(CreateConnectionString(cfg))
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	pgxCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return pool, nil
}

func CreateConnectionString(cfg Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Address, cfg.Port, cfg.DB, cfg.SSL)
}
