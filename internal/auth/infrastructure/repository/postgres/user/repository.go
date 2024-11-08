package postgres_user_repository

import (
	"auth/internal/auth/domain/models"
	queries "auth/internal/auth/infrastructure/repository/postgres/user/queries"
	pgxdb "auth/internal/pkg/database/postgres"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

type UserRepository struct {
	db        *pgxpool.Pool
	txManager *pgxdb.TransactionManager
	log       *slog.Logger
	tracer    trace.Tracer
}

func NewUserRepository(db *pgxpool.Pool, txManager *pgxdb.TransactionManager, log *slog.Logger, tracer trace.Tracer) *UserRepository {
	return &UserRepository{
		db:        db,
		txManager: txManager,
		log:       log,
		tracer:    tracer,
	}
}

func (s *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Create")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.Create, user.Username, user.HashPassword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	createdUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (s *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByID")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.GetByID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByUsername")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.GetByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserRepository) List(ctx context.Context) ([]models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.List")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.List)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userList, err := pgx.CollectRows(rows, pgx.RowTo[models.User])
	if err != nil {
		return nil, err
	}

	return userList, nil
}

func (s *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Update")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.Update, user.Username, user.HashPassword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updatedUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (s *UserRepository) Delete(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Delete")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.Delete, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deletedAt, err := pgx.RowTo[time.Time](rows)
	if err != nil {
		return nil, err
	}

	return &deletedAt, nil
}
