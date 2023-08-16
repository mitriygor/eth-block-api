package event

import (
	"encoding/json"
	"eth-transactions-store/internal/eth_transaction"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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
	log.Printf("EthTransactionsStore::Listen::topics: %v\n", topics)

	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Printf("ERROR::EthTransactionsStore::Listen::error getting channel: %v\n", err.Error())
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Printf("ERROR::EthTransactionsStore::Listen::error declaring queue: %v\n", err.Error())
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"eth_transactions",
			false,
			nil,
		)

		if err != nil {
			log.Printf("ERROR::EthTransactionsStore::Listen::error binding queue: %v\n", err.Error())
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("ERROR::EthTransactionsStore::Listen::error consuming queue: %v\n", err.Error())
		return err
	}

	for d := range messages {
		log.Printf("EthTransactionsStore::Listen::d: %v\n", d)
		log.Printf("EthTransactionsStore::Listen::d.Body: %v\n", d.Body)

		var bd eth_transaction.EthTransaction
		_ = json.Unmarshal(d.Body, &bd)

		consumer.HandlePayload(bd)
	}

	return nil
}

func (consumer *Consumer) HandlePayload(et eth_transaction.EthTransaction) {
	err := consumer.ethTransactionService.InsertEthTransactionService(et)
	if err != nil {
		log.Printf("ERROR::EthTransactionsStore::Listen::HandlePayload::error inserting eth block: %v\n", err.Error())
	}
}
