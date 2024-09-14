package schedule

type Calendar struct {
	Events []Event
	Owner  string
	Source Source
}

func NewCalendar(owner string, source Source) Calendar {
	return Calendar{
		Owner:  owner,
		Source: source,
		Events: make([]Event, 0),
	}
}
