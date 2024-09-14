package app

import (
	"log"

	"github.com/Anton-Kraev/event-timeslot-planner/internal/config"
)

func Run() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	// TODO: init logger: slog
	// TODO: init tt cache
	// TODO: init schedule service
	// TODO: init schedule controller
	// TODO: init router: chi/stdlib
	// TODO: run server
}
