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
	url := fmt.Sprintf(findEducatorEndpoint, c.host)
	queryParams := map[string]string{
		"first_name":  firstName,
		"last_name":   lastName,
		"middle_name": middleName,
	}

	return c.getID(ctx, url, queryParams)
}

func (c Client) GetEducatorEvents(ctx context.Context, educatorID uint64) ([]timetable.Event, error) {
	url := fmt.Sprintf(getEducatorEventsEndpoint, c.host, educatorID)

	return c.getEvents(ctx, url)
}
