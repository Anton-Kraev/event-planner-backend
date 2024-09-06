package timetable

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/domain/calendar"
)

type CalendarOwner struct {
	name  string
	class UserClass
}

func (o CalendarOwner) String() string {
	return fmt.Sprintf("%s %s", o.class, o.name)
}

type UserClass string

const (
	Educator  UserClass = "educator"
	Group     UserClass = "group"
	Classroom UserClass = "classroom"
)

type client interface {
	GetEducatorEvents(educatorID uint64) ([]Event, error)
	GetGroupEvents(groupID uint64) ([]Event, error)
	GetClassroomEvents(classroomName string) ([]Event, error)
	FindEducator(firstName, lastName, middleName string) (uint64, error)
	FindGroup(groupName string) (uint64, error)
}

type CalendarLoader struct {
	client client
}

var (
	errUnexpectedOwnerClass = errors.New("unexpected owner class")
	errParseEducatorName    = errors.New("can't parse educator's name")
)

func (c CalendarLoader) Load(owner CalendarOwner) (schedule calendar.Calendar, err error) {
	schedule.Owner = owner.String()
	schedule.Source = calendar.Timetable
	schedule.Events = make([]calendar.Event, 0)

	switch owner.class {
	case Educator:
		var (
			educatorID                      uint64
			events                          []Event
			firstName, lastName, middleName string
		)

		firstName, lastName, middleName, err = parseEducatorName(owner.name)
		if err != nil {
			break
		}

		educatorID, err = c.client.FindEducator(firstName, lastName, middleName)
		if err != nil {
			break
		}

		events, err = c.client.GetEducatorEvents(educatorID)
		for _, e := range events {
			schedule.Events = append(schedule.Events, e.standardize())
		}
	case Group:
		var (
			groupID uint64
			events  []Event
		)

		groupID, err = c.client.FindGroup(owner.name)
		if err != nil {
			break
		}

		events, err = c.client.GetGroupEvents(groupID)
		for _, e := range events {
			schedule.Events = append(schedule.Events, e.standardize())
		}
	case Classroom:
		var events []Event

		events, err = c.client.GetClassroomEvents(owner.name)
		if err != nil {
			break
		}

		for _, e := range events {
			schedule.Events = append(schedule.Events, e.standardize())
		}
	default:
		err = errUnexpectedOwnerClass
	}

	return
}

func parseEducatorName(name string) (string, string, string, error) {
	nameSeparated := strings.Split(strings.TrimSpace(name), " ")
	if len(nameSeparated) != 3 {
		return "", "", "", errParseEducatorName
	}

	firstName, lastName, middleName := nameSeparated[0], nameSeparated[1], nameSeparated[2]

	return firstName, lastName, middleName, nil
}

type Event struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}

func (e Event) standardize() calendar.Event {
	return calendar.Event{
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Location:    e.Location,
		Description: e.Description,
	}
}
