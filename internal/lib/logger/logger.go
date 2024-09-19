package logger

import (
	"log/slog"
	"os"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/config"
)

func Setup(env config.EnvType) *slog.Logger {
	if env == config.EnvLocal {
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
}
