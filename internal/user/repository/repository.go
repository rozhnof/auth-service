package repository

import (
	"auth/internal/user/models"
	"auth/internal/user/repository/queries"
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type PostgresDatabase interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type UserRepository struct {
	db PostgresDatabase
}

func NewUserRepository(db PostgresDatabase) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (s *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	userEntity := UserToEntity(user)

	rows, err := s.db.Query(ctx, queries.CreateQuery, userEntity.Username, userEntity.Password)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdUserEntity User
	if err := pgxscan.ScanOne(&createdUserEntity, rows); err != nil {
		return nil, err
	}

	createdUser := UserToModel(&createdUserEntity)

	if *createdUser != *user {
		return nil, errors.Wrapf(ErrDuplicate, "user with email = %s already exists", userEntity.Username)
	}

	return createdUser, nil
}

func (s *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var userEntity User

	if err := pgxscan.Get(ctx, s.db, &userEntity, queries.GetByIDQuery, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(ErrNotExists, fmt.Sprintf("user with id = %s does not exists", userID))
		}
		return nil, err
	}

	return UserToModel(&userEntity), nil
}

func (s *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var userEntity User

	if err := pgxscan.Get(ctx, s.db, &userEntity, queries.GetByUsernameQuery, username); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(ErrNotExists, fmt.Sprintf("user with username = %s does not exists", username))
		}
		return nil, err
	}

	return UserToModel(&userEntity), nil
}

func (s *UserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := s.db.Query(ctx, queries.ListQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userEntityList []User
	if err := pgxscan.ScanAll(&userEntityList, rows); err != nil {
		return nil, err
	}

	return UserListToModel(userEntityList), nil
}

func (s *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	userEntity := UserToEntity(user)

	rows, err := s.db.Query(ctx, queries.CreateQuery, userEntity.Username, userEntity.Password)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updatedUserEntity User
	if err := pgxscan.ScanOne(&updatedUserEntity, rows); err != nil {
		return nil, err
	}

	return UserToModel(&updatedUserEntity), nil
}

func (s *UserRepository) Delete(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	rows, err := s.db.Query(ctx, queries.DeleteQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrNotExists
	}

	var deletedAt time.Time
	if err := rows.Scan(&deletedAt); err != nil {
		return nil, err
	}

	return &deletedAt, nil
}
