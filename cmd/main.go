package main

import (
	http_app "auth/internal/app/http"
	"auth/internal/pkg/config"
	"context"
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

	logger, err := NewLogger(cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}

	app, err := http_app.NewApp(ctx, cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := app.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		if err := app.Stop(ctx); err != nil {
			log.Fatalf("error in app.Stop: %s", err)
		}
	}
}

func NewLogger(cfg config.LoggerConfig) (*slog.Logger, error) {
	var level slog.Leveler

	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	logFile, err := os.OpenFile(cfg.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: level})

	return slog.New(handler), nil
}
