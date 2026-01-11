package sender

import "delayed_notifier/internal/models"

// Sender — интерфейс для всех каналов отправки
type Sender interface {
	Send(n models.Notification) error
}
