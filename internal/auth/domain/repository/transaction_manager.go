package repository

import (
	"context"
)

type TransactionManager interface {
	WithTransaction(ctx context.Context, f func(ctx context.Context) error) error
}
