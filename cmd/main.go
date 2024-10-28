package main

import (
	http_app "auth/internal/app/http"
	"auth/internal/pkg/config"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "auth/docs"
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

	logger := NewLogger(cfg.Logger)
	slog.SetDefault(logger)
	logger.Info("init logger")

	app, err := http_app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("init http application")

	go func() {
		logger.Info(fmt.Sprintf("http application listening on address %s", cfg.Server.HTTP.Address))

		if err := app.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		if err := app.Stop(ctx); err != nil {
			log.Fatalf("error in app.Stop: %s", err)
		}

		logger.Info("Gracefully stopped http server")
	}
}

func NewLogger(loggerConfig config.LoggerConfig) *slog.Logger {
	var level slog.Leveler

	switch loggerConfig.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	return slog.New(handler)
}
