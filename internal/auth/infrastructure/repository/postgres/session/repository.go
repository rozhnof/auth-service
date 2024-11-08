package postgres_session_repository

import (
	"auth/internal/auth/domain/models"
	queries "auth/internal/auth/infrastructure/repository/postgres/session/queries"
	pgxdb "auth/internal/pkg/database/postgres"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

type PostgresDatabase interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type SessionRepository struct {
	db        PostgresDatabase
	txManager *pgxdb.TransactionManager
	log       *slog.Logger
	tracer    trace.Tracer
}

func NewSessionRepository(db PostgresDatabase, txManager *pgxdb.TransactionManager, log *slog.Logger, tracer trace.Tracer) *SessionRepository {
	return &SessionRepository{
		db:        db,
		txManager: txManager,
		log:       log,
		tracer:    tracer,
	}
}

func (s *SessionRepository) Create(ctx context.Context, session *models.Session) (*models.Session, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.Create")
	defer span.End()

	sessionEntity := SessionToEntity(session)

	args := []any{
		sessionEntity.UserID,
		sessionEntity.RefreshToken,
		sessionEntity.ExpiredAt,
		sessionEntity.IsRevoked,
	}

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.Create, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	createdSessionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Session])
	if err != nil {
		return nil, err
	}

	return SessionToModel(&createdSessionEntity), nil
}

func (s *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.GetByRefreshToken")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.GetByRefreshToken, refreshToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Session])
	if err != nil {
		return nil, err
	}

	return SessionToModel(&sessionEntity), nil
}

func (s *SessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.GetByID")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.GetByID, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Session])
	if err != nil {
		return nil, err
	}

	return SessionToModel(&sessionEntity), nil
}

func (s *SessionRepository) Update(ctx context.Context, session *models.Session) (*models.Session, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.Update")
	defer span.End()

	sessionEntity := SessionToEntity(session)

	args := []any{
		sessionEntity.ID,
		sessionEntity.UserID,
		sessionEntity.RefreshToken,
		sessionEntity.ExpiredAt,
		sessionEntity.IsRevoked,
	}

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.Update, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updatedSessionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Session])
	if err != nil {
		return nil, err
	}

	return SessionToModel(&updatedSessionEntity), nil
}

func (s *SessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) (*time.Time, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.Delete")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, queries.Delete, sessionID)
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
