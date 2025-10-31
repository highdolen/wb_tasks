package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"secondBlock/L2.18/internal/calendar"
	"secondBlock/L2.18/internal/models"

	"github.com/gorilla/mux"
)

type Handlers struct {
	calendar *calendar.Calendar
}

func NewHandlers(cal *calendar.Calendar) *Handlers {
	return &Handlers{
		calendar: cal,
	}
}

func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEventRequest
	if err := decodeRequest(r, &req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.Date == "" || req.Title == "" {
		sendError(w, "user_id, date and title are required", http.StatusBadRequest)
		return
	}

	event, err := h.calendar.CreateEvent(req.UserID, req.Date, req.Title)
	if err != nil {
		status := http.StatusBadRequest
		if err == calendar.ErrInvalidDate {
			status = http.StatusBadRequest
		}
		sendError(w, err.Error(), status)
		return
	}

	sendSuccess(w, "event created with id: "+strconv.Itoa(event.ID))
}

func (h *Handlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateEventRequest
	if err := decodeRequest(r, &req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == 0 || req.UserID == 0 || req.Date == "" || req.Title == "" {
		sendError(w, "id, user_id, date and title are required", http.StatusBadRequest)
		return
	}

	event, err := h.calendar.UpdateEvent(req.ID, req.UserID, req.Date, req.Title)
	if err != nil {
		status := http.StatusBadRequest
		if err == calendar.ErrEventNotFound {
			status = http.StatusServiceUnavailable
		} else if err == calendar.ErrInvalidDate {
			status = http.StatusBadRequest
		}
		sendError(w, err.Error(), status)
		return
	}

	sendSuccess(w, "event updated: "+strconv.Itoa(event.ID))
}

func (h *Handlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteEventRequest
	if err := decodeRequest(r, &req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == 0 || req.UserID == 0 {
		sendError(w, "id and user_id are required", http.StatusBadRequest)
		return
	}

	err := h.calendar.DeleteEvent(req.ID, req.UserID)
	if err != nil {
		status := http.StatusBadRequest
		if err == calendar.ErrEventNotFound {
			status = http.StatusServiceUnavailable
		}
		sendError(w, err.Error(), status)
		return
	}

	sendSuccess(w, "event deleted")
}

func (h *Handlers) EventsForDay(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseQueryParams(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.calendar.GetEventsForDay(userID, date)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendJSON(w, models.EventsResponse{Result: events})
}

func (h *Handlers) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseQueryParams(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.calendar.GetEventsForWeek(userID, date)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendJSON(w, models.EventsResponse{Result: events})
}

func (h *Handlers) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseQueryParams(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.calendar.GetEventsForMonth(userID, date)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendJSON(w, models.EventsResponse{Result: events})
}

func decodeRequest(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		return json.NewDecoder(r.Body).Decode(v)
	}

	// Default to form data
	if err := r.ParseForm(); err != nil {
		return err
	}

	// Simple form to struct mapping
	if req, ok := v.(*models.CreateEventRequest); ok {
		if idStr := r.FormValue("user_id"); idStr != "" {
			req.UserID, _ = strconv.Atoi(idStr)
		}
		req.Date = r.FormValue("date")
		req.Title = r.FormValue("event")
	}

	if req, ok := v.(*models.UpdateEventRequest); ok {
		if idStr := r.FormValue("id"); idStr != "" {
			req.ID, _ = strconv.Atoi(idStr)
		}
		if idStr := r.FormValue("user_id"); idStr != "" {
			req.UserID, _ = strconv.Atoi(idStr)
		}
		req.Date = r.FormValue("date")
		req.Title = r.FormValue("event")
	}

	if req, ok := v.(*models.DeleteEventRequest); ok {
		if idStr := r.FormValue("id"); idStr != "" {
			req.ID, _ = strconv.Atoi(idStr)
		}
		if idStr := r.FormValue("user_id"); idStr != "" {
			req.UserID, _ = strconv.Atoi(idStr)
		}
	}

	return nil
}

func parseQueryParams(r *http.Request) (int, string, error) {
	userIDStr := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	if userIDStr == "" || date == "" {
		return 0, "", calendar.ErrInvalidDate
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, "", calendar.ErrInvalidDate
	}

	return userID, date, nil
}

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendSuccess(w http.ResponseWriter, message string) {
	sendJSON(w, models.SuccessResponse{Result: message})
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}

func (h *Handlers) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/create_event", h.CreateEvent).Methods("POST")
	router.HandleFunc("/update_event", h.UpdateEvent).Methods("POST")
	router.HandleFunc("/delete_event", h.DeleteEvent).Methods("POST")
	router.HandleFunc("/events_for_day", h.EventsForDay).Methods("GET")
	router.HandleFunc("/events_for_week", h.EventsForWeek).Methods("GET")
	router.HandleFunc("/events_for_month", h.EventsForMonth).Methods("GET")

	return router
}
