package timetable

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
	"github.com/Anton-Kraev/event-planner-backend/internal/lib/api/helpers"
	"github.com/Anton-Kraev/event-planner-backend/internal/lib/api/response"
	"github.com/Anton-Kraev/event-planner-backend/internal/lib/logger"
)

type educatorsResponse struct {
	response.Response
	Educators []timetable.Educator `json:"educators"`
}

func (h Handler) Educators(w http.ResponseWriter, r *http.Request) {
	const (
		op     = "http.server.handler.timetable.Educators"
		errMsg = "failed to get educators list from timetable"
	)

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	educators, err := h.service.Educators(r.Context())
	if err != nil {
		log.Error(errMsg, logger.Err(err))

		render.JSON(w, r, response.Error(errMsg))

		return
	}

	log.Info("educators list from timetable received")

	render.JSON(w, r, educatorsResponse{
		Response:  response.OK(),
		Educators: educators,
	})
}

func (h Handler) EducatorSchedule(w http.ResponseWriter, r *http.Request) {
	const (
		op     = "http.server.handler.timetable.EducatorSchedule"
		errMsg = "failed to get educator events from timetable"
	)

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	educatorID, err := helpers.ParseAndValidateID(chi.URLParam(r, "id"))
	if err != nil {
		log.Error("invalid request", logger.Err(err))

		render.JSON(w, r, response.Error("invalid request"))

		return
	}

	log.Info("id param parsed", slog.Uint64("id", educatorID))

	events, err := h.service.EducatorSchedule(r.Context(), educatorID)
	if err != nil {
		log.Error(errMsg, logger.Err(err))

		render.JSON(w, r, response.Error(errMsg))

		return
	}

	log.Info("educator events from timetable received")

	render.JSON(w, r, scheduleResponse{
		Response: response.OK(),
		Events:   events,
	})
}
