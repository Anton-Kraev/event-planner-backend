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

type groupsResponse struct {
	response.Response
	Groups []timetable.Group `json:"groups"`
}

func (h Handler) Groups(w http.ResponseWriter, r *http.Request) {
	const (
		op     = "http.server.handler.timetable.Groups"
		errMsg = "failed to get groups list from timetable"
	)

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	groups, err := h.service.Groups(r.Context())
	if err != nil {
		log.Error(errMsg, logger.Err(err))

		render.JSON(w, r, response.Error(errMsg))

		return
	}

	log.Info("groups list from timetable received")

	render.JSON(w, r, groupsResponse{
		Response: response.OK(),
		Groups:   groups,
	})
}

func (h Handler) GroupSchedule(w http.ResponseWriter, r *http.Request) {
	const (
		op     = "http.server.handler.timetable.GroupSchedule"
		errMsg = "failed to get group events from timetable"
	)

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	groupID, err := helpers.ParseAndValidateID(chi.URLParam(r, "id"))
	if err != nil {
		log.Error("invalid request", logger.Err(err))

		render.JSON(w, r, response.Error("invalid request"))

		return
	}

	log.Info("id param parsed")

	events, err := h.service.GroupSchedule(r.Context(), groupID)
	if err != nil {
		log.Error(errMsg, logger.Err(err))

		render.JSON(w, r, response.Error(errMsg))

		return
	}

	log.Info("group events from timetable received")

	render.JSON(w, r, scheduleResponse{
		Response: response.OK(),
		Events:   events,
	})
}
