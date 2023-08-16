package consumer

import (
	"context"
	"encoding/json"
	"eth-transactions-req/internal/eth_transaction"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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
		//var eti eth_transaction.EthTransactionIdentifier

		//log.Printf("TRANSACTIONS-REQ::Consume::d.Body: %v\n", d.Body)
		//log.Printf("TRANSACTIONS-REQ::Consume::d.Body: %v\n", []byte(d.Body))
		log.Printf("TRANSACTIONS-REQ::Consume::d.Body: %v\n", string(d.Body))
		hash := string(d.Body)
		//err := json.Unmarshal([]byte(d.Body), &eti)
		//if err != nil {
		//	log.Fatalf("ERROR::TRANSACTIONS-REQ::Error unmarshaling JSON: %v\n", err)
		//}

		transactionResponse, err := service.GetEthTransaction(hash)

		log.Printf("TRANSACTIONS-REQ::Consume::err: %v\n", err)
		log.Printf("TRANSACTIONS-REQ::Consume::transactionResponse: %v\n", transactionResponse)

		if err == nil && transactionResponse != nil {
			jsonStr, err := json.MarshalIndent(transactionResponse, "", "  ")
			log.Printf("TRANSACTIONS-REQ::Consume::err: %v\n", err)
			log.Printf("TRANSACTIONS-REQ::Consume::jsonStr: %v\n", jsonStr)

			if err == nil {
				response = string(jsonStr)
			}
		}

		log.Printf("TRANSACTIONS-REQ::Consume::response: %v\n", response)

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

		log.Printf("TRANSACTIONS-REQ::Consume::err: %v\n", err)

		if err != nil {
			log.Printf("Failed to publish a message: %s", err)
		}
	}
}
