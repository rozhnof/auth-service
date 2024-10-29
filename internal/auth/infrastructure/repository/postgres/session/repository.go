package postgres_session_repository

import (
	"auth/internal/auth/domain/models"
	"auth/internal/auth/infrastructure/repository"
	queries "auth/internal/auth/infrastructure/repository/postgres/session/queries"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type PostgresDatabase interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type SessionRepository struct {
	db  PostgresDatabase
	log *slog.Logger
}

func NewSessionRepository(db PostgresDatabase, log *slog.Logger) *SessionRepository {
	log = log.With(
		slog.String("layer", "infrastructure"),
		slog.String("pkg", "postgres_session_repository"),
	)

	return &SessionRepository{
		db:  db,
		log: log,
	}
}

func (s *SessionRepository) Create(ctx context.Context, session *models.Session) (*models.Session, error) {
	log := s.log.With(
		slog.String("function", "SessionRepository.Create"),
		slog.String("user_id", session.UserID.String()),
	)

	log.Debug("create session start")

	sessionEntity := SessionToEntity(session)

	args := []any{
		sessionEntity.UserID,
		sessionEntity.RefreshToken,
		sessionEntity.ExpiredAt,
		sessionEntity.IsRevoked,
	}

	rows, err := s.db.Query(ctx, queries.Create, args...)
	if err != nil {
		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}
	defer rows.Close()

	var createdSessionEntity Session
	if err := pgxscan.ScanOne(&createdSessionEntity, rows); err != nil {
		log.Info("scan rows error", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("create user end")

	return SessionToModel(&createdSessionEntity), nil
}

func (s *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	log := s.log.With(
		slog.String("function", "SessionRepository.GetByRefreshToken"),
	)

	log.Debug("get session by refresh token start")

	var session Session

	if err := pgxscan.Get(ctx, s.db, &session, queries.GetByRefreshToken, refreshToken); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNotExists, fmt.Sprintf("session with refresh token = %s does not exists", refreshToken))
		}

		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("get session by refresh token end")

	return SessionToModel(&session), nil
}

func (s *SessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	log := s.log.With(
		slog.String("function", "SessionRepository.GetByID"),
		slog.String("session_id", sessionID.String()),
	)

	log.Debug("get session by id end")

	var session Session

	if err := pgxscan.Get(ctx, s.db, &session, queries.GetByID, sessionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNotExists, fmt.Sprintf("session with id = %s does not exists", sessionID))
		}

		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("get session by id end")

	return SessionToModel(&session), nil
}

func (s *SessionRepository) Update(ctx context.Context, session *models.Session) (*models.Session, error) {
	log := s.log.With(
		slog.String("function", "SessionRepository.Update"),
		slog.String("session_id", session.ID.String()),
	)

	log.Debug("update session start")

	sessionEntity := SessionToEntity(session)

	args := []any{
		sessionEntity.ID,
		sessionEntity.UserID,
		sessionEntity.RefreshToken,
		sessionEntity.ExpiredAt,
		sessionEntity.IsRevoked,
	}

	rows, err := s.db.Query(ctx, queries.Update, args...)
	if err != nil {
		log.Info("failed to execute postgres query", slog.Any("error", err.Error()))

		return nil, err
	}
	defer rows.Close()

	var updatedSession Session
	if err := pgxscan.ScanOne(&updatedSession, rows); err != nil {
		log.Info("scan rows error", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("update session end")

	return SessionToModel(&updatedSession), nil
}

func (s *SessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) (*time.Time, error) {
	log := s.log.With(
		slog.String("function", "UserRepository.Create"),
		slog.String("session_id", sessionID.String()),
	)

	log.Debug("delete session start")

	rows, err := s.db.Query(ctx, queries.Delete, sessionID)
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

	log.Debug("delete session end")

	return &deletedAt, nil
}
