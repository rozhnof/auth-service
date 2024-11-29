package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/rozhnof/auth-service/internal/domain/entities"
)

type UserFilters struct {
	UserIDs []uuid.UUID
}

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, userID uuid.UUID) error

	GetByID(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*entities.User, error)

	List(ctx context.Context, filters *UserFilters, pagination *Pagination) ([]entities.User, error)
}
