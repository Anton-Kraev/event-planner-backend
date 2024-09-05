package calendar

import (
	"time"
)

type Source string

const (
	Timetable Source = "timetable"
)

type Event struct {
	StartTime   time.Time
	EndTime     time.Time
	Location    string
	Description string
}

type Calendar struct {
	Events []Event
	Owner  string
	Source Source
}

type Loader interface {
	Load() Calendar
}
