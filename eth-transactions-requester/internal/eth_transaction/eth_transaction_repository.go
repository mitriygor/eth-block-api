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
	ethTransactionsRecorderQueueCh   *amqp.Channel
	ethTransactionsRecorderQueueName string
}

func NewEthTransactionRepository(ethTransactionsRecorderQueueCh *amqp.Channel, ethTransactionsRecorderQueueName string) Repository {

	return &EthTransactionRepository{
		ethTransactionsRecorderQueueCh:   ethTransactionsRecorderQueueCh,
		ethTransactionsRecorderQueueName: ethTransactionsRecorderQueueName,
	}
}

func (ebr *EthTransactionRepository) PushEthTransaction(et EthTransaction) error {

	log.Printf("eth-transactions-requester::PushEthTransaction::et: %v\n", et)

	var transaction string

	jsonStr, err := json.MarshalIndent(et, "", "  ")

	log.Printf("eth-transactions-requester::PushEthTransaction::jsonStr: %v\n", jsonStr)

	if err == nil {
		transaction = string(jsonStr)
	}

	log.Printf("eth-transactions-requester::PushEthTransaction::err: %v\n", err)
	log.Printf("eth-transactions-requester::PushEthTransaction::transaction: %v\n", transaction)

	err = ebr.ethTransactionsRecorderQueueCh.PublishWithContext(context.TODO(),
		ebr.ethTransactionsRecorderQueueName,
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
		log.Printf("eth-transactions-requester::PushEthTransaction::error publishing message: %v\n", err.Error())
	}

	return err
}
