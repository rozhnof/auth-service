package services

import (
	"auth/internal/auth/domain/models"
	"auth/internal/auth/domain/repository"
	"context"
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
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	txManager       repository.TransactionManager
	atManager       AccessTokenManager
	rtManager       RefreshTokenManager
	passwordManager PasswordManager
}

type UserService struct {
	Dependencies
}

func NewAuthService(d Dependencies) (*UserService, error) {
	if err := func() error {
		if d.userRepository == nil {
			return errors.New("user repository")
		}

		if d.sessionRepository == nil {
			return errors.New("session repository")
		}

		if d.txManager == nil {
			return errors.New("transaction manager")
		}

		if d.atManager == nil {
			return errors.New("access token manager")
		}

		if d.rtManager == nil {
			return errors.New("refresh token manager")
		}

		if d.passwordManager == nil {
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
	hashPassword, err := s.passwordManager.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:     username,
		HashPassword: hashPassword,
	}

	createdUser, err := s.userRepository.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (at string, rt string, err error) {
	if err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		user, err := s.userRepository.GetByUsername(txCtx, username)
		if err != nil {
			return err
		}

		if !s.passwordManager.CheckPassword(password, user.HashPassword) {
			return ErrInvalidPassword
		}

		at, err = s.atManager.NewAccessTokenWithTTL(atTimeout)
		if err != nil {
			return err
		}

		rt, err = s.rtManager.NewRefreshToken()
		if err != nil {
			return err
		}

		session := &models.Session{
			UserID:       user.ID,
			RefreshToken: rt,
			ExpiredAt:    time.Now().Add(rtTimeout),
		}

		if _, err := s.sessionRepository.Create(txCtx, session); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (at string, rt string, err error) {
	if err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		session, err := s.sessionRepository.GetByRefreshToken(txCtx, refreshToken)
		if err != nil {
			return err
		}

		if !session.Valid() {
			return errors.New("invalid refresh token")
		}

		at, err = s.atManager.NewAccessTokenWithTTL(atTimeout)
		if err != nil {
			return err
		}

		rt, err = s.rtManager.NewRefreshToken()
		if err != nil {
			return err
		}

		session.RefreshToken = rt
		session.ExpiredAt = time.Now().Add(rtTimeout)

		if _, err := s.sessionRepository.Update(txCtx, session); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", "", err
	}

	return at, rt, nil
}
