package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"

	ttcli "github.com/Anton-Kraev/event-timeslot-planner/internal/client/timetable"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/config"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/logger"
	ttrepo "github.com/Anton-Kraev/event-timeslot-planner/internal/repository/redis/timetable"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/service/schedule"
)

func Run() {
	cfg := config.MustInit()

	log := logger.MustSetup(cfg.Env)
	log.Info("starting event-timeslot-planner", slog.String("env", string(cfg.Env)))
	log.Debug("debug logging enabled")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	ttCache := ttrepo.NewRedisRepository(redisClient, cfg.Redis.ExpirationPeriod)

	ctx := context.Background()
	pingRes, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Error("error connecting redis", logger.Err(err))
		os.Exit(1)
	}

	log.Debug("pinging redis", slog.String("result", pingRes))

	httpClient := &http.Client{Timeout: cfg.TimetableAPI.Timeout}
	ttClient := ttcli.NewClient(cfg.TimetableAPI.Address, httpClient)

	_ = schedule.NewService(ttClient, ttCache)

	// TODO: init schedule controller
	// TODO: init router: chi/stdlib
	// TODO: run server
}
