package event

import (
	"encoding/json"
	"eth-redis-recorder/internal/eth_block"
	"eth-redis-recorder/internal/eth_transaction"
	"eth-redis-recorder/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
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
		logger.Error("eth-redis-recorder::ERROR::NewConsumer::error setting up consumer", "error", err)
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
		logger.Error("eth-redis-recorder::ERROR::Listen::error getting channel", "error", err)
		return err
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			logger.Error("eth-redis-recorder::ERROR::Listen::error closing channel", "error", err)
		}
	}(ch)

	q, err := declareRandomQueue(ch)
	if err != nil {
		logger.Error("eth-redis-recorder::ERROR::Listen::error declaring queue", "error", err)
		return err
	}

	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"eth_redis",
			false,
			nil,
		)
		if err != nil {
			logger.Error("eth-redis-recorder::ERROR::Listen::error binding queue", "error", err)
			return err
		}

		if err != nil {
			logger.Error("eth-redis-recorder::ERROR::Listen::error binding queue", "error", err)
			return err
		}
	}

	_, err = ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		logger.Error("eth-redis-recorder::ERROR::Listen::error consuming queue", "error", err)
		return err
	}

	return nil
}

func (consumer *Consumer) process(data string) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		logger.Error("eth-redis-recorder::process::ERROR::Unmarshal", "error", err)
		return
	}

	if _, ok := result["transactions"]; ok {
		var bd eth_block.BlockDetails
		if err := json.Unmarshal([]byte(data), &bd); err != nil {
			logger.Error("eth-redis-recorder::process::ERROR::Unmarshal::BlockDetails", "error", err)
		} else {
			consumer.handleBlock(bd)
		}
	} else if _, ok := result["transactionIndex"]; ok {
		var et eth_transaction.EthTransaction
		if err := json.Unmarshal([]byte(data), &et); err != nil {
			logger.Error("eth-redis-recorder::process::ERROR::Unmarshal::EthTransaction", "error", err)
		} else {
			consumer.handleTransaction(et)
		}
	} else {
		logger.Error("eth-redis-recorder::process::ERROR::Unknown type")
	}
}

func (consumer *Consumer) handleBlock(bd eth_block.BlockDetails) {
	consumer.ethBlockService.AddBlockService(bd)
}

func (consumer *Consumer) handleTransaction(et eth_transaction.EthTransaction) {
	consumer.ethTransactionService.AddTransactionService(et)
}
