package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"eth-api/app/models"
	"eth-api/pkg/queue_helper/connector"
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

	log.Printf("API::GetTransactionByHash::hash: %v\n", hash)
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
	transactions, err := etr.getEthTransactionsFromRedis(address)

	if err != nil || !etr.isTransactionsListValid(transactions) {
		transactions, err = etr.getEthTransactionsFromMongo(address)

		if err != nil || !etr.isTransactionsListValid(transactions) {
			return nil, nil
		}
	}

	return nil, nil
}

func (etr *ethTransactionRepository) getEthTransactionFromRedis(hash string) (*models.EthTransaction, error) {

	var key strings.Builder
	key.WriteString("eth_transactions_by_hash")
	key.WriteString(":")
	key.WriteString(hash)

	log.Printf("API::getEthTransactionFromRedis::key: %v\n", key)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Make sure to cancel the context when you are done to release resources

	res, err := etr.redisClient.HGetAll(ctx, key.String()).Result()

	if err != nil {
		log.Printf("ERROR::API::getEthTransactionFromRedis::res::err: %v\n", err)
		return nil, err
	}

	log.Printf("API::getEthTransactionFromRedis::res: %v\n", res)

	ethTransaction := &models.EthTransaction{}
	err = mapstructure.Decode(res, ethTransaction)

	if err != nil {
		log.Printf("ERROR::API::getEthTransactionFromRedis::ethTransaction::err: %v\n", err)
		return nil, err
	}

	log.Printf("API::getEthTransactionFromRedis::ethTransaction: %v\n", ethTransaction)

	return ethTransaction, nil
}

func (etr *ethTransactionRepository) getEthTransactionFromMongo(hash string, dbName string, collectionName string) (*models.EthTransaction, error) {

	var res models.EthTransaction
	collection := etr.mongoClient.Database(dbName).Collection(collectionName)
	filter := bson.M{"hash": hash}

	log.Printf("API::getEthTransactionFromMongo::hash: %v\n", hash)
	log.Printf("API::getEthTransactionFromMongo::dbName: %v\n", dbName)
	log.Printf("API::getEthTransactionFromMongo::collectionName: %v\n", collectionName)
	log.Printf("API::getEthTransactionFromMongo::collection: %v\n", collection)
	log.Printf("API::getEthTransactionFromMongo::filter: %v\n", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&res)

	if err != nil {
		log.Printf("ERROR::API::getEthTransactionFromMongo::res::err: %v\n", err)
		return nil, err
	}

	log.Printf("API::getEthTransactionFromMongo::res: %v\n", res)

	return &res, nil
}

func (etr *ethTransactionRepository) getEthTransactionFromApi(hash string) (*models.EthTransaction, error) {
	log.Printf("API::getEthTransactionFromApi::hash: %v\n", hash)

	conn, err := connector.ConnectToQueue(etr.queueHost)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR::API::getEthTransactionFromApi::ch::err: %v\n", err)
	}
	defer ch.Close()

	replyToQueue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Printf("ERROR::API::getEthTransactionFromApi::replyToQueue::err: %v\n", err)
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
		log.Printf("ERROR::API::getEthTransactionFromApi::msgs::err: %v\n", err)
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
		log.Printf("ERROR::API::getEthTransactionFromApi::ch::err: %v\n", err)
	}

	select {
	case d := <-msgs:
		if d.CorrelationId == correlationId {
			log.Printf("API::getEthTransactionFromApi::d.Body: %v\n", string(d.Body))
			var et models.EthTransaction

			err = json.Unmarshal(d.Body, &et)

			log.Printf("API::getEthTransactionFromApi::err: %v\n", err)
			log.Printf("API::getEthTransactionFromApi::et: %v\n", et)

			if err != nil {
				log.Printf("ERROR::API::getEthTransactionFromApi::et::err: %v\n", err)
				return nil, err
			}

			return &et, nil
		}
	case <-time.After(5 * time.Second):
		log.Printf("ERROR:API:getEthTransactionFromApi:timeout\n")
		return nil, errors.New("timeout")
	}

	return nil, errors.New("no response")
}

func (etr *ethTransactionRepository) getEthTransactionsFromRedis(address string) ([]*models.EthTransaction, error) {
	return nil, nil
}

func (etr *ethTransactionRepository) getEthTransactionsFromMongo(address string) ([]*models.EthTransaction, error) {
	return nil, nil
}

func (ebr *ethTransactionRepository) isTransactionValid(transaction *models.EthTransaction) bool {
	return transaction != nil && transaction.Hash != ""
}

func (ebr *ethTransactionRepository) isTransactionsListValid(transactions []*models.EthTransaction) bool {
	return transactions != nil && len(transactions) > 0
}
