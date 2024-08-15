package postgres_repository

import (
	repository "auth/internal/infrastructure"
	"auth/internal/infrastructure/repository/postgres/dao"
	"auth/internal/infrastructure/repository/postgres/queries/auth_queries"
	"auth/pkg/database"
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const txKey = "transaction"

type AuthRepository struct {
	db *database.Database
}

func NewAuthRepository(db *database.Database) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (s *AuthRepository) Create(ctx context.Context, user dao.User) (*uuid.UUID, error) {
	querier := s.getQuerier(ctx)

	rows, err := querier.Query(ctx, auth_queries.CreateQuery, uuid.New(), user.Email, user.TokenHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userID uuid.UUID
	if err := pgxscan.ScanOne(&userID, rows); err != nil {
		return nil, err
	}

	return &userID, nil
}

func (s *AuthRepository) GetByID(ctx context.Context, userID uuid.UUID) (*dao.User, error) {
	querier := s.getQuerier(ctx)

	var user dao.User
	err := pgxscan.Get(ctx, querier, &user, auth_queries.GetByIDQuery, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNoData, fmt.Sprintf("user with id = %s does not exists", userID))
		}
		return nil, err
	}

	return &user, nil
}

func (s *AuthRepository) List(ctx context.Context) ([]dao.User, error) {
	querier := s.getQuerier(ctx)

	rows, err := querier.Query(ctx, auth_queries.ListQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList []dao.User
	if err := pgxscan.ScanAll(&userList, rows); err != nil {
		return nil, err
	}

	return userList, nil
}

func (s *AuthRepository) Update(ctx context.Context, user dao.User) error {
	querier := s.getQuerier(ctx)

	rows, err := querier.Query(ctx, auth_queries.UpdateQuery, user.ID, user.Email, user.TokenHash)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (s *AuthRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	querier := s.getQuerier(ctx)

	rows, err := querier.Query(ctx, auth_queries.DeleteQuery, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (s *AuthRepository) getQuerier(ctx context.Context) pgxscan.Querier {
	querier, ok := ctx.Value(txKey).(pgxscan.Querier)
	if !ok {
		return s.db.GetPool()
	}
	return querier
}

func (s *AuthRepository) WithTransaction(ctx context.Context, f func(ctxTX context.Context) error) error {
	txOptions := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	pool := s.db.GetPool()
	tx, err := pool.BeginTx(ctx, txOptions)

	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)
	if err = f(ctxWithTx); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			errors.Wrap(err, errRollback.Error())
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
