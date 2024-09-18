package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	ttcli "github.com/Anton-Kraev/event-timeslot-planner/internal/client/timetable"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/config"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/logger"
	ttrepo "github.com/Anton-Kraev/event-timeslot-planner/internal/repository/redis/timetable"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/service/schedule"
)

func Run() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	env := viper.GetString("env")

	customLog, err := logger.Setup(logger.EnvType(env))
	if err != nil {
		log.Fatalf("error initializing logger: %s", err.Error())
	}

	customLog.Info("starting event-timeslot-planner", slog.String("env", env))
	customLog.Debug("debug logging enabled")

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

	customLog.Info("pinging redis", slog.String("result", pingRes))

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
