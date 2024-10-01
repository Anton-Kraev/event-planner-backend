package timetable

import (
	"github.com/Anton-Kraev/event-planner-backend/internal/domain/timetable"
	"github.com/Anton-Kraev/event-planner-backend/internal/lib/api/response"
)

type scheduleResponse struct {
	response.Response
	Events []timetable.Event `json:"events"`
}
