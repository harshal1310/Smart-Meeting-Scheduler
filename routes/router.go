package routes

import (
	"smart-scheduler/handlers"

	"github.com/julienschmidt/httprouter"
)

func SetupRoutes() *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/v1/schedule", handlers.ScheduleMeeting)

	// GET routes
	router.GET("/api/v1/calendar/:userID", handlers.GetUserCalendar)

	return router
}
