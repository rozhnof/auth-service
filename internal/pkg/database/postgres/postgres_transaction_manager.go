package postgres_database

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type txKeyType string

var txKeyValue = txKeyType("tx")

type QueryEngine interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type TransactionManager struct {
	*Database
}

func NewTransactionManager(db *Database) *TransactionManager {
	txManager := &TransactionManager{
		Database: db,
	}

	return txManager
}

func (m *TransactionManager) WithTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	txOptions := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	tx, err := m.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, txKeyValue, tx)
	if err := f(ctxWithTx); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			// TODO
			return errors.Errorf("%w and %w", err, errRollback)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (m *TransactionManager) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	queryEngine, ok := ctx.Value(txKeyValue).(QueryEngine)
	if !ok {
		return m.Database.Query(ctx, sql, args...)
	}

	return queryEngine.Query(ctx, sql, args...)
}
