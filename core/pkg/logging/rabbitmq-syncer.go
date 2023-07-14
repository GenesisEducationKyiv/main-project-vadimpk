package logging

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"go.uber.org/zap/zapcore"
)

type rabbitMQSyncer struct {
	channel *amqp.Channel
}

func NewRabbitMQSyncer(channel *amqp.Channel) (*rabbitMQSyncer, error) {
	_, err := channel.QueueDeclare(
		"logger-info",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(
		"logger-error",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &rabbitMQSyncer{channel: channel}, nil
}

func (r *rabbitMQSyncer) Write(p []byte) (n int, err error) {
	entry := zapcore.Entry{}
	err = json.Unmarshal(p, &entry)
	if err != nil {
		return 0, err
	}

	err = r.channel.Publish(
		"",
		"logger-"+entry.Level.String(),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        p,
		})
	if err != nil {
		log.Println("error publishing message to rabbitmq", err)
		return 0, err
	}

	return len(p), nil
}

func (r *rabbitMQSyncer) Sync() error {
	// RabbitMQ takes care of message persistence, so no action required here.
	return nil
}

func (r *rabbitMQSyncer) Close() error {
	return r.channel.Close()
}
