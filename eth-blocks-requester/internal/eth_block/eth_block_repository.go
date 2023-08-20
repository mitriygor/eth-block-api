package eth_block

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Repository interface {
	PushEthBlock(bd BlockDetails) error
}

type EthBlockRepository struct {
	ethBlocksRecorderQueueCh   *amqp.Channel
	ethBlocksRecorderQueueName string
}

func NewEthBlockRepository(ethBlocksRecorderQueueCh *amqp.Channel, ethBlocksRecorderQueueName string) Repository {

	return &EthBlockRepository{
		ethBlocksRecorderQueueCh:   ethBlocksRecorderQueueCh,
		ethBlocksRecorderQueueName: ethBlocksRecorderQueueName,
	}
}

func (ebr *EthBlockRepository) PushEthBlock(bd BlockDetails) error {

	var block string

	jsonStr, err := json.MarshalIndent(bd, "", "  ")
	if err == nil {
		block = string(jsonStr)
	}

	err = ebr.ethBlocksRecorderQueueCh.PublishWithContext(context.TODO(),
		ebr.ethBlocksRecorderQueueName,
		"log.INFO",
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(block),
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Printf("eth-blocks-requester::ERROR::Push::error publishing message: %v\n", err.Error())
	}

	return err
}
