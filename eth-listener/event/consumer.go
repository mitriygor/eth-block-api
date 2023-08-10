package event

import (
	"encoding/json"
	"eth-listener/internal/eth_block"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	conn            *amqp.Connection
	ethBlockService eth_block.Service
}

func NewConsumer(conn *amqp.Connection, ethBlockService eth_block.Service) (Consumer, error) {
	consumer := Consumer{
		conn:            conn,
		ethBlockService: ethBlockService,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Printf("ERROR::EthListener::Listen::error getting channel: %v\n", err.Error())
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Printf("ERROR::EthListener::Listen::error declaring queue: %v\n", err.Error())
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"eth_blocks",
			false,
			nil,
		)

		if err != nil {
			log.Printf("ERROR::EthListener::Listen::error binding queue: %v\n", err.Error())
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("ERROR::EthListener::Listen::error consuming queue: %v\n", err.Error())
		return err
	}

	for d := range messages {
		var bd eth_block.BlockDetails
		_ = json.Unmarshal(d.Body, &bd)
		consumer.HandlePayload(bd)
	}

	return nil
}

func (consumer *Consumer) HandlePayload(bd eth_block.BlockDetails) {
	log.Printf("EthListener::Listen:HandlePayload::bd: %v\n", bd)
}
