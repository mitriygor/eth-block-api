package event

import (
	"encoding/json"
	"eth-transactions-scheduler/internal/eth_transaction"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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
	log.Printf("eth-blocks-recorder::Listen::topics: %v\n", topics)

	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Printf("eth-blocks-recorder::ERROR::Listen::error getting channel: %v\n", err.Error())
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Printf("eth-blocks-recorder::ERROR::Listen::error declaring queue: %v\n", err.Error())
		return err
	}

	for _, s := range topics {
		log.Printf("eth-blocks-recorder::Listen::binding queue: %v\n", q.Name)
		ch.QueueBind(
			q.Name,
			s,
			"eth_transactions",
			false,
			nil,
		)

		if err != nil {
			log.Printf("eth-blocks-recorder::ERROR::Listen::error binding queue: %v\n", err.Error())
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("eth-transactions-scheduler::Listen::ERROR::Listen::error consuming queue: %v\n", err.Error())
	}

	for d := range messages {
		log.Printf("eth-blocks-recorder::Listen::received message: %v\n", string(d.Body))
		consumer.HandlePayload(d.Body)
	}

	return nil
}

func (consumer *Consumer) HandlePayload(transactions []byte) {

	log.Printf("eth-transactions-scheduler::HandlePayload::transactionsHashes: %v\n", transactions)

	var transactionsHashes []string
	err := json.Unmarshal(transactions, &transactionsHashes)
	if err != nil {
		log.Printf("eth-blocks-recorder::ERROR::Listen::error consuming queue: %v\n", err.Error())
	}

	for _, hash := range transactionsHashes {

		log.Printf("eth-transactions-scheduler::HandlePayload::hash: %v\n", hash)

		et, err := consumer.ethTransactionsService.GetEthTransaction(hash)
		if err != nil {
			log.Printf("eth-transactions-scheduler::ERROR::HandlePayload::GetEthTransaction::error: %v\n", err)
		}
		consumer.ethTransactionsService.PushEthTransactionService(*et)

		time.Sleep(time.Duration(consumer.requesterInterval) * time.Second)
	}
}
