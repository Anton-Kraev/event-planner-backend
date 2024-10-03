package timetable

import (
	"context"
	"log/slog"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

type (
	service interface {
		Educators(ctx context.Context) ([]timetable.Educator, error)
		Groups(ctx context.Context) ([]timetable.Group, error)
		Classrooms(ctx context.Context) ([]timetable.Classroom, error)
		EducatorSchedule(ctx context.Context, id uint64) ([]timetable.Event, error)
		GroupSchedule(ctx context.Context, id uint64) ([]timetable.Event, error)
		ClassroomSchedule(ctx context.Context, name string) ([]timetable.Event, error)
	}

	Handler struct {
		service service
		log     *slog.Logger
	}
)

func NewHandler(timetableService service, logger *slog.Logger) Handler {
	return Handler{service: timetableService, log: logger}
}
