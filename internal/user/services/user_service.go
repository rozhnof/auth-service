package services

import (
	"auth/internal/user/models"
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

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID models.UserID) (*models.User, error)
	GetByUsername(ctx context.Context, username models.Username) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID models.UserID) (*time.Time, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewAuthService(repository UserRepository) *UserService {
	return &UserService{
		userRepository: repository,
	}
}

func (s *UserService) Register(ctx context.Context, username string, password string) (*models.User, error) {
	user := models.User{
		Username: username,
		Password: models.UserPasswordFromString(password),
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

	if !user.Password.Compare(models.UserPasswordFromString(password)) {
		return "", "", ErrInvalidPassword
	}

	var (
		accessToken  = models.NewAccessToken(atTimeout)
		refreshToken = models.NewRefreshToken(rtTimeout)
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

	if !user.RefreshToken.Compare(models.RefreshTokenFromString(rt)) {
		return "", "", ErrUnauthorizedRefresh
	}

	var (
		accessToken  = models.NewAccessToken(atTimeout)
		refreshToken = models.NewRefreshToken(rtTimeout)
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
