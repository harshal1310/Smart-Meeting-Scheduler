package handlers

import (
	"encoding/json"
	"net/http"
	"smart-scheduler/api"
	"smart-scheduler/repository"
	service "smart-scheduler/service"

	"github.com/julienschmidt/httprouter"
)

func ScheduleMeeting(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req repository.ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}
	resp, err := service.ScheduleEvent(req)
	if err != nil {
		api.Error(w, r, err, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	api.SuccessJson(w, r, resp)
}

func GetUserCalendar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract userID from httprouter params
	userId := ps.ByName("userID")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	events, err := service.GetCalendarEvents(userId, start, end)
	if err != nil {
		api.Error(w, r, err, http.StatusInternalServerError)
		return
	}

	api.SuccessJson(w, r, events)
}
