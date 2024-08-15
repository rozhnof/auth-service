package main

import (
	handlers "auth/internal/application/handlers/http"
	"auth/internal/domain/services"
	repository "auth/internal/infrastructure/repository/postgres"
	"auth/internal/pkg/config"
	"auth/pkg/config_reader"
	"auth/pkg/database"
	"auth/pkg/mail"
	"auth/pkg/server"
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "auth/docs"

	"github.com/gorilla/mux"
)

const (
	cfgPath  = "config/auth_config.yaml"
	certPath = "cert.pem"
	keyPath  = "key.pem"
)

// @title           Authentication Service API
// @version         1.0

// @host      localhost:8080
// @BasePath  /

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := loadConfig(cfgPath)
	if err != nil {
		log.Fatalf("error loading config: %s", err)
	}

	secrets, err := loadSecrets()
	if err != nil {
		log.Fatalf("error loading secrets: %s", err)
	}

	cert, err := loadAuthCertificate(certPath, keyPath)
	if err != nil {
		log.Fatalf("failed to load key pair: %v", err)
	}

	db, err := database.NewDatabase(ctx, secrets.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.GetPool().Close()

	serverCfg := server.Config{
		Address:         cfg.Server.Address,
		ShutdownTimeout: cfg.Server.ShutdownTimeout,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				*cert,
			},
		},
	}

	serviceCfg := services.Config{
		SecretKey:           []byte(secrets.SecretKey),
		AccessTokenTimeout:  cfg.Service.AccessTokenTimeout,
		RefreshTokenTimeout: cfg.Service.RefreshTokenTimeout,
	}

	var (
		authMailSender = mail.NewMailSender(cfg.Mail.Sender, cfg.Mail.Password, cfg.Mail.Host)
		authRepository = repository.NewAuthRepository(db)
		authService    = services.NewAuthService(authRepository, *authMailSender, serviceCfg)
		authHandler    = handlers.NewAuthHandler(authService)
		authRouter     = mux.NewRouter()
	)

	initRoutes(authRouter, authHandler)

	authServer := server.New(serverCfg, authRouter)
	if err := authServer.RunTLS(ctx); err != nil {
		log.Fatal(err)
	}
}

func loadAuthCertificate(certPath string, keyPath string) (*tls.Certificate, error) {
	authCertificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &authCertificate, nil
}

func loadSecrets() (*config.AuthSecrets, error) {
	var secrets config.AuthSecrets
	if err := env.Parse(&secrets); err != nil {
		return nil, err
	}

	if err := secrets.Validate(); err != nil {
		return nil, err
	}

	return &secrets, nil
}

func loadConfig(authConfigPath string) (*config.AuthConfig, error) {
	var cfg config.AuthConfig
	if err := config_reader.LoadYaml(authConfigPath, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func initRoutes(router *mux.Router, handler *handlers.AuthHandler) {
	router.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.Register(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	router.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.Login(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	router.HandleFunc("/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.Refresh(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
