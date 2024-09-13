package timetable

import (
	"strings"
)

func ParseEducatorName(name string) (string, string, string, error) {
	nameSeparated := strings.Split(strings.TrimSpace(name), " ")
	if len(nameSeparated) != 3 {
		return "", "", "", ErrParseEducatorName
	}

	firstName, lastName, middleName := nameSeparated[0], nameSeparated[1], nameSeparated[2]

	return firstName, lastName, middleName, nil
}
