package timetable

import (
	"time"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/schedule"
)

type Event struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}

func (e Event) Standardize() schedule.Event {
	return schedule.Event{
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Location:    e.Location,
		Description: e.Description,
	}
}
