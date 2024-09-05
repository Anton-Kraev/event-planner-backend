package timetable

import (
	"time"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/calendar"
)

type Client interface {
}

type CalendarLoader struct {
	client Client
}

func (c CalendarLoader) Load() calendar.Calendar {
	//TODO implement me
	panic("implement me")
}

type Event struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}
