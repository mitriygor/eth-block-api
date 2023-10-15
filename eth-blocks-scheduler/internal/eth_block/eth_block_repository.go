package eth_block

import (
	"context"
	"encoding/json"
	"eth-blocks-scheduler/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
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
	j, _ := json.MarshalIndent(&bd, "", "\t")

	err := ebr.ethBlocksRecorderQueueCh.PublishWithContext(context.TODO(),
		ebr.ethBlocksRecorderQueueName,
		"log.INFO",
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         j,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		logger.Error("eth-blocks-scheduler::PushBlocksForRecording::error", "error", err)
		return err
	}

	return nil
}

func (ebr *EthBlockRepository) PushBlocksForRedis(bd BlockDetails) error {
	j, _ := json.MarshalIndent(&bd, "", "\t")

	err := ebr.ethRedisRecorderQueueCh.PublishWithContext(context.TODO(),
		ebr.ethRedisRecorderQueueName,
		"log.INFO",
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         j,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		logger.Error("eth-blocks-scheduler::PushBlocksForRedis::error", "error", err)
		return err
	}

	return nil
}

func (ebr *EthBlockRepository) PushTransactionsForScheduling(transactions []string) error {
	j, _ := json.MarshalIndent(&transactions, "", "\t")

	err := ebr.ethTransactionsSchedulerQueueCh.PublishWithContext(context.TODO(),
		ebr.ethTransactionsSchedulerQueueName,
		"log.INFO",
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         j,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		logger.Error("eth-blocks-scheduler::PushTransactionsForScheduling::error", "error", err)
		return err
	}

	return nil
}
