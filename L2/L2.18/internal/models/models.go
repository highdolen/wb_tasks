package models

import (
	"time"
)

type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Title  string    `json:"title"`
}

type CreateEventRequest struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Title  string `json:"title"`
}

type UpdateEventRequest struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Title  string `json:"title"`
}

type DeleteEventRequest struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
}

type EventsResponse struct {
	Result []Event `json:"result"`
}

type SuccessResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
