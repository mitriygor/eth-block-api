package event

import (
	"encoding/json"
	"eth-transactions-recorder/internal/eth_transaction"
	"eth-transactions-recorder/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn                  *amqp.Connection
	ethTransactionService eth_transaction.Service
}

func NewConsumer(conn *amqp.Connection, ethTransactionService eth_transaction.Service) (Consumer, error) {
	consumer := Consumer{
		conn:                  conn,
		ethTransactionService: ethTransactionService,
	}

	err := consumer.setup()
	if err != nil {
		logger.Error("eth-transactions-recorder::ERROR::NewConsumer::error setting up consumer", "error", err)
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
		logger.Error("eth-transactions-recorder::ERROR::Listen::error getting channel", "error", err)
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		logger.Error("eth-transactions-recorder::ERROR::Listen::error declaring queue", "error", err)
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
			logger.Error("eth-transactions-recorder::ERROR::Listen::error binding queue", "error", err)
			return err
		}

		if err != nil {
			logger.Error("eth-transactions-recorder::ERROR::Listen::error binding queue", "error", err)
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		logger.Error("eth-transactions-recorder::ERROR::Listen::error consuming queue", "error", err, "queue", q.Name)
		return err
	}

	for d := range messages {

		var bd eth_transaction.EthTransaction
		_ = json.Unmarshal(d.Body, &bd)

		consumer.HandlePayload(bd)
	}

	return nil
}

func (consumer *Consumer) HandlePayload(et eth_transaction.EthTransaction) {
	err := consumer.ethTransactionService.InsertEthTransactionService(et)
	if err != nil {
		logger.Error("eth-transactions-recorder::ERROR::Listen::HandlePayload::error inserting eth block", "error", err)
	}
}
