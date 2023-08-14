package consumer

import (
	"context"
	"encoding/json"
	"eth-req/internal/eth_block"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
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
		log.Fatalf("%s: %s", "Failed to register a consumer", err)
	}

	for d := range msgs {
		var response string

		var bi eth_block.BlockIdentifier

		err := json.Unmarshal([]byte(d.Body), &bi)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
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
			log.Printf("Failed to publish a message: %s", err)
		}
	}
}
