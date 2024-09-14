package schedule

import "time"

type Event struct {
	StartTime   time.Time
	EndTime     time.Time
	Location    string
	Description string
}
