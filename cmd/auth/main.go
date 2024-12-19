package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	_ "github.com/rozhnof/auth-service/docs"
	"github.com/rozhnof/auth-service/internal/app"
	"github.com/rozhnof/auth-service/internal/app/auth"
	"github.com/rozhnof/auth-service/internal/infrastructure/kafka"
	"github.com/rozhnof/auth-service/internal/pkg/config"
	"github.com/rozhnof/auth-service/internal/pkg/server"
)

const (
	EnvConfigPath = "CONFIG_PATH"
	pprofAddress  = ":6060"
)

// @title           Authentication Service API
// @version         1.0

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @host      localhost:8080
// @BasePath  /

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.NewConfig[auth.Config](os.Getenv(EnvConfigPath))
	if err != nil {
		log.Fatal(err)
	}

	logger, err := app.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("init logger success")

	postgresDatabase, err := app.NewPostgresDatabase(ctx, cfg.Postgres)
	if err != nil {
		logger.Error("init postgres failed", slog.String("error", err.Error()))
		return
	}
	defer postgresDatabase.Close()
	logger.Info("init postgres success")

	redisDatabase, err := app.NewRedisDatabase(ctx, cfg.Redis)
	if err != nil {
		logger.Error("init redis failed", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := redisDatabase.Close(); err != nil {
			logger.Error("close redis failed", slog.String("error", err.Error()))
		}
	}()
	logger.Info("init redis success")

	tracer, shutdown, err := app.NewTracer(ctx, cfg.Tracing, auth.ServiceName)
	if err != nil {
		logger.Error("init tracer failed", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			logger.Error("close tracer failed", slog.String("error", err.Error()))
		}
	}()
	logger.Info("init tracer success")

	kafkaProducer, err := kafka.NewProducer(cfg.Kafka.BrokerList)
	if err != nil {
		logger.Error("init kafka producer failed", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			logger.Error("close kafka producer failed", slog.String("error", err.Error()))
		}
	}()
	logger.Info("init kafka producer success")

	authApp := auth.NewApp(
		ctx,
		cfg,
		logger,
		postgresDatabase,
		redisDatabase,
		kafkaProducer,
		tracer,
	)
	logger.Info("init app success")

	go func() {
		logger.Info("start pprof server")

		pprofServer := server.NewHTTPServer(ctx, pprofAddress, nil, logger)

		if err := pprofServer.Run(ctx); err != nil {
			logger.Error("pprof server error", slog.String("error", err.Error()))
			return
		}

		logger.Info("shutdown pprof server")
	}()

	logger.Info("start app")
	if err := authApp.Run(ctx); err != nil {
		logger.Error("app error", slog.String("error", err.Error()))
		return
	}

	logger.Error("shutdown app")
}
