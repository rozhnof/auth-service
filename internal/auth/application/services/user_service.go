package services

import (
	"auth/internal/auth/domain/models"
	"auth/internal/auth/domain/repository"
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
)

const (
	atTimeout = time.Hour
	rtTimeout = time.Hour * 72
)

type AccessTokenManager interface {
	NewAccessTokenWithTTL(ttl time.Duration) (string, error)
}

type RefreshTokenManager interface {
	NewRefreshToken() (string, error)
}

type PasswordManager interface {
	HashPassword(password string) (string, error)
	CheckPassword(password string, hashPassword string) bool
}

type Dependencies struct {
	UserRepository    repository.UserRepository
	SessionRepository repository.SessionRepository
	TxManager         repository.TransactionManager
	AtManager         AccessTokenManager
	RtManager         RefreshTokenManager
	PasswordManager   PasswordManager
}

type UserService struct {
	Dependencies
}

func NewUserService(d Dependencies) (*UserService, error) {
	if err := func() error {
		if d.UserRepository == nil {
			return errors.New("user repository")
		}

		if d.SessionRepository == nil {
			return errors.New("session repository")
		}

		if d.TxManager == nil {
			return errors.New("transaction manager")
		}

		if d.AtManager == nil {
			return errors.New("access token manager")
		}

		if d.RtManager == nil {
			return errors.New("refresh token manager")
		}

		if d.PasswordManager == nil {
			return errors.New("password manager")
		}

		return nil
	}(); err != nil {
		return nil, errors.Wrap(err, "mising required dependency")
	}

	return &UserService{
		Dependencies: d,
	}, nil
}

func (s *UserService) Register(ctx context.Context, username string, password string) (*models.User, error) {
	hashPassword, err := s.PasswordManager.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:     username,
		HashPassword: hashPassword,
	}

	createdUser, err := s.UserRepository.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (at string, rt string, err error) {
	if err := s.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		user, err := s.UserRepository.GetByUsername(txCtx, username)
		if err != nil {
			return err
		}

		if !s.PasswordManager.CheckPassword(password, user.HashPassword) {
			return ErrInvalidPassword
		}

		at, err = s.AtManager.NewAccessTokenWithTTL(atTimeout)
		if err != nil {
			return err
		}

		rt, err = s.RtManager.NewRefreshToken()
		if err != nil {
			return err
		}

		session := &models.Session{
			UserID:       user.ID,
			RefreshToken: rt,
			ExpiredAt:    time.Now().Add(rtTimeout),
		}

		slog.Debug("creating a new session", slog.Any("session", session))

		if _, err := s.SessionRepository.Create(txCtx, session); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (at string, rt string, err error) {
	if err := s.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		session, err := s.SessionRepository.GetByRefreshToken(txCtx, refreshToken)
		if err != nil {
			return err
		}

		slog.Debug("refresh session", slog.Any("session", session))

		if !session.Valid() {
			return errors.New("invalid refresh token")
		}

		at, err = s.AtManager.NewAccessTokenWithTTL(atTimeout)
		if err != nil {
			return err
		}

		rt, err = s.RtManager.NewRefreshToken()
		if err != nil {
			return err
		}

		session.RefreshToken = rt
		session.ExpiredAt = time.Now().Add(rtTimeout)

		if _, err := s.SessionRepository.Update(txCtx, session); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", "", err
	}

	return at, rt, nil
}
