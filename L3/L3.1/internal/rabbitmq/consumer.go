package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"delayed_notifier/internal/models"
	"delayed_notifier/internal/sender"
	"delayed_notifier/internal/storage"

	"github.com/rabbitmq/amqp091-go"
	wbfrabbit "github.com/wb-go/wbf/rabbitmq"
)

type NotificationConsumer struct {
	store          storage.NotificationStorage
	emailSender    sender.Sender
	telegramSender sender.Sender
}

func NewNotificationConsumer(
	client *wbfrabbit.RabbitClient,
	store storage.NotificationStorage,
	emailSender sender.Sender,
	telegramSender sender.Sender,
) *wbfrabbit.Consumer {

	handler := &NotificationConsumer{
		store:          store,
		emailSender:    emailSender,
		telegramSender: telegramSender,
	}

	cfg := wbfrabbit.ConsumerConfig{
		Queue:         "notifications_queue",
		ConsumerTag:   "notifications-consumer",
		AutoAck:       false,
		Workers:       3,
		PrefetchCount: 5,
	}

	return wbfrabbit.NewConsumer(client, cfg, handler.Handle)
}

func (c *NotificationConsumer) Handle(ctx context.Context, msg amqp091.Delivery) error {
	// Парсинг сообщение
	var n models.Notification
	if err := json.Unmarshal(msg.Body, &n); err != nil {
		log.Printf("consumer: failed to unmarshal message, err=%v", err)
		return err
	}

	log.Printf(
		"consumer: notification received id=%s channel=%s send_at=%s",
		n.ID,
		n.Channel,
		n.SendAt.Format(time.RFC3339),
	)

	// Первая проверка статуса
	if c.isCancelled(n.ID) {
		log.Printf(
			"consumer: notification cancelled before send, skip id=%s",
			n.ID,
		)
		return nil
	}

	// Ожидания до времени отправки
	if delay := time.Until(n.SendAt); delay > 0 {
		log.Printf(
			"consumer: waiting until send time id=%s delay=%s",
			n.ID,
			delay,
		)

		select {
		case <-time.After(delay):
			// ok
		case <-ctx.Done():
			log.Printf(
				"consumer: context cancelled while waiting id=%s",
				n.ID,
			)
			return ctx.Err()
		}
	}

	// Повторная проверка статуса после sleep
	if c.isCancelled(n.ID) {
		log.Printf(
			"consumer: notification cancelled after wait, skip id=%s",
			n.ID,
		)
		return nil
	}

	// Отправка уведомления
	log.Printf(
		"consumer: sending notification id=%s channel=%s",
		n.ID,
		n.Channel,
	)

	var err error
	switch n.Channel {
	case "email":
		if c.emailSender == nil {
			err = fmt.Errorf("email sender is not configured")
		} else {
			err = c.emailSender.Send(n)
		}

	case "telegram":
		if c.telegramSender == nil {
			err = fmt.Errorf("telegram sender is not configured")
		} else {
			err = c.telegramSender.Send(n)
		}

	default:
		err = fmt.Errorf("unknown channel %s", n.Channel)
	}

	if err != nil {
		n.Meta.Attempt++
		log.Printf(
			"consumer: failed to send notification id=%s channel=%s err=%v",
			n.ID,
			n.Channel,
			err,
		)
		return err
	}

	// Обновляем статус
	n.Status = models.StatusSent
	_ = c.store.Save(n)

	log.Printf(
		"consumer: notification sent successfully id=%s channel=%s",
		n.ID,
		n.Channel,
	)

	return nil
}

// isCancelled - проверка отмены
func (c *NotificationConsumer) isCancelled(id string) bool {
	stored, ok := c.store.Get(id)
	if !ok {
		return false
	}
	return stored.Status == models.StatusCancelled
}
