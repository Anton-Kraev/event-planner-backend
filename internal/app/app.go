package app

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/viper"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/config"
)

type envType string

const (
	envLocal envType = "local"
	envProd  envType = "prod"
)

func Run() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	env := viper.GetString("env")

	logger, err := setupLogger(envType(env))
	if err != nil {
		log.Fatalf("error initializing logger: %s", err.Error())
	}

	logger.Info("starting event-timeslot-planner", slog.String("env", env))
	logger.Debug("debug logging enabled")

	// TODO: init tt cache
	// TODO: init schedule service
	// TODO: init schedule controller
	// TODO: init router: chi/stdlib
	// TODO: run server
}

func setupLogger(env envType) (*slog.Logger, error) {
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
