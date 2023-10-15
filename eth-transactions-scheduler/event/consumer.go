package event

import (
	"encoding/json"
	"eth-transactions-scheduler/internal/eth_transaction"
	"eth-transactions-scheduler/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Consumer struct {
	conn                   *amqp.Connection
	ethTransactionsService eth_transaction.Service
	requesterInterval      int
}

func NewConsumer(conn *amqp.Connection, ethTransactionsService eth_transaction.Service, requesterInterval int) (Consumer, error) {
	consumer := Consumer{
		conn:                   conn,
		ethTransactionsService: ethTransactionsService,
		requesterInterval:      requesterInterval,
	}

	err := consumer.setup()
	if err != nil {
		logger.Error("eth-transactions-scheduler:ERROR:NewConsumer", "error", err)
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
		logger.Error("eth-blocks-recorder:ERROR:Listen", "error", err)
		return err
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			logger.Error("eth-blocks-recorder:ERROR:Listen", "error", err)
		}
	}(ch)

	q, err := declareRandomQueue(ch)
	if err != nil {
		logger.Error("eth-blocks-recorder:ERROR:Listen", "error", err)
		return err
	}

	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"eth_transactions",
			false,
			nil,
		)
		if err != nil {
			logger.Error("eth-blocks-recorder:ERROR:Listen", "error", err)
			return err
		}

		if err != nil {
			logger.Error("eth-blocks-recorder:ERROR:Listen", "error", err)
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		logger.Error("eth-blocks-recorder:ERROR:Listen", "error", err)
	}

	for d := range messages {
		logger.Error("eth-blocks-recorder:Listen", "message", string(d.Body))
		consumer.HandlePayload(d.Body)
	}

	return nil
}

func (consumer *Consumer) HandlePayload(transactions []byte) {
	var transactionsHashes []string
	err := json.Unmarshal(transactions, &transactionsHashes)
	if err != nil {
		logger.Error("eth-blocks-recorder:ERROR:Listen", "error", err)
	}

	for _, hash := range transactionsHashes {
		et, err := consumer.ethTransactionsService.GetEthTransaction(hash)
		if err != nil {
			logger.Error("eth-transactions-scheduler::ERROR::HandlePayload::GetEthTransaction::error", "error", err)
		}
		consumer.ethTransactionsService.PushEthTransactionService(*et)

		time.Sleep(time.Duration(consumer.requesterInterval) * time.Second)
	}
}
