package timetable

import (
	"context"
	"errors"
	"fmt"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

func (s Service) Educators(ctx context.Context) ([]timetable.Educator, error) {
	const op = "service.timetable.Educators"

	educators, err := s.cache.GetEducators(ctx)
	if errors.Is(err, timetable.ErrNotCachedYet) {
		educators, err = s.client.Educators(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.cache.SetEducators(ctx, educators); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return educators, nil
}

func (s Service) EducatorSchedule(ctx context.Context, id uint64) ([]timetable.Event, error) {
	const op = "service.timetable.EducatorSchedule"

	var key = fmt.Sprintf("educator:%d", id)

	return s.getScheduleWithCache(ctx, op, key, func(ctx context.Context) ([]timetable.Event, error) {
		return s.client.EducatorSchedule(ctx, id)
	})
}
