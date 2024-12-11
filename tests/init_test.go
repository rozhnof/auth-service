package tests

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rozhnof/auth-service/internal/infrastructure/database/postgres"
	"github.com/rozhnof/auth-service/internal/pkg/config"
)

const (
	configPath = "../config/test-config.yaml"
)

var (
	authServiceDatabase AuthServiceDatabase
	authServiceClient   AuthServiceClient
)

const (
	authServiceAddress = "http://localhost:9090"
	authServiceTimeout = time.Second * 5
)

type Config struct {
	Postgres config.Postgres
}

func init() {
	cfg, err := config.NewConfig[Config](configPath)
	if err != nil {
		log.Fatal(err)
	}

	postgresCfg := postgres.DatabaseConfig{
		Address:  cfg.Postgres.Address,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DB:       cfg.Postgres.DB,
		SSL:      cfg.Postgres.SSL,
	}

	db, err := postgres.NewDatabase(context.Background(), postgresCfg)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{
		Timeout: authServiceTimeout,
	}

	authServiceDatabase = NewAuthServiceDatabase(db)
	authServiceClient = NewAuthServiceClient(httpClient, authServiceAddress)
}

func SetUp() error {
	if err := Erase(); err != nil {
		return err
	}

	return nil
}

func Erase() error {
	if err := authServiceDatabase.Truncate(context.Background()); err != nil {
		return err
	}

	return nil
}
