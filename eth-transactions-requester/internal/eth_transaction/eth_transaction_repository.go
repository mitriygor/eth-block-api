package eth_transaction

import (
	"context"
	"encoding/json"
	"eth-transactions-requester/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
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
	var transaction string
	jsonStr, err := json.MarshalIndent(et, "", "  ")

	if err == nil {
		transaction = string(jsonStr)
	} else {
		logger.Error("eth-transactions-requester::PushEthTransaction::json.MarshalIndent", "err", err)
	}

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
		logger.Error("eth-transactions-requester::PushEthTransaction::ch.PublishWithContext", "err", err, "transaction", transaction)
	}

	return err
}
