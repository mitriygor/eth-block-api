package event

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emitter struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	confirm    chan amqp.Confirmation
}

func (e *Emitter) Close() error {
	return e.channel.Close()
}

func (e *Emitter) Push(event string, severity string) error {
	err := e.channel.PublishWithContext(context.TODO(),
		"eth_blocks",
		severity,
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(event),
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Printf("eth-blocks-scheduler::ERROR::Push::error publishing message: %v\n", err.Error())
	}

	if confirmed := <-e.confirm; !confirmed.Ack {
		log.Printf("eth-blocks-scheduler::UNCOFIRMED MESSAGE!!!\n")
	}

	return err
}
