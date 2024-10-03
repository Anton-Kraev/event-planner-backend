package timetable

import (
	"context"
	"errors"
	"fmt"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

func (s Service) getScheduleWithCache(
	ctx context.Context, op, key string, getScheduleFromClient func(ctx context.Context) ([]timetable.Event, error),
) ([]timetable.Event, error) {
	events, err := s.cache.GetEvents(ctx, key)
	if errors.Is(err, timetable.ErrNotCachedYet) {
		events, err = getScheduleFromClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.cache.SetEvents(ctx, key, events); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
