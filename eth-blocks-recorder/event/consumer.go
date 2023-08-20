package event

import (
	"encoding/json"
	"eth-blocks-recorder/internal/eth_block"
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
	log.Printf("eth-blocks-recorder::Listen::topics: %v\n", topics)

	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Printf("eth-blocks-recorder::ERROR::Listen::error getting channel: %v\n", err.Error())
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Printf("eth-blocks-recorder::ERROR::Listen::error declaring queue: %v\n", err.Error())
		return err
	}

	for _, s := range topics {
		log.Printf("eth-blocks-recorder::Listen::binding queue: %v\n", q.Name)
		ch.QueueBind(
			q.Name,
			s,
			"eth_blocks",
			false,
			nil,
		)

		if err != nil {
			log.Printf("eth-blocks-recorder::ERROR::Listen::error binding queue: %v\n", err.Error())
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("eth-blocks-recorder::ERROR::Listen::error consuming queue: %v\n", err.Error())
		return err
	}

	for d := range messages {
		log.Printf("eth-blocks-recorder::Listen::received message: %v\n", string(d.Body))
		var bd eth_block.BlockDetails
		_ = json.Unmarshal(d.Body, &bd)
		consumer.HandlePayload(bd)
	}

	return nil
}

func (consumer *Consumer) HandlePayload(bd eth_block.BlockDetails) {
	err := consumer.ethBlockService.InsertEthBlockService(bd)
	if err != nil {
		log.Printf("eth-blocks-recorder::ERROR::Listen::HandlePayload::error inserting eth block: %v\n", err.Error())
	}
}
