package schedule

import (
	"context"
	"errors"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule"
	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule/timetable"
)

func (s Service) GetTimetableSchedule(
	ctx context.Context, owner timetable.CalendarOwner,
) (schedule.Calendar, error) {
	events, err := s.ttCache.GetEvents(ctx, owner)
	if errors.Is(err, timetable.ErrNotCachedYet) {
		events, err = s.getEventsFromTimetable(ctx, owner)
	}

	if err != nil {
		return schedule.Calendar{}, err
	}

	if err = s.ttCache.SetEvents(ctx, owner, events); err != nil {
		return schedule.Calendar{}, err
	}

	calendar := schedule.NewCalendar(owner.String(), schedule.Timetable)

	for _, event := range events {
		calendar.Events = append(calendar.Events, event.Standardize())
	}

	return calendar, nil
}

func (s Service) getEventsFromTimetable(
	ctx context.Context, owner timetable.CalendarOwner,
) (events []timetable.Event, err error) {
	switch owner.Class {
	case timetable.Educator:
		var (
			firstName, lastName, middleName string
			educatorID                      uint64
		)

		firstName, lastName, middleName, err = timetable.ParseEducatorName(owner.Name)
		if err != nil {
			return
		}

		educatorID, err = s.ttClient.FindEducator(ctx, firstName, lastName, middleName)
		if err != nil {
			return
		}

		events, err = s.ttClient.GetEducatorEvents(ctx, educatorID)
	case timetable.Group:
		var groupID uint64

		groupID, err = s.ttClient.FindGroup(ctx, owner.Name)
		if err != nil {
			return
		}

		events, err = s.ttClient.GetGroupEvents(ctx, groupID)
	case timetable.Classroom:
		events, err = s.ttClient.GetClassroomEvents(ctx, owner.Name)
	default:
		err = timetable.ErrUnexpectedOwnerClass
	}

	return
}
