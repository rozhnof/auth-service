package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiredAt    time.Time `db:"expired_at"`
	IsRevoked    bool      `db:"is_revoked"`
}

func (s Session) Valid() bool {
	return s.ExpiredAt.Before(time.Now()) && !s.IsRevoked
}
