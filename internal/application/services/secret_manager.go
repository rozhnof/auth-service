package services

import "github.com/rozhnof/auth-service/internal/domain"

type SecretManager interface {
	SecretKey() domain.Secret
}
