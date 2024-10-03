package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/Anton-Kraev/event-planner-backend/internal/config"
	ttcli "github.com/Anton-Kraev/event-planner-backend/internal/http/client/timetable"
	tthndl "github.com/Anton-Kraev/event-planner-backend/internal/http/server/handler/timetable"
	rtr "github.com/Anton-Kraev/event-planner-backend/internal/http/server/router"
	"github.com/Anton-Kraev/event-planner-backend/internal/lib/logger"
	ttrepo "github.com/Anton-Kraev/event-planner-backend/internal/repository/redis/timetable"
	ttsrvc "github.com/Anton-Kraev/event-planner-backend/internal/service/timetable"
)

func Run() {
	cfg := config.MustInit()

	log := logger.Setup(cfg.Env)
	log.Info("starting event-timeslot-planner", slog.String("env", string(cfg.Env)))
	log.Debug("debug logging enabled")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			log.Warn("failed to close redis connection", logger.Err(err))
		}
	}(redisClient)

	ttCache := ttrepo.NewRedisRepository(redisClient, cfg.Redis.UserExpPeriod, cfg.Redis.EventExpPeriod)

	ctx := context.Background()
	pingRes, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Error("error connecting redis", logger.Err(err))
		os.Exit(1)
	}

	log.Debug("pinging redis", slog.String("result", pingRes))

	httpClient := &http.Client{Timeout: cfg.TimetableAPI.Timeout}
	ttClient := ttcli.NewClient(cfg.TimetableAPI.Address, httpClient)

	scheduleService := ttsrvc.NewService(ttClient, ttCache)
	scheduleHandler := tthndl.NewHandler(scheduleService, log)

	router := rtr.SetupRouter(log, scheduleHandler)

	log.Info("starting server", slog.String("address", cfg.HttpServer.Address))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.Timeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
