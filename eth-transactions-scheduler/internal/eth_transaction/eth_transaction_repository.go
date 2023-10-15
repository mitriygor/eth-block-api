package eth_transaction

import (
	"context"
	"encoding/json"
	"eth-transactions-scheduler/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Repository interface {
	PushEthTransaction(et EthTransaction) error
}

type EthTransactionRepository struct {
	ethTransactionsRecorderQueueCh   *amqp.Channel
	ethTransactionsRecorderQueueName string
	ethRedisRecorderQueueCh          *amqp.Channel
	ethRedisRecorderQueueName        string
}

func NewEthTransactionRepository(ethTransactionsRecorderQueueCh *amqp.Channel, ethTransactionsRecorderQueueName string, ethRedisRecorderQueueCh *amqp.Channel, ethRedisRecorderQueueName string) Repository {

	return &EthTransactionRepository{
		ethTransactionsRecorderQueueCh:   ethTransactionsRecorderQueueCh,
		ethTransactionsRecorderQueueName: ethTransactionsRecorderQueueName,
		ethRedisRecorderQueueCh:          ethRedisRecorderQueueCh,
		ethRedisRecorderQueueName:        ethRedisRecorderQueueName,
	}
}

func (ebr *EthTransactionRepository) PushEthTransaction(et EthTransaction) error {

	var transaction string

	jsonStr, err := json.MarshalIndent(et, "", "  ")
	if err == nil {
		transaction = string(jsonStr)
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
		logger.Error("eth-transactions-scheduler:ERROR:PushEthTransaction", "error", err)
	}

	err = ebr.ethRedisRecorderQueueCh.PublishWithContext(context.TODO(),
		ebr.ethRedisRecorderQueueName,
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
		logger.Error("eth-transactions-scheduler:ERROR:PushEthTransaction", "error", err)
	}

	return err
}
