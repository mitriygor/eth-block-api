package eth_block

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Repository interface {
	PushBlocksForRecording(bd BlockDetails) error
	PushBlocksForRedis(bd BlockDetails) error
	PushTransactionsForScheduling(transactions []string) error
}

type EthBlockRepository struct {
	ethBlocksRecorderQueueCh          *amqp.Channel
	ethBlocksRecorderQueueName        string
	ethRedisRecorderQueueCh           *amqp.Channel
	ethRedisRecorderQueueName         string
	ethTransactionsSchedulerQueueCh   *amqp.Channel
	ethTransactionsSchedulerQueueName string
}

func NewEthBlockRepository(ethBlocksRecorderQueueCh *amqp.Channel, ethBlocksRecorderQueueName string, ethRedisRecorderQueueCh *amqp.Channel, ethRedisRecorderQueueName string, ethTransactionsSchedulerQueueCh *amqp.Channel, ethTransactionsSchedulerQueueName string) Repository {

	return &EthBlockRepository{
		ethBlocksRecorderQueueCh:          ethBlocksRecorderQueueCh,
		ethBlocksRecorderQueueName:        ethBlocksRecorderQueueName,
		ethRedisRecorderQueueCh:           ethRedisRecorderQueueCh,
		ethRedisRecorderQueueName:         ethRedisRecorderQueueName,
		ethTransactionsSchedulerQueueCh:   ethTransactionsSchedulerQueueCh,
		ethTransactionsSchedulerQueueName: ethTransactionsSchedulerQueueName,
	}
}

func (ebr *EthBlockRepository) PushBlocksForRecording(bd BlockDetails) error {
	log.Printf("eth-blocks-scheduler::PushBlocksForRecording::transactions::%v\n", bd)

	j, _ := json.MarshalIndent(&bd, "", "\t")

	log.Printf("eth-blocks-scheduler::PushBlocksForRecording::j::%v\n", j)

	err := ebr.ethBlocksRecorderQueueCh.PublishWithContext(context.TODO(),
		ebr.ethBlocksRecorderQueueName,
		"log.INFO",
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(j),
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (ebr *EthBlockRepository) PushBlocksForRedis(bd BlockDetails) error {
	log.Printf("eth-blocks-scheduler::PushBlocksForRedis::transactions::%v\n", bd)

	j, _ := json.MarshalIndent(&bd, "", "\t")

	log.Printf("eth-blocks-scheduler::PushBlocksForRecording::j::%v\n", j)

	err := ebr.ethRedisRecorderQueueCh.PublishWithContext(context.TODO(),
		ebr.ethRedisRecorderQueueName,
		"log.INFO",
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(j),
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (ebr *EthBlockRepository) PushTransactionsForScheduling(transactions []string) error {

	log.Printf("eth-blocks-scheduler::PushTransactionsForScheduling::transactions::%v\n", transactions)
	j, _ := json.MarshalIndent(&transactions, "", "\t")

	log.Printf("eth-blocks-scheduler::PushTransactionsForScheduling::j::%v\n", j)
	log.Printf("eth-blocks-scheduler::PushTransactionsForScheduling::ebr.ethTransactionsSchedulerQueueName::%v\n", ebr.ethTransactionsSchedulerQueueName)

	err := ebr.ethTransactionsSchedulerQueueCh.PublishWithContext(context.TODO(),
		ebr.ethTransactionsSchedulerQueueName,
		"log.INFO",
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(j),
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
