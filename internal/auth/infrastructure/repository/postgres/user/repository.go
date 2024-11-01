package postgres_user_repository

import (
	"auth/internal/auth/domain/models"
	"auth/internal/auth/infrastructure/repository"
	queries "auth/internal/auth/infrastructure/repository/postgres/user/queries"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type PostgresDatabase interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type UserRepository struct {
	db  PostgresDatabase
	log *slog.Logger
}

func NewUserRepository(db PostgresDatabase, log *slog.Logger) *UserRepository {
	log = log.With(
		slog.String("layer", "infrastructure"),
		slog.String("pkg", "postgres_user_repository"),
	)

	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (s *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	log := s.log.With(
		slog.String("function", "UserRepository.Create"),
		slog.String("username", user.Username),
	)

	log.Debug("create user start", slog.String("password", user.HashPassword))

	rows, err := s.db.Query(ctx, queries.Create, user.Username, user.HashPassword)
	if err != nil {
		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}
	defer rows.Close()

	var createdUser models.User
	if err := pgxscan.ScanOne(&createdUser, rows); err != nil {
		var pgErr *pgconn.PgError

		log.Info("scan rows error", slog.Any("error", err.Error()))

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, errors.Wrapf(repository.ErrDuplicate, "user with username = %s already exists", user.Username)
		}

		return nil, err
	}

	log.Debug("create user end")

	return &createdUser, nil
}

func (s *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	log := s.log.With(
		slog.String("function", "UserRepository.GetByID"),
		slog.String("user_id", userID.String()),
	)

	log.Debug("get user by id start")

	var userEntity models.User

	if err := pgxscan.Get(ctx, s.db, &userEntity, queries.GetByID, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNotExists, fmt.Sprintf("user with id = %s does not exists", userID))
		}

		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("get user by id end")

	return &userEntity, nil
}

func (s *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	log := s.log.With(
		slog.String("function", "UserRepository.GetByUsername"),
		slog.String("username", username),
	)

	log.Debug("get user by username start")

	var userEntity models.User

	if err := pgxscan.Get(ctx, s.db, &userEntity, queries.GetByUsername, username); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNotExists, fmt.Sprintf("user with username = %s does not exists", username))
		}

		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("get user by username end")

	return &userEntity, nil
}

func (s *UserRepository) List(ctx context.Context) ([]models.User, error) {
	log := s.log.With(
		slog.String("function", "UserRepository.List"),
	)

	log.Debug("list user start")

	rows, err := s.db.Query(ctx, queries.List)
	if err != nil {
		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}
	defer rows.Close()

	var userList []models.User
	if err := pgxscan.ScanAll(&userList, rows); err != nil {
		log.Info("scan rows error", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("list user end")

	return userList, nil
}

func (s *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	log := s.log.With(
		slog.String("function", "UserRepository.Update"),
		slog.String("username", user.Username),
	)

	log.Debug("update user start")

	rows, err := s.db.Query(ctx, queries.Create, user.Username, user.HashPassword)
	if err != nil {
		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}
	defer rows.Close()

	var updatedUser models.User
	if err := pgxscan.ScanOne(&updatedUser, rows); err != nil {
		log.Info("scan rows error", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("update user end")

	return &updatedUser, nil
}

func (s *UserRepository) Delete(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	log := s.log.With(
		slog.String("function", "UserRepository.Create"),
		slog.String("user_id", userID.String()),
	)

	log.Debug("delete user start")

	rows, err := s.db.Query(ctx, queries.Delete, userID)
	if err != nil {
		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, repository.ErrNotExists
	}

	var deletedAt time.Time
	if err := rows.Scan(&deletedAt); err != nil {
		log.Info("scan rows error", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("delete user end")

	return &deletedAt, nil
}
