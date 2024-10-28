package postgres_database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Address  string
	Port     int
	User     string
	Password string
	DB       string
	SSL      string
}

type Database struct {
	cluster *pgxpool.Pool
}

func (db Database) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return db.cluster.BeginTx(ctx, txOptions)
}

func (db Database) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.cluster.Query(ctx, sql, args...)
}

func (db Database) Close() {
	db.cluster.Close()
}

func NewDatabase(ctx context.Context, cfg Config) (*Database, error) {
	cluster, err := pgxpool.Connect(ctx, CreateConnectionString(cfg))
	if err != nil {
		return nil, err
	}
	return &Database{cluster: cluster}, nil
}

func CreateConnectionString(cfg Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Address, cfg.Port, cfg.DB, cfg.SSL)
}
