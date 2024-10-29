package services

import (
	"auth/internal/auth/domain/models"
	"auth/internal/auth/domain/repository"
	"context"
	"github.com/pkg/errors"
	"log/slog"
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

type Dependencies struct {
	UserRepository    repository.UserRepository
	SessionRepository repository.SessionRepository
	TxManager         repository.TransactionManager
	AtManager         AccessTokenManager
	RtManager         RefreshTokenManager
	PasswordManager   PasswordManager
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

	return nil
}

type UserService struct {
	Dependencies
	log *slog.Logger
}

func NewUserService(d Dependencies, log *slog.Logger) (*UserService, error) {
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
	log := s.log.With(
		slog.String("function", "UserService.Login"),
		slog.String("username", username),
	)

	log.Debug("login user start")

	if err := s.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		user, err := s.UserRepository.GetByUsername(txCtx, username)
		if err != nil {
			log.Info("failed to get user", slog.Any("error", err.Error()))

			return err
		}

		if !s.PasswordManager.CheckPassword(password, user.HashPassword) {
			log.Info("invalid password")

			return ErrInvalidPassword
		}

		accessToken, err := s.AtManager.NewToken()
		if err != nil {
			log.Error("failed to create access token", slog.Any("error", err.Error()))

			return err
		}

		refreshToken, err := s.RtManager.NewToken()
		if err != nil {
			log.Error("failed to create refresh token", slog.Any("error", err.Error()))

			return err
		}

		session := &models.Session{
			UserID:       user.ID,
			RefreshToken: refreshToken,
		}

		slog.Debug("creating a new session", slog.Any("session", session))

		if _, err := s.SessionRepository.Create(txCtx, session); err != nil {
			log.Info("failed to create session", slog.Any("error", err.Error()))

			return err
		}

		at = accessToken.Token
		rt = refreshToken.Token

		return nil
	}); err != nil {
		return "", "", err
	}

	log.Debug("login user end")

	return at, rt, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (at string, rt string, err error) {
	log := s.log.With(
		slog.String("function", "UserService.Refresh"),
	)

	log.Debug("refresh tokens start")

	if err := s.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		session, err := s.SessionRepository.GetByRefreshToken(txCtx, refreshToken)
		if err != nil {
			log.Info("failed to get session", slog.Any("error", err.Error()))

			return err
		}

		log.Debug("refresh session", slog.Any("session", session))

		if !session.Valid() {
			log.Info("session is unavailable", slog.Any("session", session))

			return ErrUnauthorizedRefresh
		}

		newAccessToken, err := s.AtManager.NewToken()
		if err != nil {
			log.Error("failed to create access token", slog.Any("error", err.Error()))

			return err
		}

		newRefreshToken, err := s.RtManager.NewToken()
		if err != nil {
			log.Error("failed to create refresh token", slog.Any("error", err.Error()))

			return err
		}

		session.RefreshToken = newRefreshToken

		if _, err := s.SessionRepository.Update(txCtx, session); err != nil {
			log.Info("failed to update session", slog.Any("error", err.Error()))

			return err
		}

		at = newAccessToken.Token
		rt = newRefreshToken.Token

		return nil
	}); err != nil {
		return "", "", err
	}

	log.Debug("refresh tokens end")

	return at, rt, nil
}
