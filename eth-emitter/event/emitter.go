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

func NewEmitter(conn *amqp.Connection) (*Emitter, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err = channel.Confirm(false); err != nil {
		return nil, err
	}

	confirm := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	emitter := &Emitter{
		connection: conn,
		channel:    channel,
		confirm:    confirm,
	}

	return emitter, nil
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
		log.Printf("ERROR::EthEmitter::Push::error publishing message: %v\n", err.Error())
	}

	if confirmed := <-e.confirm; !confirmed.Ack {
		log.Printf("\nUNCOFIRMED MESSAGE!!!\n")
	}

	return err
}
