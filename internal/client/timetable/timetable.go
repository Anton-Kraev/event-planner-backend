package timetable

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/calendar/timetable"
)

const (
	host = "TODO" // TODO
	api  = "/api"

	educatorRoute  = host + api + "/educator"
	groupRoute     = host + api + "/group"
	classroomRoute = host + api + "/classroom"

	getEducatorEventsEndpoint  = educatorRoute + "/%d/events"
	getGroupEventsEndpoint     = groupRoute + "/%d/events"
	getClassroomEventsEndpoint = classroomRoute + "/%s/events"
	findEducatorEndpoint       = educatorRoute + "/find"
	findGroupEndpoint          = groupRoute + "/find"
)

type Client struct {
	client *http.Client
}

func (c Client) doHTTP(
	ctx context.Context, method string, endpoint string, queryParams map[string]string,
) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}

	for k, v := range queryParams {
		req.URL.Query().Add(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Error: %w", err)
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respB, nil
}

func (c Client) getID(ctx context.Context, endpoint string, queryParams map[string]string) (uint64, error) {
	respB, err := c.doHTTP(ctx, "HTTP", endpoint, queryParams)
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseUint(string(respB), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response body: %w", err)
	}

	return id, nil
}

func (c Client) getEvents(ctx context.Context, endpoint string) ([]timetable.Event, error) {
	respB, err := c.doHTTP(ctx, "HTTP", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var events []timetable.Event

	if err = json.Unmarshal(respB, &events); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return events, nil
}

func (c Client) GetEducatorEvents(ctx context.Context, educatorID uint64) ([]timetable.Event, error) {
	return c.getEvents(ctx, fmt.Sprintf(getEducatorEventsEndpoint, educatorID))
}

func (c Client) GetGroupEvents(ctx context.Context, groupID uint64) ([]timetable.Event, error) {
	return c.getEvents(ctx, fmt.Sprintf(getGroupEventsEndpoint, groupID))
}

func (c Client) GetClassroomEvents(ctx context.Context, classroomName string) ([]timetable.Event, error) {
	return c.getEvents(ctx, fmt.Sprintf(getClassroomEventsEndpoint, classroomName))
}

func (c Client) FindEducator(ctx context.Context, firstName, lastName, middleName string) (uint64, error) {
	queryParams := map[string]string{
		"first_name":  firstName,
		"last_name":   lastName,
		"middle_name": middleName,
	}

	return c.getID(ctx, findEducatorEndpoint, queryParams)
}

func (c Client) FindGroup(ctx context.Context, groupName string) (uint64, error) {
	return c.getID(ctx, findGroupEndpoint, map[string]string{"name": groupName})
}
