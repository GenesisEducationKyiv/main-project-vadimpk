package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var level = flag.String("level", "info", "log level")
var broker = flag.String("broker", "rabbitmq", "broker name")

type Consumer interface {
	Consume(out chan<- []byte)
	Close()
}

func main() {
	flag.Parse()

	if level == nil || broker == nil || *level == "" || *broker == "" {
		log.Fatal("invalid args")
	}

	consumer := consumerFactory("rabbitmq", "logger-"+*level)

	out := make(chan []byte)

	go consumer.Consume(out)

	for {
		select {
		case msg := <-out:
			io.WriteString(os.Stdout, string(msg))
		}
	}
}

func consumerFactory(consumer, queue string) Consumer {
	switch consumer {
	case "rabbitmq":
		return newRabbitMQConsumer(queue)
	}
	return nil
}

type rabbitMQConsumer struct {
	queue string

	conn *amqp.Connection
	ch   *amqp.Channel
}

func newRabbitMQConsumer(queue string) *rabbitMQConsumer {
	rabbitmqConn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatal("failed to init rabbitmq connection", "err", err)
	}

	rabbitmqChannel, err := rabbitmqConn.Channel()
	if err != nil {
		log.Fatal("failed to init rabbitmq channel", "err", err)
	}

	return &rabbitMQConsumer{
		queue: queue,
		conn:  rabbitmqConn,
		ch:    rabbitmqChannel,
	}
}

func (r *rabbitMQConsumer) Consume(out chan<- []byte) {
	msg, err := r.ch.Consume(
		r.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("failed to init rabbitmq channel", "err", err)
		close(out)
		return
	}

	for {
		select {
		case m := <-msg:
			out <- m.Body
			if err != nil {
				log.Println("failed to write", "err", err)
				continue
			}
			m.Ack(false)
		case <-time.After(100 * time.Millisecond): // if no message received for 100ms, close the channel
			close(out)
			return
		}
	}
}

func (r *rabbitMQConsumer) Close() {
	close := func(c io.Closer) {
		err := c.Close()
		if err != nil {
			log.Println("failed to close rabbitmq", "err", err)
		}
	}

	close(r.ch)
	close(r.conn)
}
