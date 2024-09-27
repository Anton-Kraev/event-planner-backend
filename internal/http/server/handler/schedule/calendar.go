package schedule

import (
	"time"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/schedule"
)

type (
	event struct {
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		Location    string    `json:"location,omitempty"`
		Description string    `json:"description,omitempty"`
	}

	calendar struct {
		Source string  `json:"source"`
		Owner  string  `json:"owner"`
		Events []event `json:"events"`
	}
)

func calendarFromDomain(d schedule.Calendar) calendar {
	var events []event

	for _, e := range d.Events {
		events = append(events, event{
			StartTime:   e.StartTime,
			EndTime:     e.EndTime,
			Location:    e.Location,
			Description: e.Description,
		})
	}

	return calendar{
		Source: string(d.Source),
		Owner:  d.Owner,
		Events: events,
	}
}
