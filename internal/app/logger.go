package app

import (
	"io"
	"log/slog"
	"os"

	"github.com/go-slog/otelslog"
	"github.com/rozhnof/auth-service/internal/pkg/config"
)

func NewLogger(cfg config.Logger) (*slog.Logger, error) {
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

	var logOutput io.Writer

	if cfg.Path != "" {
		logFile, err := os.OpenFile(cfg.Path, os.O_RDWR|os.O_SYNC, 0644)
		if err != nil {
			return nil, err
		}

		logOutput = logFile
	} else {
		logOutput = os.Stdout
	}

	handler := otelslog.NewHandler(
		slog.NewJSONHandler(
			logOutput,
			&slog.HandlerOptions{Level: level},
		),
	)

	return slog.New(handler), nil
}
