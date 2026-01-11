package rabbitmq

import (
	"context"
	"encoding/json"

	"delayed_notifier/internal/models"

	wbfrabbit "github.com/wb-go/wbf/rabbitmq"
)

type Producer struct {
	pub *wbfrabbit.Publisher
}

const (
	exchangeName = "notifications_queue"
	contentType  = "application/json"
	routingKey   = "notifications_queue"
)

func NewProducer(c *Client) *Producer {
	p := wbfrabbit.NewPublisher(c.Rmq, exchangeName, contentType)
	return &Producer{pub: p}
}

// PublishNotification отправляет уведомление в RabbitMQ
func (p *Producer) PublishNotification(ctx context.Context, n models.Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return p.pub.Publish(ctx, body, routingKey)
}
