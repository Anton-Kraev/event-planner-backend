package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
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

	ctx := context.Background()

	redisConfig := viper.GetStringMap("redis")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig["addr"].(string),
		Password: redisConfig["password"].(string),
		DB:       redisConfig["db"].(int),
	})

	pingRes, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("error connecting redis: %s", err.Error())
	}

	logger.Info("pinging redis", slog.String("result", pingRes))

	defer redisClient.Close()

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
