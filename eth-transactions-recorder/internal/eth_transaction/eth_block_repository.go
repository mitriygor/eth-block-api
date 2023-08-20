package eth_transaction

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Repository interface {
	InsertEthTransaction(bd EthTransaction) error
}

type EthTransactionRepository struct {
	rabbitConn  *amqp.Connection
	mongoClient *mongo.Client
}

func NewEthTransactionRepository(rabbitConn *amqp.Connection, mongoClient *mongo.Client) Repository {

	return &EthTransactionRepository{
		rabbitConn:  rabbitConn,
		mongoClient: mongoClient,
	}
}

func (ebr *EthTransactionRepository) InsertEthTransaction(et EthTransaction) error {
	log.Printf("EthTransactionsStore::InsertEthTransaction::et: %v\n", et)

	collection := ebr.mongoClient.Database("eth_transactions").Collection("eth_transactions")
	log.Printf("EthTransactionsStore::InsertEthTransaction::collection: %v\n", collection)

	_, err := collection.InsertOne(context.TODO(), et)
	if err != nil {
		log.Printf("eth-transactions-recorder::ERROR::InsertEthTransaction::err: %v\n", err)
		return err
	}

	return nil
}
