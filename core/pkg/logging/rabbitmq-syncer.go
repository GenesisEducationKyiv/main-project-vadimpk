package logging

import (
	"github.com/streadway/amqp"
)

type rabbitMQSyncer struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQSyncer(channel *amqp.Channel, queue amqp.Queue) *rabbitMQSyncer {
	return &rabbitMQSyncer{channel: channel, queue: queue}
}

func (r *rabbitMQSyncer) Write(p []byte) (n int, err error) {
	err = r.channel.Publish(
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        p,
		})
	if err != nil {
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
