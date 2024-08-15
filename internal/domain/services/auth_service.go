package services

import (
	"auth/internal/domain/models"
	"auth/internal/infrastructure/repository/postgres/dao"
	"auth/pkg/mail"
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type authRepository interface {
	Create(ctx context.Context, user dao.User) (*uuid.UUID, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*dao.User, error)
	List(ctx context.Context) ([]dao.User, error)
	Update(ctx context.Context, user dao.User) error
	Delete(ctx context.Context, userID uuid.UUID) error
	WithTransaction(ctx context.Context, f func(ctxWithTransaction context.Context) error) error
}

type Config struct {
	SecretKey           []byte
	AccessTokenTimeout  time.Duration
	RefreshTokenTimeout time.Duration
}

type AuthService struct {
	repository authRepository
	mailSender mail.MailSender
	cfg        Config
}

func NewAuthService(repository authRepository, mailSender mail.MailSender, cfg Config) *AuthService {
	return &AuthService{
		repository: repository,
		mailSender: mailSender,
		cfg:        cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, email string) (*models.User, error) {
	userDAO := dao.User{
		Email: email,
	}

	err := s.repository.WithTransaction(ctx, func(ctxWithTransaction context.Context) error {
		userID, err := s.repository.Create(ctxWithTransaction, userDAO)
		if err != nil {
			return err
		}

		userDAO.ID = *userID

		return nil
	})

	if err != nil {
		return nil, err
	}

	user := userDAO.ToModel()
	return &user, nil
}

func (s *AuthService) Login(ctx context.Context, userID uuid.UUID, userIP string) (*models.TokenPair, error) {
	var tokenPair *models.TokenPair

	err := s.repository.WithTransaction(ctx, func(ctxWithTransaction context.Context) error {
		userDAO, err := s.repository.GetByID(ctxWithTransaction, userID)
		if err != nil {
			return err
		}

		tokenPair, err = generateTokenPair(userID, userIP, s.cfg)
		if err != nil {
			return err
		}

		userDAO.TokenHash = tokenPair.RefreshToken.Hash

		if err := s.repository.Update(ctxWithTransaction, *userDAO); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *AuthService) Refresh(ctx context.Context, ip string, accessToken string, refreshToken string) (*models.TokenPair, error) {
	accessTokenClaims, err := getTokenClaims(accessToken, s.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	refreshTokenClaims, err := getTokenClaims(refreshToken, s.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	if err := refreshTokenClaims.Valid(); err != nil {
		return nil, errors.Wrap(ErrUnauthorizedRefresh, err.Error())
	}

	if accessTokenClaims.TokenID != refreshTokenClaims.TokenID {
		return nil, errors.Wrap(ErrUnauthorizedRefresh, "the given tokens are not paired")
	}

	var (
		userID       = refreshTokenClaims.UserID
		newTokenPair *models.TokenPair
	)

	if err := s.repository.WithTransaction(ctx, func(ctxWithTransaction context.Context) error {
		userDAO, err := s.repository.GetByID(ctxWithTransaction, userID)
		if err != nil {
			return err
		}

		if ip != refreshTokenClaims.UserIP {
			var (
				sender    = "auth-service"
				recipient = userDAO.Email
				subject   = "Notification of the security system"
				message   = fmt.Sprintf("There was a recent login to your account. IP: %s", ip)
			)
			s.mailSender.SendMessage(sender, recipient, subject, message)
		}

		if !cmpTokenWithHash([]byte(refreshToken), userDAO.TokenHash) {
			return errors.Wrap(ErrUnauthorizedRefresh, "refresh token has expired")
		}

		newTokenPair, err = generateTokenPair(userID, ip, s.cfg)
		if err != nil {
			return err
		}

		userDAO.TokenHash = newTokenPair.RefreshToken.Hash

		if err := s.repository.Update(ctxWithTransaction, *userDAO); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return newTokenPair, nil
}

func generateBcryptHash(in []byte) ([]byte, error) {
	sha512Hash := generateSHA512Hash(in)

	hash, err := bcrypt.GenerateFromPassword(sha512Hash[:], bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func generateTokenPair(userID uuid.UUID, userIP string, cfg Config) (*models.TokenPair, error) {
	tokenID := uuid.New()

	at, err := generateAccessToken(tokenID, userID, userIP, cfg.AccessTokenTimeout, cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	rt, err := generateRefreshToken(tokenID, userID, userIP, cfg.RefreshTokenTimeout, cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	tokenPair := models.TokenPair{
		AccessToken:  *at,
		RefreshToken: *rt,
	}

	return &tokenPair, nil
}

func generateAccessToken(tokenID uuid.UUID, userID uuid.UUID, userIP string, timeout time.Duration, secretKey []byte) (*models.AccessToken, error) {
	token, err := generateToken(tokenID, timeout, userID, userIP, secretKey)
	if err != nil {
		return nil, err
	}

	accessToken := models.AccessToken{
		Token: token,
	}

	return &accessToken, nil
}

func generateRefreshToken(tokenID uuid.UUID, userID uuid.UUID, userIP string, timeout time.Duration, secretKey []byte) (*models.RefreshToken, error) {
	token, err := generateToken(tokenID, timeout, userID, userIP, secretKey)
	if err != nil {
		return nil, err
	}

	hashToken, err := generateBcryptHash([]byte(token))
	if err != nil {
		return nil, err
	}

	refreshToken := models.RefreshToken{
		Token: token,
		Hash:  hashToken,
	}

	return &refreshToken, nil
}

func generateToken(tokenID uuid.UUID, timeout time.Duration, userID uuid.UUID, userIP string, secretKey []byte) (string, error) {
	expiredAt := time.Now().Add(timeout)

	claims := &models.TokenClaims{
		TokenID: tokenID,
		UserID:  userID,
		UserIP:  userIP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getTokenClaims(t string, secretKey []byte) (*models.TokenClaims, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, err
	}

	bytes, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, errors.Wrap(err, "error marshal token claims")
	}

	var claims models.TokenClaims
	if err := json.Unmarshal(bytes, &claims); err != nil {
		return nil, errors.Wrap(err, "error unmarshal token claims")
	}

	return &claims, nil
}

func cmpTokenWithHash(token []byte, hash []byte) bool {
	sha512Hash := generateSHA512Hash(token)
	return bcrypt.CompareHashAndPassword(hash, sha512Hash[:]) == nil
}

func generateSHA512Hash(in []byte) [64]byte {
	return sha512.Sum512(in)
}
