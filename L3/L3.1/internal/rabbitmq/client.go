package rabbitmq

import (
	"delayed_notifier/internal/config"
	"time"

	wbfrabbit "github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

type Client struct {
	Rmq *wbfrabbit.RabbitClient
}

// NewRabbitClient создаёт RabbitMQ клиент и объявляет очередь + exchange
func NewRabbitClient(cfg config.AppConfig) (*Client, error) {
	rmqCfg := wbfrabbit.ClientConfig{
		URL: cfg.RabbitMQ.URL,
		PublishRetry: retry.Strategy{
			Attempts: 5,
			Delay:    100 * time.Millisecond,
			Backoff:  2,
		},
		ConsumeRetry: retry.Strategy{
			Attempts: 0, // бесконечно
			Delay:    time.Second,
			Backoff:  2,
		},
	}

	rmq, err := wbfrabbit.NewClient(rmqCfg)
	if err != nil {
		return nil, err
	}

	// Создаём отдельный exchange и очередь
	queueName := "notifications_queue"
	exchangeName := "notifications_queue"

	if err := rmq.DeclareQueue(
		queueName,    // queueName
		exchangeName, // exchangeName
		queueName,    // routingKey
		true,         // queueDurable
		false,        // queueAutoDelete
		true,         // exchangeDurable
		nil,          // queueArgs
	); err != nil {
		return nil, err
	}

	return &Client{Rmq: rmq}, nil
}
