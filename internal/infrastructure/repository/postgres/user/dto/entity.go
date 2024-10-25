package postgres_user_dto

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	Password     string    `db:"password"`
	RefreshToken string    `db:"refresh_token"`
}
