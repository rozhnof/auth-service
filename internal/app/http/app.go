package http_app

import (
	"auth/internal/auth/application/services"
	postgres_session_repository "auth/internal/auth/infrastructure/repository/postgres/session"
	postgres_user_repository "auth/internal/auth/infrastructure/repository/postgres/user"
	http_handlers "auth/internal/auth/presentation/handlers/http"
	"auth/internal/pkg/config"
	postgres_database "auth/internal/pkg/database/postgres"
	"auth/internal/pkg/password_manager"
	http_server "auth/internal/pkg/server/http"
	"auth/internal/pkg/token_manager"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type App struct {
	server *http_server.HTTPServer
	log    *slog.Logger
}

func NewApp(ctx context.Context, cfg *config.Config, log *slog.Logger) (*App, error) {
	httpServerConfig := http_server.Config{
		Address:         cfg.Server.HTTP.Address,
		ShutdownTimeout: cfg.Server.HTTP.ShutdownTimeout,
		TLSConfig:       nil,
	}

	databaseConfig := postgres_database.Config{
		Address:  cfg.Repository.Address,
		Port:     cfg.Repository.Port,
		User:     cfg.Repository.User,
		Password: cfg.Repository.Password,
		DB:       cfg.Repository.DB,
		SSL:      cfg.Repository.SSL,
	}

	postgresDatabase, err := postgres_database.NewDatabase(ctx, databaseConfig)
	if err != nil {
		return nil, err
	}

	var (
		transactionManager = postgres_database.NewTransactionManager(postgresDatabase)
		userRepository     = postgres_user_repository.NewUserRepository(transactionManager, log)
		sessionRepository  = postgres_session_repository.NewSessionRepository(transactionManager, log)
		atManager          = token_manager.NewAccessTokenManager(cfg.Service.Tokens.Access.Timeout, []byte(cfg.Service.Tokens.Access.SecretKey))
		rtManager          = token_manager.NewRefreshTokenManager(cfg.Service.Tokens.Refresh.Timeout)
		passwordManager    = password_manager.NewPasswordManager()
	)

	userServiceDependencies := services.Dependencies{
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
		TxManager:         transactionManager,
		AtManager:         atManager,
		RtManager:         rtManager,
		PasswordManager:   passwordManager,
	}

	userService, err := services.NewUserService(userServiceDependencies, log)
	if err != nil {
		return nil, err
	}

	authHandler := http_handlers.NewAuthHandler(userService, log)
	router := gin.New()

	InitRoutes(router, authHandler)

	server := http_server.New(httpServerConfig, router)

	return &App{
		server: server,
		log:    log,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("starting http server")

	return a.server.Run(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info("stopping http server")

	return a.server.Shutdown(ctx)
}
