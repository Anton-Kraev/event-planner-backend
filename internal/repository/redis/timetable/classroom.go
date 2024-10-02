package timetable

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

const classroomsKey = "classrooms"

func (c RedisRepository) GetClassrooms(ctx context.Context) ([]timetable.Classroom, error) {
	const op = "repository.redis.timetable.GetClassrooms"

	bVal, err := c.getValue(ctx, classroomsKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var classrooms []timetable.Classroom

	if err = json.Unmarshal(bVal, &classrooms); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return classrooms, nil
}

func (c RedisRepository) SetClassrooms(ctx context.Context, classrooms []timetable.Classroom) error {
	const op = "repository.redis.timetable.SetClassrooms"

	if c.exists(ctx, classroomsKey) {
		return nil
	}

	bVal, err := json.Marshal(classrooms)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = c.client.Set(ctx, classroomsKey, bVal, c.expUser).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
