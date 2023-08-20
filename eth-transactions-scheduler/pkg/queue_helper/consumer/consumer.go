package consumer

import (
	"eth-transactions-scheduler/internal/eth_transaction"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func Consume(ch *amqp.Channel, name string, service eth_transaction.Service) {

	log.Println("eth-transactions-scheduler::Consume::Consume")
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

	log.Println("eth-transactions-scheduler::Consume::Consume::for!!!!!!!!!!!")
	log.Printf("eth-transactions-scheduler::Consume::Consume::for::msgs: %v\n", msgs)
	log.Printf("eth-transactions-scheduler::Consume::Consume::for::q: %v\n", q)
	log.Printf("eth-transactions-scheduler::Consume::Consume::for::q.Name: %v\n", q.Name)

	for d := range msgs {

		log.Printf("eth-transactions-scheduler::Consume::d: %v\n", d)

		//var response string
		//var eti eth_transaction.EthTransactionIdentifier

		//log.Printf("eth-transactions-scheduler::Consume::d.Body: %v\n", d.Body)
		//log.Printf("eth-transactions-scheduler::Consume::d.Body: %v\n", []byte(d.Body))
		//log.Printf("eth-transactions-scheduler::Consume::d.Body: %v\n", string(d.Body))
		//body := string(d.Body)

		//log.Printf("eth-transactions-scheduler::Consume::hash: %v\n", body)

		//err := json.Unmarshal([]byte(d.Body), &eti)
		//if err != nil {
		//	log.Fatalf("eth-transactions-scheduler::ERROR::Error unmarshaling JSON: %v\n", err)
		//}
		//
		//transactionResponse, err := service.GetEthTransaction(hash)
		//
		//log.Printf("eth-transactions-scheduler::Consume::err: %v\n", err)
		//log.Printf("eth-transactions-scheduler::Consume::transactionResponse: %v\n", transactionResponse)
		//
		//if err == nil && transactionResponse != nil {
		//	jsonStr, err := json.MarshalIndent(transactionResponse, "", "  ")
		//	log.Printf("eth-transactions-scheduler::Consume::err: %v\n", err)
		//	log.Printf("eth-transactions-scheduler::Consume::jsonStr: %v\n", jsonStr)
		//
		//	if err == nil {
		//		response = string(jsonStr)
		//	}
		//}
		//
		//log.Printf("eth-transactions-scheduler::Consume::response: %v\n", response)
		//
		//err = ch.PublishWithContext(context.TODO(),
		//	"",
		//	d.ReplyTo,
		//	false,
		//	false,
		//	amqp.Publishing{
		//		ContentType:   "text/plain",
		//		CorrelationId: d.CorrelationId,
		//		Body:          []byte(response),
		//	})
		//
		//log.Printf("eth-transactions-scheduler::Consume::err: %v\n", err)

		if err != nil {
			log.Printf("Failed to publish a message: %s", err)
		}
	}
}
