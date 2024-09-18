package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type EnvType string

const (
	envLocal EnvType = "local"
	envProd  EnvType = "prod"
)

func Setup(env EnvType) (*slog.Logger, error) {
	switch env {
	case envLocal:
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		), nil
	case envProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		), nil
	}

	return nil, fmt.Errorf("invalid env: %s", env)
}
