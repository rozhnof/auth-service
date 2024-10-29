package postgres_session_repository

import (
	"auth/internal/auth/domain/models"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiredAt    time.Time `db:"expired_at"`
	IsRevoked    bool      `db:"is_revoked"`
}

func SessionToModel(session *Session) *models.Session {
	return &models.Session{
		ID:     session.ID,
		UserID: session.UserID,
		RefreshToken: models.RefreshToken{
			Token:     session.RefreshToken,
			ExpiredAt: session.ExpiredAt,
			IsRevoked: session.IsRevoked,
		},
	}
}

func SessionListToModel(sessionEntityList []Session) []models.Session {
	sessionList := make([]models.Session, 0, len(sessionEntityList))
	for _, sessionEntity := range sessionEntityList {
		sessionList = append(sessionList, *SessionToModel(&sessionEntity))
	}

	return sessionList
}

func SessionToEntity(session *models.Session) *Session {
	return &Session{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken.Token,
		ExpiredAt:    session.RefreshToken.ExpiredAt,
		IsRevoked:    session.RefreshToken.IsRevoked,
	}
}

func SessionListToEntity(sessionList []models.Session) []Session {
	sessionEntityList := make([]Session, 0, len(sessionList))
	for _, session := range sessionList {
		sessionEntityList = append(sessionEntityList, *SessionToEntity(&session))
	}

	return sessionEntityList
}
