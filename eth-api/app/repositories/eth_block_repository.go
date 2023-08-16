package repositories

import (
	"context"
	"eth-api/app/models"
	"eth-api/pkg/queue_helper/connector"
	"fmt"
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

type EthBlockRepository interface {
	GetEthBlockByIdentifier(identifier string, identifierType string) (*models.BlockDetails, error)
}

type ethBlockRepository struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
	queueHost   string
	queueName   string
}

func NewEthBlockRepository(mongoClient *mongo.Client, redisClient *redis.Client, queueHost string, queueName string) EthBlockRepository {
	return &ethBlockRepository{
		mongoClient: mongoClient,
		redisClient: redisClient,
		queueHost:   queueHost,
		queueName:   queueName,
	}
}

func (ebr *ethBlockRepository) GetEthBlockByIdentifier(identifier string, identifierType string) (*models.BlockDetails, error) {

	mongoDb := "eth_blocks"
	mongoCollection := "eth_blocks"

	mongoField := "number"
	redisCollection := "eth_blocks_by_number"

	if identifierType == "hash" {
		mongoField = "hash"
		redisCollection = "eth_blocks_by_hash"
	}

	blockDetails, err := ebr.getEthBlockFromRedis(identifier, redisCollection)

	if err != nil || !ebr.isBlockValid(blockDetails) {
		blockDetails, err = ebr.getEthBlockFromMongo(identifier, mongoField, mongoDb, mongoCollection)

		if err != nil || !ebr.isBlockValid(blockDetails) {
			blockDetails, err = ebr.getEthBlockFromApi(identifier, identifierType)

			if err != nil || !ebr.isBlockValid(blockDetails) {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func (ebr *ethBlockRepository) getEthBlockFromRedis(identifier string, collectionName string) (*models.BlockDetails, error) {
	var key strings.Builder
	key.WriteString(collectionName)
	key.WriteString(":")
	key.WriteString(identifier)

	log.Printf("API::getEthBlockFromRedis::key: %v;", key)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Make sure to cancel the context when you are done to release resources

	res, err := ebr.redisClient.HGetAll(ctx, key.String()).Result()

	if err != nil {
		log.Printf("ERROR::API::getEthBlockFromRedis::res::err: %v;", err)
		return nil, err
	}

	log.Printf("API::getEthBlockFromRedis::res: %v", res)

	blockDetails := &models.BlockDetails{}
	err = mapstructure.Decode(res, blockDetails)

	if err != nil {
		log.Printf("ERROR::API::getEthBlockFromRedis::blockDetails::err: %v;", err)
		return nil, err
	}

	log.Printf("API::getEthBlockFromRedis::blockDetails: %v", blockDetails)

	return blockDetails, nil
}

func (ebr *ethBlockRepository) getEthBlockFromMongo(identifier string, identifierType string, dbName string, collectionName string) (*models.BlockDetails, error) {
	var res models.BlockDetails
	collection := ebr.mongoClient.Database(dbName).Collection(collectionName)
	filter := bson.M{identifierType: identifier}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&res)

	if err != nil {
		log.Printf("ERROR::API::getEthBlockFromMongo::res::err: %v;", err)
		return nil, err
	}

	log.Printf("API::getEthBlockFromMongo::res: %v;", res)

	return &res, nil
}

func (ebr *ethBlockRepository) getEthBlockFromApi(identifier string, identifierType string) (*models.BlockDetails, error) {
	conn, err := connector.ConnectToQueue(ebr.queueHost)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR::API::getEthBlockFromApi::ch::err: %v\n", err)
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
		log.Printf("ERROR::API::getEthBlockFromApi::replyToQueue::err: %v\n", err)
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
		log.Printf("ERROR::API::getEthBlockFromApi::msgs::err: %v\n", err)
	}

	correlationId := strconv.Itoa(rand.Int())

	blockIdentifier := models.BlockIdentifier{
		Identifier:     identifier,
		IdentifierType: identifierType,
	}
	blockIdentifierStr := fmt.Sprintf("%+v", blockIdentifier)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		ebr.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationId,
			ReplyTo:       replyToQueue.Name,
			Body:          []byte(blockIdentifierStr),
		})
	if err != nil {
		log.Printf("ERROR::API::getEthBlockFromApi::ch::err: %v\n", err)
	}

	select {
	case d := <-msgs:
		if d.CorrelationId == correlationId {
			log.Printf("API::getEthBlockFromApi::d.Body: %v\n", string(d.Body))
			return nil, nil
		}
	case <-time.After(5 * time.Second):
		log.Printf("ERROR:API:getEthBlockFromApi:timeout\n")
	}

	return nil, nil
}

func (ebr *ethBlockRepository) isBlockValid(block *models.BlockDetails) bool {
	return block != nil && block.Number != ""
}
