package eth_transaction

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Repository interface {
	PushEthTransaction(et EthTransaction) error
}

type EthTransactionRepository struct {
	chForStorage        *amqp.Channel
	queueForStorageName string
}

func NewEthTransactionRepository(chForStorage *amqp.Channel, queueForStorageName string) Repository {

	return &EthTransactionRepository{
		chForStorage:        chForStorage,
		queueForStorageName: queueForStorageName,
	}
}

func (ebr *EthTransactionRepository) PushEthTransaction(et EthTransaction) error {

	log.Printf("TRANSACTIONS-REQPushEthTransaction::et: %v\n", et)

	var transaction string

	jsonStr, err := json.MarshalIndent(et, "", "  ")

	log.Printf("TRANSACTIONS-REQPushEthTransaction::jsonStr: %v\n", jsonStr)

	if err == nil {
		transaction = string(jsonStr)
	}

	log.Printf("TRANSACTIONS-REQPushEthTransaction::err: %v\n", err)
	log.Printf("TRANSACTIONS-REQPushEthTransaction::transaction: %v\n", transaction)

	//err = ebr.chForStorage.PublishWithContext(context.TODO(),
	//	ebr.queueForStorageName,
	//	"log.INFO",
	//	true,
	//	false,
	//	amqp.Publishing{
	//		ContentType:  "text/plain",
	//		Body:         []byte(transaction),
	//		DeliveryMode: amqp.Persistent,
	//	},
	//)
	//
	//if err != nil {
	//	log.Printf("ERROR::TRANSACTIONS-REQ::Push::error publishing message: %v\n", err.Error())
	//}

	return err
}
