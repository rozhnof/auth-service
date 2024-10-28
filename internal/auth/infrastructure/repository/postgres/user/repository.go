package postgres_user_repository

import (
	"auth/internal/auth/domain/models"
	"auth/internal/auth/infrastructure/repository"
	queries "auth/internal/auth/infrastructure/repository/postgres/user/queries"
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
	rows, err := s.db.Query(ctx, queries.Create, user.Username, user.HashPassword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdUser models.User
	if err := pgxscan.ScanOne(&createdUser, rows); err != nil {
		return nil, err
	}

	user.ID = createdUser.ID
	if createdUser != *user {
		return nil, errors.Wrapf(repository.ErrDuplicate, "user with email = %s already exists", user.Username)
	}

	return &createdUser, nil
}

func (s *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var userEntity models.User

	if err := pgxscan.Get(ctx, s.db, &userEntity, queries.GetByID, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNotExists, fmt.Sprintf("user with id = %s does not exists", userID))
		}
		return nil, err
	}

	return &userEntity, nil
}

func (s *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var userEntity models.User

	if err := pgxscan.Get(ctx, s.db, &userEntity, queries.GetByUsername, username); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNotExists, fmt.Sprintf("user with username = %s does not exists", username))
		}
		return nil, err
	}

	return &userEntity, nil
}

func (s *UserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := s.db.Query(ctx, queries.List)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList []models.User
	if err := pgxscan.ScanAll(&userList, rows); err != nil {
		return nil, err
	}

	return userList, nil
}

func (s *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	rows, err := s.db.Query(ctx, queries.Create, user.Username, user.HashPassword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updatedUser models.User
	if err := pgxscan.ScanOne(&updatedUser, rows); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (s *UserRepository) Delete(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	rows, err := s.db.Query(ctx, queries.Delete, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, repository.ErrNotExists
	}

	var deletedAt time.Time
	if err := rows.Scan(&deletedAt); err != nil {
		return nil, err
	}

	return &deletedAt, nil
}
