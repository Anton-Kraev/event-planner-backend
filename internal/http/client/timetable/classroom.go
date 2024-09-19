package timetable

import (
	"context"
	"fmt"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule/timetable"
)

const (
	classroomRoute             = api + "/classroom"
	getClassroomEventsEndpoint = classroomRoute + "/%s/events"
)

func (c Client) GetClassroomEvents(ctx context.Context, classroomName string) ([]timetable.Event, error) {
	const op = "http.client.timetable.GetClassroomEvents"

	url := fmt.Sprintf(getClassroomEventsEndpoint, c.host, classroomName)

	events, err := c.getEvents(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
