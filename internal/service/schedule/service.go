//go:generate mockgen -package=mocks -source=./service.go -destination=./mocks/service_mock.go

package schedule

import (
	"context"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/schedule/timetable"
)

type (
	timetableClient interface {
		GetEducatorEvents(ctx context.Context, educatorID uint64) ([]timetable.Event, error)
		GetGroupEvents(ctx context.Context, groupID uint64) ([]timetable.Event, error)
		GetClassroomEvents(ctx context.Context, classroomName string) ([]timetable.Event, error)
		FindEducator(ctx context.Context, firstName, lastName, middleName string) (uint64, error)
		FindGroup(ctx context.Context, groupName string) (uint64, error)
	}

	timetableCache interface {
		GetEvents(ctx context.Context, owner timetable.CalendarOwner) ([]timetable.Event, error)
		SetEvents(ctx context.Context, owner timetable.CalendarOwner, events []timetable.Event) error
	}

	Service struct {
		ttClient timetableClient
		ttCache  timetableCache
	}
)

func NewService(ttClient timetableClient, ttCache timetableCache) Service {
	return Service{ttClient: ttClient, ttCache: ttCache}
}
