package event

import (
	"encoding/json"
	"eth-redis-recorder/internal/eth_block"
	"eth-redis-recorder/internal/eth_transaction"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	conn                  *amqp.Connection
	ethBlockService       eth_block.Service
	ethTransactionService eth_transaction.Service
}

func NewConsumer(conn *amqp.Connection, ethBlockService eth_block.Service, ethTransactionService eth_transaction.Service) (Consumer, error) {
	consumer := Consumer{
		conn:                  conn,
		ethBlockService:       ethBlockService,
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
	log.Printf("eth-redis-recorder::Listen::topics: %v\n", topics)

	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Printf("eth-redis-recorder::ERROR::Listen::error getting channel: %v\n", err.Error())
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Printf("eth-redis-recorder::ERROR::Listen::error declaring queue: %v\n", err.Error())
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"eth_redis",
			false,
			nil,
		)

		if err != nil {
			log.Printf("eth-redis-recorder::ERROR::Listen::error binding queue: %v\n", err.Error())
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("eth-redis-recorder::ERROR::Listen::error consuming queue: %v\n", err.Error())
		return err
	}

	for d := range messages {
		log.Printf("eth-redis-recorder::Listen::message: %v\n", d.Body)

	}

	return nil
}

func (consumer *Consumer) process(data string) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Printf("eth-redis-recorder::process::ERROR::Unmarshal: %v", err)
		return
	}

	if _, ok := result["transactions"]; ok {
		var bd eth_block.BlockDetails
		if err := json.Unmarshal([]byte(data), &bd); err != nil {
			log.Printf("eth-redis-recorder::process::ERROR::Unmarshal::BlockDetails: %v", err)
		} else {
			consumer.handleBlock(bd)
		}
	} else if _, ok := result["transactionIndex"]; ok {
		var et eth_transaction.EthTransaction
		if err := json.Unmarshal([]byte(data), &et); err != nil {
			log.Printf("eth-redis-recorder::process::ERROR::Unmarshal::EthTransaction: %v", err)
		} else {
			consumer.handleTransaction(et)
		}
	} else {
		fmt.Println("eth-redis-recorder::process::ERROR::Unknown type")
	}
}

func (consumer *Consumer) handleBlock(bd eth_block.BlockDetails) {
	log.Printf("eth-redis-recorder::HandleBlock::block: %v\n", bd)
	consumer.ethBlockService.AddBlockService(bd)
}

func (consumer *Consumer) handleTransaction(et eth_transaction.EthTransaction) {
	log.Printf("eth-redis-recorder::HandleTransaction::block: %v\n", et)
	consumer.ethTransactionService.AddTransactionService(et)
}
