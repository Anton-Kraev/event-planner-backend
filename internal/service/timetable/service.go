package timetable

import (
	"context"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

type (
	client interface {
		Educators(ctx context.Context) ([]timetable.Educator, error)
		Groups(ctx context.Context) ([]timetable.Group, error)
		Classrooms(ctx context.Context) ([]timetable.Classroom, error)
		EducatorSchedule(ctx context.Context, id uint64) ([]timetable.Event, error)
		GroupSchedule(ctx context.Context, id uint64) ([]timetable.Event, error)
		ClassroomSchedule(ctx context.Context, name string) ([]timetable.Event, error)
	}

	cache interface {
		GetEducators(ctx context.Context) ([]timetable.Educator, error)
		SetEducators(ctx context.Context, educators []timetable.Educator) error
		GetGroups(ctx context.Context) ([]timetable.Group, error)
		SetGroups(ctx context.Context, groups []timetable.Group) error
		GetClassrooms(ctx context.Context) ([]timetable.Classroom, error)
		SetClassrooms(ctx context.Context, classrooms []timetable.Classroom) error
		GetEvents(ctx context.Context, ownerKey string) ([]timetable.Event, error)
		SetEvents(ctx context.Context, ownerKey string, events []timetable.Event) error
	}

	Service struct {
		client client
		cache  cache
	}
)

func NewService(timetableClient client, timetableCache cache) Service {
	return Service{client: timetableClient, cache: timetableCache}
}
