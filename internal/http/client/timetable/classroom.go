package timetable

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

type classroomsResp struct {
	Classrooms []timetable.Classroom `json:"classrooms"`
}

const (
	classroomsRoute            = api + "/classrooms"
	getClassroomEventsEndpoint = classroomsRoute + "/%s/events"
)

func (c Client) Classrooms(ctx context.Context) ([]timetable.Classroom, error) {
	const op = "http.client.timetable.Classrooms"

	respB, err := c.doHTTP(ctx, http.MethodGet, fmt.Sprintf(classroomsRoute, c.host))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var classrooms classroomsResp

	if err = json.Unmarshal(respB, &classrooms); err != nil {
		return nil, fmt.Errorf("%s: failed to parse response body: %w", op, err)
	}

	return classrooms.Classrooms, nil
}

func (c Client) ClassroomSchedule(ctx context.Context, classroomName string) ([]timetable.Event, error) {
	const op = "http.client.timetable.ClassroomSchedule"

	url := fmt.Sprintf(getClassroomEventsEndpoint, c.host, classroomName)

	events, err := c.getEvents(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
