package pgrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/rozhnof/auth-service/internal/domain/entities"
	vobjects "github.com/rozhnof/auth-service/internal/domain/value_objects"
	db_queries "github.com/rozhnof/auth-service/internal/infrastructure/repository/queries"
)

type userRow = struct {
	User                  db_queries.User
	RefreshTokenID        *uuid.UUID
	RefreshToken          *string
	RefreshTokenExpiredAt *time.Time
}

func userRowToUser(row userRow) *entities.User {
	var refreshToken *vobjects.RefreshToken

	if row.RefreshTokenID != nil {
		rt := vobjects.NewExistingRefreshToken(*row.RefreshToken, *row.RefreshTokenExpiredAt)
		refreshToken = &rt
	}

	return entities.NewExistingUser(
		row.User.ID,
		row.User.Email,
		vobjects.NewExistingPassword(row.User.HashPassword),
		refreshToken,
	)
}
