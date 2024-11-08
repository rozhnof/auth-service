package services

import (
	"auth/internal/auth/domain/models"
	"auth/internal/auth/domain/repository"
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

type AccessTokenManager interface {
	NewToken() (models.AccessToken, error)
}

type RefreshTokenManager interface {
	NewToken() (models.RefreshToken, error)
}

type PasswordManager interface {
	HashPassword(password string) (string, error)
	CheckPassword(password string, hashPassword string) bool
}

const ttl = time.Hour * 60

type Dependencies struct {
	UserRepository    repository.UserRepository
	SessionRepository repository.SessionRepository
	TxManager         repository.TransactionManager
	AtManager         AccessTokenManager
	RtManager         RefreshTokenManager
	PasswordManager   PasswordManager
	UserCache         repository.UserCache
	SessionCache      repository.SessionCache
}

func (d Dependencies) Valid() error {
	if d.UserRepository == nil {
		return errors.New("missing user repository")
	}

	if d.SessionRepository == nil {
		return errors.New("missing session repository")
	}

	if d.TxManager == nil {
		return errors.New("missing transaction manager")
	}

	if d.AtManager == nil {
		return errors.New("missing access token manager")
	}

	if d.RtManager == nil {
		return errors.New("missing refresh token manager")
	}

	if d.PasswordManager == nil {
		return errors.New("missing password manager")
	}

	if d.UserCache == nil {
		return errors.New("missing user cache")
	}

	if d.SessionCache == nil {
		return errors.New("missing session cache")
	}

	return nil
}

type UserService struct {
	Dependencies
	log    *slog.Logger
	tracer trace.Tracer
}

func NewUserService(d Dependencies, log *slog.Logger, tracer trace.Tracer) (*UserService, error) {
	if err := d.Valid(); err != nil {
		return nil, errors.Wrap(err, "missing required dependency")
	}

	return &UserService{
		Dependencies: d,
		log:          log,
		tracer:       tracer,
	}, nil
}

func (s *UserService) Register(ctx context.Context, username string, password string) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserService.Register")
	defer span.End()

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

	s.UserCache.Set(ctx, username, *createdUser, ttl)

	return createdUser, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (at string, rt string, err error) {
	ctx, span := s.tracer.Start(ctx, "UserService.Login")
	defer span.End()

	getUser := func(username string) (*models.User, error) {
		if cacheUser, err := s.UserCache.Get(ctx, username); err == nil {
			return &cacheUser, nil
		}

		dbUser, err := s.UserRepository.GetByUsername(ctx, username)
		if err != nil {
			return nil, err
		}

		return dbUser, nil
	}

	user, err := getUser(username)
	if err != nil {
		return "", "", err
	}

	if !s.PasswordManager.CheckPassword(password, user.HashPassword) {
		return "", "", ErrInvalidPassword
	}

	accessToken, err := s.AtManager.NewToken()
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.RtManager.NewToken()
	if err != nil {
		return "", "", err
	}

	session := &models.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	}

	if _, err := s.SessionRepository.Create(ctx, session); err != nil {
		return "", "", err
	}

	at = accessToken.Token
	rt = refreshToken.Token

	s.UserCache.Set(ctx, username, *user, ttl)

	return at, rt, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (at string, rt string, err error) {
	ctx, span := s.tracer.Start(ctx, "UserService.Refresh")
	defer span.End()

	session, err := s.SessionRepository.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if !session.Valid() {
		return "", "", ErrUnauthorizedRefresh
	}

	newAccessToken, err := s.AtManager.NewToken()
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.RtManager.NewToken()
	if err != nil {
		return "", "", err
	}

	// session.RefreshToken = newRefreshToken

	if _, err := s.SessionRepository.Update(ctx, session); err != nil {
		return "", "", err
	}

	at = newAccessToken.Token
	rt = newRefreshToken.Token

	return at, rt, nil
}
