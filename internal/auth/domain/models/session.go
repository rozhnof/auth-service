package models

import (
	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken RefreshToken
}

func (s *Session) Valid() bool {
	return s.RefreshToken.Valid()
}
