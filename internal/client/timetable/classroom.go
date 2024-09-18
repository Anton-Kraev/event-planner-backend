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
	url := fmt.Sprintf(getClassroomEventsEndpoint, c.host, classroomName)

	return c.getEvents(ctx, url)
}
