package timetable

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/schedule/timetable"
)

type RedisRepository struct {
	client *redis.Client
	exp    time.Duration
}

func NewRedisRepository(client *redis.Client, expirationPeriod time.Duration) RedisRepository {
	return RedisRepository{client: client, exp: expirationPeriod}
}

func (c RedisRepository) GetEvents(ctx context.Context, owner timetable.CalendarOwner) ([]timetable.Event, error) {
	key := owner.String()

	bVal, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, timetable.ErrNotCachedYet
	}

	if err != nil {
		return nil, err
	}

	var events []timetable.Event

	if err = json.Unmarshal(bVal, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (c RedisRepository) SetEvents(ctx context.Context, owner timetable.CalendarOwner, events []timetable.Event) error {
	key := owner.String()

	if c.client.Exists(ctx, key).Val() != 0 {
		return nil
	}

	bVal, err := json.Marshal(events)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, bVal, c.exp).Err()
}
