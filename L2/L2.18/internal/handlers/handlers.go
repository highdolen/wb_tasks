package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"secondBlock/L2.18/internal/calendar"
	"secondBlock/L2.18/internal/models"

	"github.com/gorilla/mux"
)

// Handlers предоставляет HTTP-обработчики для календаря.
type Handlers struct {
	calendar *calendar.Calendar
}

// NewHandlers - конструктор для Handlers.
func NewHandlers(cal *calendar.Calendar) *Handlers {
	return &Handlers{
		calendar: cal,
	}
}

// CreateEvent - хэндлер для создания события.
func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEventRequest

	if err := decodeRequest(w, r, &req); err != nil {
		sendError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.Date == "" || req.Title == "" {
		sendError(w, "user_id, date and title are required", http.StatusBadRequest)
		return
	}

	event, err := h.calendar.CreateEvent(req.UserID, req.Date, req.Title)
	if err != nil {
		switch {
		case errors.Is(err, calendar.ErrInvalidDate):
			sendError(w, err.Error(), http.StatusBadRequest)
		default:
			sendError(w, err.Error(), http.StatusServiceUnavailable)
		}
		return
	}

	sendSuccess(w, "event created with id: "+strconv.Itoa(event.ID))
}

// UpdateEvent - хэндлер для обновления события.
func (h *Handlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateEventRequest
	if err := decodeRequest(w, r, &req); err != nil {
		sendError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ID == 0 || req.UserID == 0 || req.Date == "" || req.Title == "" {
		sendError(w, "id, user_id, date and title are required", http.StatusBadRequest)
		return
	}

	event, err := h.calendar.UpdateEvent(req.ID, req.UserID, req.Date, req.Title)
	if err != nil {
		switch {
		case errors.Is(err, calendar.ErrInvalidDate):
			sendError(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, calendar.ErrEventNotFound):
			sendError(w, err.Error(), http.StatusNotFound)
		default:
			sendError(w, err.Error(), http.StatusServiceUnavailable)
		}
		return
	}

	sendSuccess(w, "event updated: "+strconv.Itoa(event.ID))
}

// DeleteEvent - хэндлер для удаления события.
func (h *Handlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteEventRequest
	if err := decodeRequest(w, r, &req); err != nil {
		sendError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ID == 0 || req.UserID == 0 {
		sendError(w, "id and user_id are required", http.StatusBadRequest)
		return
	}

	if err := h.calendar.DeleteEvent(req.ID, req.UserID); err != nil {
		if errors.Is(err, calendar.ErrEventNotFound) {
			sendError(w, err.Error(), http.StatusNotFound)
			return
		}
		sendError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	sendSuccess(w, "event deleted")
}

// EventsForDay - возвращает события за день.
func (h *Handlers) EventsForDay(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseQueryParams(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.calendar.GetEventsForDay(userID, date)
	if err != nil {
		if errors.Is(err, calendar.ErrInvalidDate) {
			sendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		sendError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	sendJSON(w, models.EventsResponse{Result: events})
}

// EventsForWeek - возвращает события за неделю.
func (h *Handlers) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseQueryParams(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.calendar.GetEventsForWeek(userID, date)
	if err != nil {
		if errors.Is(err, calendar.ErrInvalidDate) {
			sendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		sendError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	sendJSON(w, models.EventsResponse{Result: events})
}

// EventsForMonth - возвращает события за месяц.
func (h *Handlers) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseQueryParams(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.calendar.GetEventsForMonth(userID, date)
	if err != nil {
		if errors.Is(err, calendar.ErrInvalidDate) {
			sendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		sendError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	sendJSON(w, models.EventsResponse{Result: events})
}

// decodeRequest читает тело запроса в структуру v.
// Поддерживает JSON и form-urlencoded.
// Ограничивает размер тела и корректно закрывает body.
func decodeRequest(w http.ResponseWriter, r *http.Request, v interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		defer func() {
			_, _ = io.Copy(io.Discard, r.Body)
			_ = r.Body.Close()
		}()

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		return dec.Decode(v)
	}

	// form-urlencoded
	if err := r.ParseForm(); err != nil {
		return err
	}

	// Маппинг form → структура
	switch req := v.(type) {
	case *models.CreateEventRequest:
		req.UserID, _ = strconv.Atoi(r.FormValue("user_id"))
		req.Date = r.FormValue("date")
		req.Title = r.FormValue("event")

	case *models.UpdateEventRequest:
		req.ID, _ = strconv.Atoi(r.FormValue("id"))
		req.UserID, _ = strconv.Atoi(r.FormValue("user_id"))
		req.Date = r.FormValue("date")
		req.Title = r.FormValue("event")

	case *models.DeleteEventRequest:
		req.ID, _ = strconv.Atoi(r.FormValue("id"))
		req.UserID, _ = strconv.Atoi(r.FormValue("user_id"))
	}

	return nil
}

func parseQueryParams(r *http.Request) (int, string, error) {
	userIDStr := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	if userIDStr == "" {
		return 0, "", errors.New("missing user_id")
	}
	if date == "" {
		return 0, "", errors.New("missing date")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, "", errors.New("invalid user_id")
	}

	return userID, date, nil
}

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func sendSuccess(w http.ResponseWriter, message string) {
	sendJSON(w, models.SuccessResponse{Result: message})
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}

// SetupRoutes регистрирует маршруты.
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
