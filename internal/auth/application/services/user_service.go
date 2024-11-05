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

const ttl = time.Second * 60

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

	log = log.With(
		slog.String("layer", "application"),
		slog.String("pkg", "services"),
	)

	return &UserService{
		Dependencies: d,
		log:          log,
		tracer:       tracer,
	}, nil
}

func (s *UserService) Register(ctx context.Context, username string, password string) (*models.User, error) {
	log := s.log.With(
		slog.String("function", "UserService.Register"),
		slog.String("username", username),
	)

	log.Debug("register user start")

	hashPassword, err := s.PasswordManager.HashPassword(password)
	if err != nil {
		log.Error("failed to hash password", slog.Any("error", err.Error()))

		return nil, err
	}

	user := models.User{
		Username:     username,
		HashPassword: hashPassword,
	}

	createdUser, err := s.UserRepository.Create(ctx, &user)
	if err != nil {
		log.Info("failed to create user", slog.Any("error", err.Error()))

		return nil, err
	}

	log.Debug("register user end")

	return createdUser, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (at string, rt string, err error) {
	ctx, span := s.tracer.Start(ctx, "AuthHandler.Login")
	defer span.End()

	log := s.log.With(
		slog.String("function", "UserService.Login"),
		slog.String("username", username),
	)

	log.Debug("login user start")

	getUser := func(username string) (*models.User, error) {
		if cacheUser, err := s.UserCache.Get(ctx, username); err == nil {
			return &cacheUser, nil
		}

		log.Info("cache miss")

		dbUser, err := s.UserRepository.GetByUsername(ctx, username)
		if err != nil {
			return nil, err
		}

		return dbUser, nil
	}

	user, err := getUser(username)
	if err != nil {
		log.Info("failed to get user", slog.Any("error", err.Error()))

		return "", "", err
	}

	if !s.PasswordManager.CheckPassword(password, user.HashPassword) {
		log.Info("invalid password")

		return "", "", ErrInvalidPassword
	}

	accessToken, err := s.AtManager.NewToken()
	if err != nil {
		log.Error("failed to create access token", slog.Any("error", err.Error()))

		return "", "", err
	}

	refreshToken, err := s.RtManager.NewToken()
	if err != nil {
		log.Error("failed to create refresh token", slog.Any("error", err.Error()))

		return "", "", err
	}

	session := &models.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	}

	slog.Debug("creating a new session", slog.Any("session", session))

	if _, err := s.SessionRepository.Create(ctx, session); err != nil {
		log.Info("failed to create session", slog.Any("error", err.Error()))

		return "", "", err
	}

	at = accessToken.Token
	rt = refreshToken.Token

	if err := s.UserCache.Set(ctx, username, *user, ttl); err != nil {
		log.Info("failed add user to cache", slog.String("error", err.Error()))
	}

	log.Debug("login user end")

	return at, rt, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (at string, rt string, err error) {
	log := s.log.With(
		slog.String("function", "UserService.Refresh"),
	)

	log.Debug("refresh tokens start")

	session, err := s.SessionRepository.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Info("failed to get session", slog.Any("error", err.Error()))

		return "", "", err
	}

	log.Debug("refresh session", slog.Any("session", session))

	if !session.Valid() {
		log.Info("session is unavailable", slog.Any("session", session))

		return "", "", ErrUnauthorizedRefresh
	}

	newAccessToken, err := s.AtManager.NewToken()
	if err != nil {
		log.Error("failed to create access token", slog.Any("error", err.Error()))

		return "", "", err
	}

	newRefreshToken, err := s.RtManager.NewToken()
	if err != nil {
		log.Error("failed to create refresh token", slog.Any("error", err.Error()))

		return "", "", err
	}

	session.RefreshToken = newRefreshToken

	if _, err := s.SessionRepository.Update(ctx, session); err != nil {
		log.Info("failed to update session", slog.Any("error", err.Error()))

		return "", "", err
	}

	at = newAccessToken.Token
	rt = newRefreshToken.Token

	log.Debug("refresh tokens end")

	return at, rt, nil
}
