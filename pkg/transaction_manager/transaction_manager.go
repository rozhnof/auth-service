package trm

import (
	"context"

	"errors"

	"github.com/avast/retry-go"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	retries = 3
)

type txKeyType string

var txKeyValue = txKeyType("tx")

type Transaction interface {
	Begin(ctx context.Context) (pgx.Tx, error)

	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TransactionManager struct {
	db *pgxpool.Pool
}

func NewTransactionManager(db *pgxpool.Pool) TransactionManager {
	txManager := TransactionManager{
		db: db,
	}

	return txManager
}

func (m TransactionManager) WithTransaction(ctx context.Context, f func(context.Context) error) error {
	txOptions := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	if _, ok := ctx.Value(txKeyValue).(Transaction); ok {
		return f(ctx)
	}

	if err := retry.Do(func() error {
		tx, err := m.db.BeginTx(ctx, txOptions)
		if err != nil {
			return err
		}

		ctxWithTx := context.WithValue(ctx, txKeyValue, tx)
		if err := f(ctxWithTx); err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				return errors.Join(err, errRollback)
			}

			return err
		}

		if err := tx.Commit(ctx); err != nil {
			return err
		}

		return nil
	}, retry.Attempts(retries), retry.RetryIf(isSerializationError)); err != nil {
		return err
	}

	return nil
}

func (m TransactionManager) TxOrDB(ctx context.Context) Transaction {
	tx, ok := ctx.Value(txKeyValue).(Transaction)
	if !ok {
		return m.db
	}

	return tx
}

func isSerializationError(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.SerializationFailure {
		return true
	}

	return false
}
