package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	ttcli "github.com/Anton-Kraev/event-timeslot-planner/internal/client/timetable"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/config"
	ttrepo "github.com/Anton-Kraev/event-timeslot-planner/internal/repository/redis/timetable"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/service/schedule"
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

	redisConfig := viper.Sub("redis")
	if redisConfig == nil {
		log.Fatal("invalid config: redis is empty")
	}

	redisAddr := redisConfig.GetString("address")
	if redisAddr == "" {
		log.Fatal("invalid redis config: address is empty")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisConfig.GetString("password"),
		DB:       redisConfig.GetInt("db"),
	})

	ctx := context.Background()
	pingRes, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("error connecting redis: %s", err.Error())
	}

	logger.Info("pinging redis", slog.String("result", pingRes))

	defer redisClient.Close()

	redisExpirationPeriod := redisConfig.GetDuration("expiration_period")
	if redisExpirationPeriod == 0 {
		log.Fatal("invalid redis config: expiration_period is empty")
	}

	ttCache := ttrepo.NewRedisRepository(redisClient, redisExpirationPeriod)

	ttConfig := viper.Sub("timetable_api")
	if ttConfig == nil {
		log.Fatal("invalid config: timetable_api is empty")
	}

	ttAddress := ttConfig.GetString("address")
	if ttAddress == "" {
		log.Fatal("invalid timetable_api config: address is empty")
	}

	ttTimeout := ttConfig.GetDuration("timeout")
	if ttTimeout == 0 {
		log.Fatal("invalid timetable_api config: timeout is empty")
	}

	httpClient := &http.Client{Timeout: ttTimeout}
	ttClient := ttcli.NewClient(ttAddress, httpClient)

	_ = schedule.NewService(ttClient, ttCache)

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
