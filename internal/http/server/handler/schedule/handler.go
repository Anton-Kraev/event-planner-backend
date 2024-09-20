package schedule

import (
	"context"
	"log/slog"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule/timetable"
)

type (
	service interface {
		GetTimetableSchedule(ctx context.Context, owner timetable.CalendarOwner) (schedule.Calendar, error)
	}

	Handler struct {
		service service
		log     *slog.Logger
	}
)

func NewHandler(scheduleService service, logger *slog.Logger) Handler {
	return Handler{service: scheduleService, log: logger}
}
