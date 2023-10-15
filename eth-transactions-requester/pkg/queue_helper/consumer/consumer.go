package consumer

import (
	"context"
	"encoding/json"
	"eth-transactions-requester/internal/eth_transaction"
	"eth-transactions-requester/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume(ch *amqp.Channel, name string, service eth_transaction.Service) {
	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("eth-transactions-requester::Consume::ch.QueueDeclare", "err", err, "name", name)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("eth-transactions-requester::Consume::ch.Consume", "err", err, "name", name)
	}

	for d := range msgs {
		var response string
		hash := string(d.Body)

		transactionResponse, err := service.GetEthTransaction(hash)

		if err == nil && transactionResponse != nil {
			jsonStr, err := json.MarshalIndent(transactionResponse, "", "  ")
			if err == nil {
				response = string(jsonStr)
			} else {
				logger.Error("eth-transactions-requester::Consume::json.MarshalIndent", "err", err)
			}
		} else if err != nil {
			logger.Error("eth-transactions-requester::Consume::service.GetEthTransaction", "err", err)
		}

		err = ch.PublishWithContext(context.TODO(),
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(response),
			})

		if err != nil {
			logger.Error("eth-transactions-requester::Consume::ch.PublishWithContext", "err", err)
		}
	}
}
