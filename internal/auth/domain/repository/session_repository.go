package repository

import (
	"auth/internal/auth/domain/models"
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) (*models.Session, error)
	GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	List(ctx context.Context) ([]models.Session, error)
	Update(ctx context.Context, session *models.Session) (*models.Session, error)
	Delete(ctx context.Context, sessionID uuid.UUID) (*time.Time, error)
}
