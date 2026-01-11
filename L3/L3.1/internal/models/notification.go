package models

import "time"

type NotificationStatus string

const (
	StatusPending   NotificationStatus = "pending"
	StatusSent      NotificationStatus = "sent"
	StatusFailed    NotificationStatus = "failed"
	StatusCancelled NotificationStatus = "cancelled"
)

type Notification struct {
	ID        string    `json:"id"`
	SendAt    time.Time `json:"send_at"`
	Channel   string    `json:"channel"`
	Recipient string    `json:"recipient"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`

	Status NotificationStatus `json:"status"`
	Meta   NotificationMeta   `json:"meta"`
}

type NotificationMeta struct {
	Attempt    int       `json:"attempt"`
	MaxAttempt int       `json:"max_attempts"`
	CreatedAt  time.Time `json:"created_at"`
}
