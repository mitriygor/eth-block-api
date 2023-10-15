package event

import (
	"encoding/json"
	"eth-blocks-recorder/internal/eth_block"
	"eth-blocks-recorder/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
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
		logger.Error("eth-blocks-recorder::Listen::error getting channel", "error", err)
		return err
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			logger.Error("eth-blocks-recorder::Listen::error closing channel", "error", err)
		}
	}(ch)

	q, err := declareRandomQueue(ch)
	if err != nil {
		logger.Error("eth-blocks-recorder::Listen::error declaring queue", "error", err)
		return err
	}

	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"eth_blocks",
			false,
			nil,
		)
		if err != nil {
			logger.Error("eth-blocks-recorder::Listen::error binding queue", "error", err)
			return err
		}

		if err != nil {
			logger.Error("eth-blocks-recorder::Listen::error binding queue", "error", err)
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		logger.Error("eth-blocks-recorder::Listen::error consuming queue", "error", err)
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
	err := consumer.ethBlockService.InsertEthBlockService(bd)
	if err != nil {
		logger.Error("eth-blocks-recorder::Listen::HandlePayload::error inserting eth block", "error", err)
	}
}
