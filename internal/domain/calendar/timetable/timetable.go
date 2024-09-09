package timetable

import (
	"context"
	"errors"
	"strings"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/calendar"
)

type (
	client interface {
		GetEducatorEvents(ctx context.Context, educatorID uint64) ([]Event, error)
		GetGroupEvents(ctx context.Context, groupID uint64) ([]Event, error)
		GetClassroomEvents(ctx context.Context, classroomName string) ([]Event, error)
		FindEducator(ctx context.Context, firstName, lastName, middleName string) (uint64, error)
		FindGroup(ctx context.Context, groupName string) (uint64, error)
	}

	cache interface {
		GetEvents(ctx context.Context, owner CalendarOwner) ([]Event, error)
		SaveEvents(ctx context.Context, owner CalendarOwner, events []Event) error
	}

	CalendarLoader struct {
		client client
		cache  cache
	}
)

func NewCalendarLoader(client client, cache cache) CalendarLoader {
	return CalendarLoader{client: client, cache: cache}
}

func (c CalendarLoader) Load(ctx context.Context, owner CalendarOwner) (calendar.Calendar, error) {
	events, err := c.cache.GetEvents(ctx, owner)
	if err != nil && !errors.Is(err, errNotCachedYet) {
		return calendar.Calendar{}, err
	}

	if errors.Is(err, errNotCachedYet) {
		events, err = c.getEventsFromTimetable(ctx, owner)
		if err != nil {
			return calendar.Calendar{}, err
		}
	}

	schedule := calendar.NewCalendar(owner.String(), calendar.Timetable)

	for _, event := range events {
		schedule.Events = append(schedule.Events, event.standardize())
	}

	return schedule, nil
}

func (c CalendarLoader) getEventsFromTimetable(ctx context.Context, owner CalendarOwner) (events []Event, err error) {
	switch owner.class {
	case Educator:
		var (
			firstName, lastName, middleName string
			educatorID                      uint64
		)

		firstName, lastName, middleName, err = parseEducatorName(owner.name)
		if err != nil {
			return
		}

		educatorID, err = c.client.FindEducator(ctx, firstName, lastName, middleName)
		if err != nil {
			return
		}

		events, err = c.client.GetEducatorEvents(ctx, educatorID)
	case Group:
		var groupID uint64

		groupID, err = c.client.FindGroup(ctx, owner.name)
		if err != nil {
			return
		}

		events, err = c.client.GetGroupEvents(ctx, groupID)
	case Classroom:
		events, err = c.client.GetClassroomEvents(ctx, owner.name)
	default:
		err = errUnexpectedOwnerClass
	}

	return
}

func parseEducatorName(name string) (string, string, string, error) {
	nameSeparated := strings.Split(strings.TrimSpace(name), " ")
	if len(nameSeparated) != 3 {
		return "", "", "", errParseEducatorName
	}

	firstName, lastName, middleName := nameSeparated[0], nameSeparated[1], nameSeparated[2]

	return firstName, lastName, middleName, nil
}
