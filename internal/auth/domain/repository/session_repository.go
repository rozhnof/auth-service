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
	Update(ctx context.Context, session *models.Session) (*models.Session, error)
	Delete(ctx context.Context, sessionID uuid.UUID) (*time.Time, error)
	RevokeByUserID(ctx context.Context, userID uuid.UUID) error
}

type SessionCache interface {
	Get(ctx context.Context, key string) (models.Session, error)
	Set(ctx context.Context, key string, value models.Session, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
