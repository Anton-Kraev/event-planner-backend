package schedule

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule/timetable"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/lib/api/response"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/lib/logger"
)

type (
	getTimetableScheduleRequest struct {
		OwnerClass string `json:"owner_class"`
		OwnerName  string `json:"owner_name"`
	}

	getTimetableScheduleResponse struct {
		response.Response
		Calendar calendar `json:"calendar"`
	}
)

func (h Handler) GetTimetableSchedule(w http.ResponseWriter, r *http.Request) {
	const op = "http.server.handler.schedule.GetTimetableSchedule"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req getTimetableScheduleRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", logger.Err(err))

		render.JSON(w, r, response.Error("failed to decode request"))

		return
	}

	log.Info("request body decoded", slog.Any("request", req))

	if err := req.Validate(); err != nil {
		log.Error("invalid request", logger.Err(err))

		render.JSON(w, r, response.Error("invalid request"))

		return
	}

	schedule, err := h.service.GetTimetableSchedule(
		r.Context(),
		timetable.CalendarOwner{
			Class: timetable.UserClass(req.OwnerClass),
			Name:  req.OwnerName,
		},
	)
	if err != nil {
		log.Error("failed to get timetable schedule", logger.Err(err))

		render.JSON(w, r, response.Error("failed to get timetable schedule"))

		return
	}

	log.Info("timetable schedule received")

	render.JSON(w, r, getTimetableScheduleResponse{
		Response: response.OK(),
		Calendar: calendarFromDomain(schedule),
	})
}

func (r getTimetableScheduleRequest) Validate() error {
	if r.OwnerClass == "" {
		return errors.New("field owner_class is empty")
	}

	if r.OwnerName == "" {
		return errors.New("field owner_name is empty")
	}

	return nil
}
