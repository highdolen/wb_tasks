package broker

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	wbfkafka "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

type Broker struct {
	producer *wbfkafka.Producer
	consumer *wbfkafka.Consumer
}

// New - инициализация Broker
func New(brokers []string, topic, groupID string) *Broker {
	return &Broker{
		producer: wbfkafka.NewProducer(brokers, topic),
		consumer: wbfkafka.NewConsumer(brokers, topic, groupID),
	}
}

// PublishTask - отправка id изображения в Kafka
func (b *Broker) PublishTask(ctx context.Context, id string) error {
	return b.producer.Send(ctx, nil, []byte(id))
}

// StartWorker - запуск фонового чтение из Kafka
func (b *Broker) StartWorker(ctx context.Context, handler func(id string)) {
	msgChan := make(chan kafka.Message)

	// стратегия ретраев при проблемах с Kafka
	strategy := retry.Strategy{
		Attempts: 5,
		Delay:    2 * time.Second,
		Backoff:  2,
	}

	// запуск чтения сообщений
	b.consumer.StartConsuming(ctx, msgChan, strategy)

	go func() {
		for msg := range msgChan {
			id := string(msg.Value)

			// бизнес-логика обработки
			handler(id)

			// подтверждаем обработку сообщения
			if err := b.consumer.Commit(ctx, msg); err != nil {
				log.Println("failed to commit message:", err)
			}
		}
	}()
}

// Close - закрытие соединения с Kafka
func (b *Broker) Close() {
	_ = b.producer.Close()
	_ = b.consumer.Close()
}
