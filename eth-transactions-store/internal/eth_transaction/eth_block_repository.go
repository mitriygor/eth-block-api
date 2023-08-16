package eth_transaction

import (
	"context"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Repository interface {
	InsertEthTransaction(bd EthTransaction) error
}

type EthTransactionRepository struct {
	rabbitConn  *amqp.Connection
	redisClient *redis.Client
	mongoClient *mongo.Client
}

func NewEthTransactionRepository(redisClient *redis.Client, rabbitConn *amqp.Connection, mongoClient *mongo.Client) Repository {

	return &EthTransactionRepository{
		rabbitConn:  rabbitConn,
		redisClient: redisClient,
		mongoClient: mongoClient,
	}
}

func (ebr *EthTransactionRepository) InsertEthTransaction(et EthTransaction) error {
	log.Printf("EthTransactionsStore::InsertEthTransaction::et: %v\n", et)

	collection := ebr.mongoClient.Database("eth_transactions").Collection("eth_transactions")
	log.Printf("EthTransactionsStore::InsertEthTransaction::collection: %v\n", collection)

	_, err := collection.InsertOne(context.TODO(), et)
	if err != nil {
		log.Printf("ERROR::EthTransactionsStore::InsertEthTransaction::err: %v\n", err)
		return err
	}

	return nil
}

func addToQueue(rdb *redis.Client, queueName string, item string, queueLimit int64) {
	rdb.LPush(context.Background(), queueName, item)
	rdb.LTrim(context.Background(), queueName, -queueLimit, -1)
}
