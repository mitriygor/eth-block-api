package consumer

import (
	"context"
	"encoding/json"
	"eth-blocks-requester/internal/eth_block"
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

		err := json.Unmarshal(d.Body, &bi)
		if err != nil {
			log.Printf("eth-blocks-requester::d.Body: %v\n", d.Body)
			log.Printf("eth-blocks-requester::[]byte(d.Body): %v\n", d.Body)
			log.Printf("eth-blocks-requester::string([]byte(d.Body)): %v\n", string(d.Body))
			log.Fatalf("eth-blocks-requester::Error unmarshaling JSON: %v\n", err)
		}

		log.Printf("eth-blocks-requester::bi: %v\n", bi)

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
			log.Printf("eth-blocks-requester::Failed to publish a message: %s", err)
		}
	}
}
