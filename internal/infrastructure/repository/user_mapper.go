package pgrepo

import (
	"github.com/rozhnof/auth-service/internal/domain/entities"
	vobjects "github.com/rozhnof/auth-service/internal/domain/value_objects"
	db_queries "github.com/rozhnof/auth-service/internal/infrastructure/repository/queries"
)

func dtoToRefreshToken(rt *db_queries.RefreshToken) *vobjects.RefreshToken {
	if rt == nil {
		return nil
	}

	return vobjects.NewExistingRefreshToken(rt.Token, rt.ExpiredAt)
}

func dtoToRegisterToken(rt *db_queries.RegisterToken) *vobjects.RegisterToken {
	if rt == nil {
		return nil
	}

	return vobjects.NewExistingRegisterToken(rt.Token, rt.ExpiredAt)
}

func dtoToUser(user db_queries.User, refreshToken *db_queries.RefreshToken, registerToken *db_queries.RegisterToken) *entities.User {
	return entities.NewExistingUser(
		user.ID,
		user.Email,
		vobjects.NewExistingPassword(user.HashPassword),
		user.Confirmed,
		dtoToRefreshToken(refreshToken),
		dtoToRegisterToken(registerToken),
	)
}
