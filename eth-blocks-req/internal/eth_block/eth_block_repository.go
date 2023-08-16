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
	chForStorage        *amqp.Channel
	queueForStorageName string
}

func NewEthBlockRepository(chForStorage *amqp.Channel, queueForStorageName string) Repository {

	return &EthBlockRepository{
		chForStorage:        chForStorage,
		queueForStorageName: queueForStorageName,
	}
}

func (ebr *EthBlockRepository) PushEthBlock(bd BlockDetails) error {

	var block string

	jsonStr, err := json.MarshalIndent(bd, "", "  ")
	if err == nil {
		block = string(jsonStr)
	}

	err = ebr.chForStorage.PublishWithContext(context.TODO(),
		ebr.queueForStorageName,
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
		log.Printf("ERROR::EthEmitter::Push::error publishing message: %v\n", err.Error())
	}

	return err
}
