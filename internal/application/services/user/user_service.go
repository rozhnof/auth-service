package user_services

import (
	"auth/internal/domain/models"
	repository_interfaces "auth/internal/domain/repository"
	"auth/internal/domain/value_object"
	"context"
	"time"
)

const (
	atTimeout = time.Hour
	rtTimeout = time.Hour * 72
)

var (
	secretKey = []byte{'a', 'v'}
)

type UserService struct {
	userRepository repository_interfaces.UserRepository
}

func NewAuthService(repository repository_interfaces.UserRepository) *UserService {
	return &UserService{
		userRepository: repository,
	}
}

func (s *UserService) Register(ctx context.Context, email string) (*models.User, error) {
	user := models.User{
		Username: email,
	}

	createdUser, err := s.userRepository.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (string, string, error) {
	user, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if !user.Password.Compare(value_object.PasswordFromString(password)) {
		return "", "", ErrInvalidPassword
	}

	var (
		accessToken  = value_object.NewAccessToken(atTimeout)
		refreshToken = value_object.NewRefreshToken(rtTimeout)
	)

	signedAccessToken, err := accessToken.Sign(secretKey)
	if err != nil {
		return "", "", err
	}

	user.RefreshToken = refreshToken

	if _, err := s.userRepository.Update(ctx, user); err != nil {
		return "", "", err
	}

	return signedAccessToken, user.RefreshToken.String(), nil
}

func (s *UserService) Refresh(ctx context.Context, username string, rt string) (string, string, error) {
	user, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if !user.RefreshToken.Compare(value_object.RefreshTokenFromString(rt)) {
		return "", "", ErrUnauthorizedRefresh
	}

	var (
		accessToken  = value_object.NewAccessToken(atTimeout)
		refreshToken = value_object.NewRefreshToken(rtTimeout)
	)

	signedAccessToken, err := accessToken.Sign(secretKey)
	if err != nil {
		return "", "", err
	}

	user.RefreshToken = refreshToken

	if _, err := s.userRepository.Update(ctx, user); err != nil {
		return "", "", err
	}

	return signedAccessToken, user.RefreshToken.String(), nil
}
