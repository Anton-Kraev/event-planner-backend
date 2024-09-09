package timetable

import "errors"

var (
	errUnexpectedOwnerClass = errors.New("unexpected owner class")
	errParseEducatorName    = errors.New("can't parse educator's name")
	errNotCachedYet         = errors.New("calendar not cached yet")
)
