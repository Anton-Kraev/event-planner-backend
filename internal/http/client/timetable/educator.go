package timetable

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

type educatorsResp struct {
	Educators []timetable.Educator `json:"educators"`
}

const (
	educatorsRoute            = api + "/educators"
	getEducatorEventsEndpoint = educatorsRoute + "/%d/events"
)

func (c Client) Educators(ctx context.Context) ([]timetable.Educator, error) {
	const op = "http.client.timetable.Educators"

	respB, err := c.doHTTP(ctx, http.MethodGet, fmt.Sprintf(educatorsRoute, c.host))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var educators educatorsResp

	if err = json.Unmarshal(respB, &educators); err != nil {
		return nil, fmt.Errorf("%s: failed to parse response body: %w", op, err)
	}

	return educators.Educators, nil
}

func (c Client) EducatorSchedule(ctx context.Context, educatorID uint64) ([]timetable.Event, error) {
	const op = "http.client.timetable.EducatorSchedule"

	url := fmt.Sprintf(getEducatorEventsEndpoint, c.host, educatorID)

	events, err := c.getEvents(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
