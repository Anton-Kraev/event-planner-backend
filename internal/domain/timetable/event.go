package timetable

type Event struct {
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Location    string `json:"location"`
	Description string `json:"description"`
}
