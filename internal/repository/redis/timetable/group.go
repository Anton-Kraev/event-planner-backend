package timetable

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

const groupsKey = "groups"

func (c RedisRepository) GetGroups(ctx context.Context) ([]timetable.Group, error) {
	const op = "repository.redis.timetable.GetGroups"

	bVal, err := c.getValue(ctx, groupsKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var groups []timetable.Group

	if err = json.Unmarshal(bVal, &groups); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return groups, nil
}

func (c RedisRepository) SetGroups(ctx context.Context, groups []timetable.Group) error {
	const op = "repository.redis.timetable.SetGroups"

	if c.exists(ctx, groupsKey) {
		return nil
	}

	bVal, err := json.Marshal(groups)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = c.client.Set(ctx, groupsKey, bVal, c.expUser).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
