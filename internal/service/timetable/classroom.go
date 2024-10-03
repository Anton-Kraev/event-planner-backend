package timetable

import (
	"context"
	"errors"
	"fmt"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
)

func (s Service) Classrooms(ctx context.Context) ([]timetable.Classroom, error) {
	const op = "service.timetable.Classrooms"

	classrooms, err := s.cache.GetClassrooms(ctx)
	if errors.Is(err, timetable.ErrNotCachedYet) {
		classrooms, err = s.client.Classrooms(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.cache.SetClassrooms(ctx, classrooms); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return classrooms, nil
}

func (s Service) ClassroomSchedule(ctx context.Context, name string) ([]timetable.Event, error) {
	const op = "service.timetable.ClassroomSchedule"

	var key = fmt.Sprintf("classroom:%s", name)

	return s.getScheduleWithCache(ctx, op, key, func(ctx context.Context) ([]timetable.Event, error) {
		return s.client.ClassroomSchedule(ctx, name)
	})
}
