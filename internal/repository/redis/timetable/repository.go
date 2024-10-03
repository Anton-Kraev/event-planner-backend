package timetable

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

type RedisRepository struct {
	client   *redis.Client
	expUser  time.Duration
	expEvent time.Duration
}

func NewRedisRepository(client *redis.Client, userExpPeriod, eventExpPeriod time.Duration) RedisRepository {
	return RedisRepository{client: client, expUser: userExpPeriod, expEvent: eventExpPeriod}
}

func (c RedisRepository) exists(ctx context.Context, key string) bool {
	return c.client.Exists(ctx, key).Val() != 0
}

func (c RedisRepository) getValue(ctx context.Context, key string) ([]byte, error) {
	bVal, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, timetable.ErrNotCachedYet
	}

	if err != nil {
		return nil, err
	}

	return bVal, nil
}
