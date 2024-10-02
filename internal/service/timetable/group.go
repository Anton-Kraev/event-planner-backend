package timetable

import (
	"context"
	"errors"
	"fmt"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

func (s Service) Groups(ctx context.Context) ([]timetable.Group, error) {
	const op = "service.timetable.Groups"

	groups, err := s.cache.GetGroups(ctx)
	if errors.Is(err, timetable.ErrNotCachedYet) {
		groups, err = s.client.Groups(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.cache.SetGroups(ctx, groups); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return groups, nil
}

func (s Service) GroupSchedule(ctx context.Context, id uint64) ([]timetable.Event, error) {
	const op = "service.timetable.GroupSchedule"

	var key = fmt.Sprintf("group:%d", id)

	return s.getScheduleWithCache(ctx, op, key, func(ctx context.Context) ([]timetable.Event, error) {
		return s.client.GroupSchedule(ctx, id)
	})
}
