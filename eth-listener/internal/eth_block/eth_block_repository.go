package eth_block

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Repository interface {
	InsertEthBlock(bd BlockDetails) error
}

type EthBlockRepository struct {
	DBW         *sql.DB
	rabbitConn  *amqp.Connection
	redisClient *redis.Client
	mongoClient *mongo.Client
}

func NewEthBlockRepository(dbw *sql.DB, redisClient *redis.Client, rabbitConn *amqp.Connection, mongoClient *mongo.Client) Repository {

	return &EthBlockRepository{
		DBW:         dbw,
		rabbitConn:  rabbitConn,
		redisClient: redisClient,
		mongoClient: mongoClient,
	}
}

func (ebr *EthBlockRepository) InsertEthBlock(bd BlockDetails) error {
	fmt.Printf("CreateEthBlock::bd: %v\n", bd)

	collection := ebr.mongoClient.Database("eth_blocks").Collection("eth_blocks")

	_, err := collection.InsertOne(context.TODO(), bd)
	if err != nil {
		log.Println("Error inserting into eth_blocks:", err)
		return err
	}

	return nil
}

func addToQueue(rdb *redis.Client, queueName string, item string, queueLimit int64) {
	rdb.LPush(context.Background(), queueName, item)

	rdb.LTrim(context.Background(), queueName, -queueLimit, -1)
}
