package timetable

import (
	"time"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/calendar"
)

type Event struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}

func (e Event) Standardize() calendar.Event {
	return calendar.Event{
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Location:    e.Location,
		Description: e.Description,
	}
}
