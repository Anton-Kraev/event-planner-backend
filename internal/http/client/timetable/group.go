package timetable

import (
	"context"
	"fmt"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule/timetable"
)

const (
	groupRoute             = api + "/group"
	findGroupEndpoint      = groupRoute + "/find"
	getGroupEventsEndpoint = groupRoute + "/%d/events"
)

func (c Client) FindGroup(ctx context.Context, groupName string) (uint64, error) {
	url := fmt.Sprintf(findGroupEndpoint, c.host)
	queryParams := map[string]string{"name": groupName}

	return c.getID(ctx, url, queryParams)
}

func (c Client) GetGroupEvents(ctx context.Context, groupID uint64) ([]timetable.Event, error) {
	url := fmt.Sprintf(getGroupEventsEndpoint, c.host, groupID)

	return c.getEvents(ctx, url)
}
