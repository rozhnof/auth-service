package main

import (
	http_app "auth/internal/app/http"
	user_services "auth/internal/application/services/user"
	postgres_user_repository "auth/internal/infrastructure/repository/postgres/user"
	"auth/internal/pkg/config"
	postgres_database "auth/internal/pkg/database/postgres"
	http_server "auth/internal/pkg/server/http"
	handlers "auth/internal/presentation/http/user/handlers"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "auth/docs"

	"github.com/gin-gonic/gin"
)

const (
	EnvConfigPath = "CONFIG_PATH"
)

// @title           Authentication Service API
// @version         1.0

// @host      localhost:8080
// @BasePath  /

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.NewConfig(os.Getenv(EnvConfigPath))
	if err != nil {
		log.Fatal(err)
	}

	InitLogging(cfg.Logger.Level)

	postgresDatabase, err := postgres_database.NewDatabase(ctx, postgres_database.CreateConnectionString(cfg.Repository))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userRepository = postgres_user_repository.NewUserRepository(postgresDatabase)
		userService    = user_services.NewAuthService(userRepository)
		userHandler    = handlers.NewAuthHandler(userService)
		router         = gin.New()
	)

	http_app.InitRoutes(router, userHandler)

	log.Println("http server started successfully")
	if err := RunHTTP(ctx, cfg, router); err != nil {
		log.Fatal(err)
	}
	log.Println("http server stopped successfully")
}

func InitLogging(logLevel string) {
	var level slog.Leveler

	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	var (
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
		logger  = slog.New(handler)
	)

	slog.SetDefault(logger)
}

func RunHTTP(ctx context.Context, cfg *config.Config, router *gin.Engine) error {
	serverCfg := http_server.Config{
		Address:         cfg.Server.HTTP.Address,
		ShutdownTimeout: cfg.Server.HTTP.ShutdownTimeout,
	}

	server := http_server.New(
		serverCfg,
		router,
	)

	return server.Run(ctx)
}
