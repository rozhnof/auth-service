package repository_interfaces

import (
	"auth/internal/domain/models"
	"auth/internal/domain/value_object"
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID value_object.UserID) (*models.User, error)
	GetByUsername(ctx context.Context, username value_object.Username) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID value_object.UserID) (*time.Time, error)
}
