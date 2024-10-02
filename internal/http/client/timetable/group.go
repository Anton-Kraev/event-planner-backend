package timetable

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

const (
	groupsRoute            = api + "/groups"
	getGroupEventsEndpoint = groupsRoute + "/%d/events"
)

func (c Client) Groups(ctx context.Context) ([]timetable.Group, error) {
	const op = "http.client.timetable.Groups"

	respB, err := c.doHTTP(ctx, http.MethodGet, groupsRoute)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var groups []timetable.Group

	if err = json.Unmarshal(respB, &groups); err != nil {
		return nil, fmt.Errorf("%s: failed to parse response body: %w", op, err)
	}

	return groups, nil
}

func (c Client) GroupSchedule(ctx context.Context, groupID uint64) ([]timetable.Event, error) {
	const op = "http.client.timetable.GroupSchedule"

	url := fmt.Sprintf(getGroupEventsEndpoint, c.host, groupID)

	events, err := c.getEvents(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
