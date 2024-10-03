package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "github.com/Anton-Kraev/event-planner-backend/internal/http/server/middleware"
)

type timetableHandler interface {
	Educators(w http.ResponseWriter, r *http.Request)
	EducatorSchedule(w http.ResponseWriter, r *http.Request)
	Groups(w http.ResponseWriter, r *http.Request)
	GroupSchedule(w http.ResponseWriter, r *http.Request)
	Classrooms(w http.ResponseWriter, r *http.Request)
	ClassroomSchedule(w http.ResponseWriter, r *http.Request)
}

func SetupRouter(log *slog.Logger, timetableHandler timetableHandler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mw.NewLogger(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/api/timetable", func(r chi.Router) {
		r.Route("/educators", func(r chi.Router) {
			r.Get("/", timetableHandler.Educators)
			r.Get("/{id}", timetableHandler.EducatorSchedule)
		})

		r.Route("/groups", func(r chi.Router) {
			r.Get("/", timetableHandler.Groups)
			r.Get("/{id}", timetableHandler.GroupSchedule)
		})

		r.Get("/classrooms", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("name") == "" {
				timetableHandler.Classrooms(w, r)
			} else {
				timetableHandler.ClassroomSchedule(w, r)
			}
		})
	})

	return router
}
