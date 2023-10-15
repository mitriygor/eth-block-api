package consumer

import (
	"context"
	"encoding/json"
	"eth-blocks-requester/internal/eth_block"
	"eth-blocks-requester/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume(ch *amqp.Channel, name string, service eth_block.Service) {
	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("eth-blocks-requester::Failed to declare a queue", "error", err)
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
		logger.Error("eth-blocks-requester::Failed to register a consumer", "error", err)
	}

	for d := range msgs {
		var response string

		var bi eth_block.BlockIdentifier

		err := json.Unmarshal(d.Body, &bi)
		if err != nil {
			logger.Error("eth-blocks-requester::Error unmarshaling JSON", "error", err)
		}

		blockResponse, err := service.GetEthBlock(bi.Identifier, bi.IdentifierType)

		if err == nil && blockResponse != nil {
			jsonStr, err := json.MarshalIndent(blockResponse, "", "  ")
			if err == nil {
				response = string(jsonStr)
			}
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
			logger.Error("eth-blocks-requester::Failed to publish a message", "error", err)
		}
	}
}
