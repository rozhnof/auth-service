package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"github.com/rozhnof/auth-service/internal/domain"
	"github.com/rozhnof/auth-service/internal/domain/entities"
	vobjects "github.com/rozhnof/auth-service/internal/domain/value_objects"
	"go.opentelemetry.io/otel/trace"
)

type AuthServiceConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type AuthService struct {
	repository        UserRepository
	loginMsgSender    LoginMessageSender
	registerMsgSender RegisterMessageSender
	txManager         TransactionManager
	secretManager     SecretManager
	log               *slog.Logger
	tracer            trace.Tracer
	cfg               AuthServiceConfig
}

func NewAuthService(
	repository UserRepository,
	loginMsgSender LoginMessageSender,
	registerMsgSender RegisterMessageSender,
	txManager TransactionManager,
	secretManager SecretManager,
	log *slog.Logger,
	tracer trace.Tracer,
	cfg AuthServiceConfig,
) *AuthService {
	return &AuthService{
		repository:        repository,
		loginMsgSender:    loginMsgSender,
		registerMsgSender: registerMsgSender,
		txManager:         txManager,
		secretManager:     secretManager,
		log:               log,
		tracer:            tracer,
		cfg:               cfg,
	}
}

func (s *AuthService) OAuthLogin(ctx context.Context, email string) (*entities.User, error) {
	ctx, span := s.tracer.Start(ctx, "AuthService.OAuthLogin")
	defer span.End()

	if user, err := s.repository.GetByEmail(ctx, email); err == nil {
		if err := user.UpdateTokens(s.cfg.AccessTokenTTL, s.cfg.RefreshTokenTTL, s.secretManager.SecretKey().Get()); err != nil {
			return nil, err
		}

		if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
			if err := s.repository.Update(ctx, user); err != nil {
				return err
			}

			loginMsg := LoginMessage{
				Email: email,
			}

			if err := s.loginMsgSender.SendMessage(ctx, loginMsg); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, err
		}

		return user, nil
	}

	const randomPasswordLen = 72

	passwordStr, err := domain.GenerateRandomString(randomPasswordLen)
	if err != nil {
		return nil, err
	}

	password, err := vobjects.NewPassword(passwordStr)
	if err != nil {
		return nil, err
	}

	user := entities.NewUser(email, password)

	if err := user.UpdateTokens(s.cfg.AccessTokenTTL, s.cfg.RefreshTokenTTL, s.secretManager.SecretKey().Get()); err != nil {
		return nil, err
	}

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.repository.Create(ctx, user); err != nil {
			return err
		}

		registerMsg := RegisterMessage{
			Email: email,
		}

		if err := s.registerMsgSender.SendMessage(ctx, registerMsg); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Register(ctx context.Context, email string, passwordStr string) (*entities.User, error) {
	ctx, span := s.tracer.Start(ctx, "AuthService.Register")
	defer span.End()

	var user *entities.User

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		password, err := vobjects.NewPassword(passwordStr)
		if err != nil {
			return err
		}

		user = entities.NewUser(email, password)

		if err := s.repository.Create(ctx, user); err != nil {
			return err
		}

		registerMsg := RegisterMessage{
			Email: email,
		}

		if err := s.registerMsgSender.SendMessage(ctx, registerMsg); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (at string, rt string, err error) {
	ctx, span := s.tracer.Start(ctx, "AuthService.Login")
	defer span.End()

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		user, err := s.repository.GetByEmail(ctx, email)
		if err != nil {
			return err
		}

		if !user.Password().Compare(password) {
			return ErrInvalidPassword
		}

		if err := user.UpdateTokens(s.cfg.AccessTokenTTL, s.cfg.RefreshTokenTTL, s.secretManager.SecretKey().Get()); err != nil {
			return err
		}

		if err := s.repository.Update(ctx, user); err != nil {
			return err
		}

		loginMsg := LoginMessage{
			Email: email,
		}

		if err := s.loginMsgSender.SendMessage(ctx, loginMsg); err != nil {
			return err
		}

		at = user.AccessToken().Token()
		rt = user.RefreshToken().Token()

		return nil
	}); err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (at string, rt string, err error) {
	ctx, span := s.tracer.Start(ctx, "AuthService.Refresh")
	defer span.End()

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		user, err := s.repository.GetByRefreshToken(ctx, refreshToken)
		if err != nil {
			return err
		}

		if !user.RefreshToken().Valid() {
			return errors.Wrap(ErrUnauthorizedRefresh, "invalid refresh token")
		}

		if err := user.UpdateTokens(s.cfg.AccessTokenTTL, s.cfg.RefreshTokenTTL, s.secretManager.SecretKey().Get()); err != nil {
			return err
		}

		if err := s.repository.Update(ctx, user); err != nil {
			return err
		}

		at = user.AccessToken().Token()
		rt = user.RefreshToken().Token()

		return nil
	}); err != nil {
		return "", "", err
	}

	return at, rt, nil
}
