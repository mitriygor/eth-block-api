package eth_transaction

import (
	"context"
	"eth-transactions-recorder/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
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
	collection := ebr.mongoClient.Database("eth_transactions").Collection("eth_transactions")

	_, err := collection.InsertOne(context.TODO(), et)
	if err != nil {
		logger.Error("eth-transactions-recorder::ERROR::InsertEthTransaction::err", "error", err)
		return err
	}

	return nil
}
