package timetable

import (
	"context"
	"fmt"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/schedule/timetable"
)

const (
	groupRoute             = api + "/group"
	findGroupEndpoint      = groupRoute + "/find"
	getGroupEventsEndpoint = groupRoute + "/%d/events"
)

func (c Client) FindGroup(ctx context.Context, groupName string) (uint64, error) {
	const op = "http.client.timetable.FindGroup"

	url := fmt.Sprintf(findGroupEndpoint, c.host)
	queryParams := map[string]string{"name": groupName}

	id, err := c.getID(ctx, url, queryParams)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (c Client) GetGroupEvents(ctx context.Context, groupID uint64) ([]timetable.Event, error) {
	const op = "http.client.timetable.GetGroupEvents"

	url := fmt.Sprintf(getGroupEventsEndpoint, c.host, groupID)

	events, err := c.getEvents(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
