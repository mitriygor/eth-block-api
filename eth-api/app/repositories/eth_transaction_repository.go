package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"eth-api/app/helpers/logger"
	"eth-api/app/models"
	"eth-helpers/queue_helper/connector"
	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type EthTransactionRepository interface {
	GetTransactionByHash(hash string) (*models.EthTransaction, error)
	GetTransactionsByAddress(address string) ([]*models.EthTransaction, error)
}

type ethTransactionRepository struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
	queueHost   string
	queueName   string
}

func NewEthTransactionRepository(mongoClient *mongo.Client, redisClient *redis.Client, queueHost string, queueName string) EthTransactionRepository {
	return &ethTransactionRepository{
		mongoClient: mongoClient,
		redisClient: redisClient,
		queueHost:   queueHost,
		queueName:   queueName,
	}
}

func (etr *ethTransactionRepository) GetTransactionByHash(hash string) (*models.EthTransaction, error) {
	mongoDb := "eth_transactions"
	mongoCollection := "eth_transactions"

	transactionDetails, err := etr.getEthTransactionFromRedis(hash)

	if err != nil || !etr.isTransactionValid(transactionDetails) {
		transactionDetails, err = etr.getEthTransactionFromMongo(hash, mongoDb, mongoCollection)

		if err != nil || !etr.isTransactionValid(transactionDetails) {
			transactionDetails, err = etr.getEthTransactionFromApi(hash)

			if err != nil || !etr.isTransactionValid(transactionDetails) {
				return nil, errors.New("transaction not found")
			}
		}
	}

	return transactionDetails, nil
}

func (etr *ethTransactionRepository) GetTransactionsByAddress(address string) ([]*models.EthTransaction, error) {

	mongoDb := "eth_transactions"
	mongoCollection := "eth_transactions"

	transactions, err := etr.getEthTransactionsFromRedis(address)

	if err != nil || !etr.isTransactionsListValid(transactions) {
		transactions, err = etr.getEthTransactionsFromMongo(address, mongoDb, mongoCollection)

		if err != nil || !etr.isTransactionsListValid(transactions) {
			return nil, errors.New("transactions not found")
		}
	}

	return transactions, nil
}

func (etr *ethTransactionRepository) getEthTransactionFromRedis(hash string) (*models.EthTransaction, error) {

	var key strings.Builder
	key.WriteString("eth_transactions_by_hash")
	key.WriteString(":")
	key.WriteString(hash)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Make sure to cancel the context when you are done to release resources

	res, err := etr.redisClient.HGetAll(ctx, key.String()).Result()

	if err != nil {
		return nil, err
	}

	ethTransaction := &models.EthTransaction{}
	err = mapstructure.Decode(res, ethTransaction)

	if err != nil {
		return nil, err
	}

	return ethTransaction, nil
}

func (etr *ethTransactionRepository) getEthTransactionFromMongo(hash string, dbName string, collectionName string) (*models.EthTransaction, error) {

	var res models.EthTransaction
	collection := etr.mongoClient.Database(dbName).Collection(collectionName)
	filter := bson.M{"hash": hash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (etr *ethTransactionRepository) getEthTransactionFromApi(hash string) (*models.EthTransaction, error) {
	conn, err := connector.ConnectToQueue(etr.queueHost)
	if err != nil {
		logger.Error("eth-api::ERROR::getEthTransactionFromApi::conn", err)
		return nil, err
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			logger.Error("eth-api::ERROR::getEthTransactionFromApi::conn::defer", err)
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("eth-api::ERROR::getEthTransactionFromApi::ch", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			logger.Error("eth-api::ERROR::getEthTransactionFromApi::ch::defer", err)
		}
	}(ch)

	replyToQueue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		logger.Error("eth-api::ERROR::getEthTransactionFromApi::replyToQueue", err)
	}

	msgs, err := ch.Consume(
		replyToQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("eth-api::ERROR::getEthTransactionFromApi::msgs", err)
	}

	correlationId := strconv.Itoa(rand.Int())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		etr.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationId,
			ReplyTo:       replyToQueue.Name,
			Body:          []byte(hash),
		})
	if err != nil {
		logger.Error("eth-api::ERROR::getEthTransactionFromApi::ch", err)
	}

	select {
	case d := <-msgs:
		if d.CorrelationId == correlationId {
			var et models.EthTransaction
			err = json.Unmarshal(d.Body, &et)

			if err != nil {
				logger.Error("eth-api::ERROR::getEthTransactionFromApi::json.Unmarshal", err)
				return nil, err
			}

			return &et, nil
		}
	case <-time.After(5 * time.Second):
		logger.Error("eth-api::ERROR::getEthTransactionFromApi::timeout", err)
		return nil, errors.New("timeout")
	}

	return nil, errors.New("no response")
}

func (etr *ethTransactionRepository) getEthTransactionsFromRedis(address string) ([]*models.EthTransaction, error) {

	return nil, nil
}

func (etr *ethTransactionRepository) getEthTransactionsFromMongo(address string, dbName string, collectionName string) ([]*models.EthTransaction, error) {
	collection := etr.mongoClient.Database(dbName).Collection(collectionName)
	filter := bson.D{{"accessList.address", address}}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			logger.Error("eth-api::ERROR::getEthTransactionsFromMongo::defer", err)
		}
	}(cur, context.TODO())

	var results []*models.EthTransaction

	// Iterate through the results
	for cur.Next(context.TODO()) {
		var result models.EthTransaction
		err := cur.Decode(&result)
		if err != nil {
			logger.Error("eth-api::ERROR::getEthTransactionsFromMongo::cur", err)
		}
		results = append(results, &result)
	}

	if err := cur.Err(); err != nil {
		logger.Error("eth-api::ERROR::getEthTransactionsFromMongo::cur", err)
	}

	return results, nil
}

func (etr *ethTransactionRepository) isTransactionValid(transaction *models.EthTransaction) bool {
	return transaction != nil && transaction.Hash != ""
}

func (etr *ethTransactionRepository) isTransactionsListValid(transactions []*models.EthTransaction) bool {
	return transactions != nil && len(transactions) > 0
}
