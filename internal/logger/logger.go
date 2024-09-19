package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/config"
)

func MustSetup(env config.EnvType) *slog.Logger {
	switch env {
	case config.EnvLocal:
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log.Fatalf("error initializing logger: invalid env: %s", env)

		return nil
	}
}
