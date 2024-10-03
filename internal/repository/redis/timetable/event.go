package timetable

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

func (c RedisRepository) GetEvents(ctx context.Context, ownerKey string) ([]timetable.Event, error) {
	const op = "repository.redis.timetable.GetEvents"

	bVal, err := c.getValue(ctx, ownerKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var events []timetable.Event

	if err = json.Unmarshal(bVal, &events); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

func (c RedisRepository) SetEvents(ctx context.Context, ownerKey string, events []timetable.Event) error {
	const op = "repository.redis.timetable.SetEvents"

	if c.exists(ctx, ownerKey) {
		return nil
	}

	bVal, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = c.client.Set(ctx, ownerKey, bVal, c.expEvent).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
