package timetable

import "errors"

var (
	ErrParseEducatorName    = errors.New("can't parse educator's name")
	ErrUnexpectedOwnerClass = errors.New("unexpected owner class")
	ErrNotCachedYet         = errors.New("calendar not cached yet")
)
