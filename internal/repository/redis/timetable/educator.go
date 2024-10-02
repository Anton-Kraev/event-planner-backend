package timetable

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

const educatorsKey = "educators"

func (c RedisRepository) GetEducators(ctx context.Context) ([]timetable.Educator, error) {
	const op = "repository.redis.timetable.GetEducators"

	bVal, err := c.getValue(ctx, educatorsKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var educators []timetable.Educator

	if err = json.Unmarshal(bVal, &educators); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return educators, nil
}

func (c RedisRepository) SetEducators(ctx context.Context, educators []timetable.Educator) error {
	const op = "repository.redis.timetable.SetEducators"

	if c.exists(ctx, educatorsKey) {
		return nil
	}

	bVal, err := json.Marshal(educators)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = c.client.Set(ctx, educatorsKey, bVal, c.expUser).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
