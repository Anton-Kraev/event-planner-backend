package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/config"
	ttcli "github.com/Anton-Kraev/event-timeslot-planner/internal/http/client/timetable"
	schedhndl "github.com/Anton-Kraev/event-timeslot-planner/internal/http/server/handler/schedule"
	mw "github.com/Anton-Kraev/event-timeslot-planner/internal/http/server/middleware"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/lib/logger"
	ttrepo "github.com/Anton-Kraev/event-timeslot-planner/internal/repository/redis/timetable"
	schedsrvc "github.com/Anton-Kraev/event-timeslot-planner/internal/service/schedule"
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

	scheduleService := schedsrvc.NewService(ttClient, ttCache)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mw.NewLogger(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	scheduleHandler := schedhndl.NewHandler(scheduleService, log)

	router.Get("/timetable/get_schedule", scheduleHandler.GetTimetableSchedule)

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
