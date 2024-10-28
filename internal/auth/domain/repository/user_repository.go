package repository

import (
	"auth/internal/auth/domain/models"
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) (*time.Time, error)
}
