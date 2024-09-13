package timetable

import "fmt"

type UserClass string

const (
	Educator  UserClass = "educator"
	Group     UserClass = "group"
	Classroom UserClass = "classroom"
)

type CalendarOwner struct {
	Name  string
	Class UserClass
}

func (o CalendarOwner) String() string {
	return fmt.Sprintf("%s %s", o.Class, o.Name)
}
