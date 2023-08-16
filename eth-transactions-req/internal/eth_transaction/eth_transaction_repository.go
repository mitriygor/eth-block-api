package eth_transaction

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Repository interface {
	PushEthTransaction(et EthTransaction) error
}

type EthTransactionRepository struct {
	chForStorage        *amqp.Channel
	queueForStorageName string
}

func NewEthTransactionRepository(chForStorage *amqp.Channel, queueForStorageName string) Repository {

	return &EthTransactionRepository{
		chForStorage:        chForStorage,
		queueForStorageName: queueForStorageName,
	}
}

func (ebr *EthTransactionRepository) PushEthTransaction(et EthTransaction) error {

	log.Printf("EthTransactionsReq::PushEthTransaction::et: %v\n", et)

	var transaction string

	jsonStr, err := json.MarshalIndent(et, "", "  ")

	log.Printf("EthTransactionsReq::PushEthTransaction::jsonStr: %v\n", jsonStr)

	if err == nil {
		transaction = string(jsonStr)
	}

	log.Printf("EthTransactionsReq::PushEthTransaction::err: %v\n", err)
	log.Printf("EthTransactionsReq::PushEthTransaction::transaction: %v\n", transaction)

	err = ebr.chForStorage.PublishWithContext(context.TODO(),
		ebr.queueForStorageName,
		"log.INFO",
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(transaction),
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Printf("EthTransactionsReq::PushEthTransaction::error publishing message: %v\n", err.Error())
	}

	return err
}
