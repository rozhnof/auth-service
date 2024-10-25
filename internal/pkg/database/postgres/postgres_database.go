package postgres_database

import (
	"auth/internal/pkg/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	cluster *pgxpool.Pool
}

func (db Database) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.cluster.Query(ctx, sql, args...)
}

func (db Database) Close() {
	db.cluster.Close()
}

func NewDatabase(ctx context.Context, connString string) (*Database, error) {
	cluster, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	return &Database{cluster: cluster}, nil
}

func CreateConnectionString(cfg config.RepositoryConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Address, cfg.Port, cfg.DB, cfg.SSL)
}
