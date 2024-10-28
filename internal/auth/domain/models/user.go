package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	HashPassword string    `db:"hash_password"`
}
