package timetable

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/schedule/timetable"
)

type RedisRepository struct {
	client *redis.Client
	exp    time.Duration
}

func NewRedisRepository(client *redis.Client, expirationPeriod time.Duration) RedisRepository {
	return RedisRepository{client: client, exp: expirationPeriod}
}

func (c RedisRepository) GetEvents(ctx context.Context, owner timetable.CalendarOwner) ([]timetable.Event, error) {
	const op = "repository.redis.timetable.GetEvents"

	key := owner.String()

	bVal, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("%s: %w", op, timetable.ErrNotCachedYet)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var events []timetable.Event

	if err = json.Unmarshal(bVal, &events); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

func (c RedisRepository) SetEvents(ctx context.Context, owner timetable.CalendarOwner, events []timetable.Event) error {
	const op = "repository.redis.timetable.SetEvents"

	key := owner.String()

	if c.client.Exists(ctx, key).Val() != 0 {
		return nil
	}

	bVal, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = c.client.Set(ctx, key, bVal, c.exp).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
