package postgres_session_repository

import (
	"auth/internal/auth/domain/models"
	"auth/internal/auth/infrastructure/repository"
	queries "auth/internal/auth/infrastructure/repository/postgres/session/queries"
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

type SessionRepository struct {
	db PostgresDatabase
}

func NewSessionRepository(db PostgresDatabase) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (s *SessionRepository) Create(ctx context.Context, session *models.Session) (*models.Session, error) {
	rows, err := s.db.Query(ctx, queries.Create)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdSession models.Session
	if err := pgxscan.ScanOne(&createdSession, rows); err != nil {
		return nil, err
	}

	return &createdSession, nil
}

func (s *SessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session

	if err := pgxscan.Get(ctx, s.db, &session, queries.GetByID, sessionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(repository.ErrNotExists, fmt.Sprintf("session with id = %s does not exists", sessionID))
		}
		return nil, err
	}

	return &session, nil
}

func (s *SessionRepository) Update(ctx context.Context, session *models.Session) (*models.Session, error) {
	rows, err := s.db.Query(ctx, queries.Create)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updatedSession models.Session
	if err := pgxscan.ScanOne(&updatedSession, rows); err != nil {
		return nil, err
	}

	return &updatedSession, nil
}

func (s *SessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) (*time.Time, error) {
	rows, err := s.db.Query(ctx, queries.Delete, sessionID)
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
