package user

import "time"

// User - структура пользователя
type User struct {
	ID         int64
	Email      string
	Name       string
	TelegramID string
	Role       string
	CreatedAt  time.Time
}
