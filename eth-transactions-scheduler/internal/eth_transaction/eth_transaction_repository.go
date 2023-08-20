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

	log.Printf("eth-transactions-scheduler::PushEthTransaction::et: %v\n", et)

	var transaction string

	jsonStr, err := json.MarshalIndent(et, "", "  ")

	log.Printf("eth-transactions-scheduler::PushEthTransaction::jsonStr: %v\n", jsonStr)

	if err == nil {
		transaction = string(jsonStr)
	}

	log.Printf("eth-transactions-scheduler::PushEthTransaction::err: %v\n", err)
	log.Printf("eth-transactions-scheduler::PushEthTransaction::transaction: %v\n", transaction)
	log.Printf("eth-transactions-scheduler::PushEthTransaction::jsonStr: %v\n", jsonStr)
	log.Printf("eth-transactions-scheduler::PushEthTransaction::et: %v\n", et)

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
		log.Printf("eth-transactions-scheduler::PushEthTransaction::error publishing message: %v\n", err.Error())
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
		log.Printf("eth-transactions-scheduler::PushEthTransaction::error publishing message: %v\n", err.Error())
	}

	return err
}
