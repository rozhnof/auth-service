package postgres_user_repository

import (
	"auth/internal/domain/models"
	"auth/internal/domain/repository/filters"
	repository "auth/internal/infrastructure"
	postgres_user_dto "auth/internal/infrastructure/repository/postgres/user/dto"
	postgres_user_queries "auth/internal/infrastructure/repository/postgres/user/queries"
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
	userEntity := postgres_user_dto.UserToEntity(user)

	rows, err := s.db.Query(ctx, postgres_user_queries.CreateQuery, userEntity.Username, userEntity.Password)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdUserEntity postgres_user_dto.User
	if err := pgxscan.ScanOne(&createdUserEntity, rows); err != nil {
		return nil, err
	}

	createdUser := postgres_user_dto.UserToModel(&createdUserEntity)

	if *createdUser != *user {
		return nil, errors.Wrapf(repository.ErrDuplicate, "user with email = %s already exists", userEntity.Username)
	}

	return createdUser, nil
}

func (s *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var userEntity postgres_user_dto.User

	if err := pgxscan.Get(ctx, s.db, &userEntity, postgres_user_queries.GetByIDQuery, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNotExists, fmt.Sprintf("user with id = %s does not exists", userID))
		}
		return nil, err
	}

	return postgres_user_dto.UserToModel(&userEntity), nil
}

func (s *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var userEntity postgres_user_dto.User

	if err := pgxscan.Get(ctx, s.db, &userEntity, postgres_user_queries.GetByUsernameQuery, username); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNotExists, fmt.Sprintf("user with username = %s does not exists", username))
		}
		return nil, err
	}

	return postgres_user_dto.UserToModel(&userEntity), nil
}

func (s *UserRepository) List(ctx context.Context, filter *filters.UserFilter, pagination *filters.Pagination) ([]models.User, error) {
	var (
		query = postgres_user_queries.ListQuery
		args  []any
	)

	if filter != nil {
		query, args = addFilters(query, args, *filter)
	}

	if pagination != nil {
		query, args = addPagination(query, args, *pagination)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userEntityList []postgres_user_dto.User
	if err := pgxscan.ScanAll(&userEntityList, rows); err != nil {
		return nil, err
	}

	return postgres_user_dto.UserListToModel(userEntityList), nil
}

func (s *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	userEntity := postgres_user_dto.UserToEntity(user)

	rows, err := s.db.Query(ctx, postgres_user_queries.CreateQuery, userEntity.Username, userEntity.Password)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updatedUserEntity postgres_user_dto.User
	if err := pgxscan.ScanOne(&updatedUserEntity, rows); err != nil {
		return nil, err
	}

	return postgres_user_dto.UserToModel(&updatedUserEntity), nil
}

func (s *UserRepository) Delete(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	rows, err := s.db.Query(ctx, postgres_user_queries.DeleteQuery, userID)
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

func addFilters(query string, args []any, filter filters.UserFilter) (string, []any) {
	if len(filter.Email) > 0 {
		args, query = addAnyOfCondition(args, query, filter.Email, "email")
	}

	return query, args
}

func addPagination(query string, args []any, pagination filters.Pagination) (string, []any) {
	if pagination.Limit > 0 {
		args = append(args, pagination.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
	}

	if pagination.Offset > 0 {
		args = append(args, pagination.Offset)
		query += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	return query, args
}

func addAnyOfCondition(args []any, query string, filter any, name string) ([]any, string) {
	args = append(args, filter)
	query += fmt.Sprintf(" AND %s = ANY($%d)", name, len(args))
	return args, query
}

func addGreaterEqualCondition(args []any, query string, filter any, name string) ([]any, string) {
	args = append(args, filter)
	query += fmt.Sprintf(" AND %s >= $%d", name, len(args))
	return args, query
}

func addLessEqualCondition(args []any, query string, filter any, name string) ([]any, string) {
	args = append(args, filter)
	query += fmt.Sprintf(" AND %s <= $%d", name, len(args))
	return args, query
}
