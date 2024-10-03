package timetable

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
	"github.com/Anton-Kraev/event-planner-backend/internal/lib/api/response"
	"github.com/Anton-Kraev/event-planner-backend/internal/lib/logger"
)

type classroomsResponse struct {
	response.Response
	Classrooms []timetable.Classroom `json:"classrooms"`
}

func (h Handler) Classrooms(w http.ResponseWriter, r *http.Request) {
	const (
		op     = "http.server.handler.timetable.Classrooms"
		errMsg = "failed to get classrooms list from timetable"
	)

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	classrooms, err := h.service.Classrooms(r.Context())
	if err != nil {
		log.Error(errMsg, logger.Err(err))

		render.JSON(w, r, response.Error(errMsg))

		return
	}

	log.Info("classrooms list from timetable received")

	render.JSON(w, r, classroomsResponse{
		Response:   response.OK(),
		Classrooms: classrooms,
	})
}

func (h Handler) ClassroomSchedule(w http.ResponseWriter, r *http.Request) {
	const (
		op     = "http.server.handler.timetable.ClassroomSchedule"
		errMsg = "failed to get classroom events from timetable"
	)

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	classroomName := r.URL.Query().Get("name")
	if classroomName == "" {
		log.Error("invalid request: name is required")

		render.JSON(w, r, response.Error("invalid request: name is required"))

		return
	}

	log.Info("name param parsed", slog.String("name", classroomName))

	events, err := h.service.ClassroomSchedule(r.Context(), classroomName)
	if err != nil {
		log.Error(errMsg, logger.Err(err))

		render.JSON(w, r, response.Error(errMsg))

		return
	}

	log.Info("classroom events from timetable received")

	render.JSON(w, r, scheduleResponse{
		Response: response.OK(),
		Events:   events,
	})
}
