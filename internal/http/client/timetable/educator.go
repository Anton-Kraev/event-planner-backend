package timetable

import (
	"context"
	"fmt"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule/timetable"
)

const (
	educatorRoute             = api + "/educator"
	findEducatorEndpoint      = educatorRoute + "/find"
	getEducatorEventsEndpoint = educatorRoute + "/%d/events"
)

func (c Client) FindEducator(ctx context.Context, firstName, lastName, middleName string) (uint64, error) {
	const op = "http.client.timetable.FindEducator"

	url := fmt.Sprintf(findEducatorEndpoint, c.host)
	queryParams := map[string]string{
		"first_name":  firstName,
		"last_name":   lastName,
		"middle_name": middleName,
	}

	id, err := c.getID(ctx, url, queryParams)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (c Client) GetEducatorEvents(ctx context.Context, educatorID uint64) ([]timetable.Event, error) {
	const op = "http.client.timetable.GetEducatorEvents"

	url := fmt.Sprintf(getEducatorEventsEndpoint, c.host, educatorID)

	events, err := c.getEvents(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
