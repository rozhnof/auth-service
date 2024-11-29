package postgres

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfig struct {
	Address  string
	Port     int
	User     string
	Password string
	DB       string
	SSL      string
}

type Database struct {
	*pgxpool.Pool
}

func NewDatabase(ctx context.Context, cfg DatabaseConfig) (Database, error) {
	pgxCfg, err := pgxpool.ParseConfig(CreateConnectionString(cfg))
	if err != nil {
		return Database{}, fmt.Errorf("create connection pool: %w", err)
	}

	pgxCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return Database{}, fmt.Errorf("connect to database: %w", err)
	}

	return Database{
		Pool: pool,
	}, nil
}

func CreateConnectionString(cfg DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Address, cfg.Port, cfg.DB, cfg.SSL)
}
