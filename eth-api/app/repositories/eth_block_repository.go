package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"eth-api/app/models"
	"eth-helpers/queue_helper/connector"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type EthBlockRepository interface {
	GetEthBlockByIdentifier(identifier string, identifierType string) (*models.BlockDetails, error)
	GetLatestEthBlocks() ([]*models.BlockDetails, error)
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

func (ebr *ethBlockRepository) GetLatestEthBlocks() ([]*models.BlockDetails, error) {

	log.Printf("eth-api::EthBlockRepository::GetLatestEthBlocks")

	mongoDb := "eth_blocks"
	mongoCollection := "eth_blocks"

	latestBlockDetails, err := ebr.getLatestEthBlocksFromRedis()

	if err != nil || len(latestBlockDetails) == 0 {
		latestBlockDetails, err = ebr.getLatestEthBlockFromMongo(mongoDb, mongoCollection)
		if err != nil {
			return nil, err
		}
	}

	return latestBlockDetails, nil
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

		log.Printf("eth-api::GetEthBlockByIdentifier::blockDetails: %v\n", blockDetails)
		log.Printf("eth-api::GetEthBlockByIdentifier::err: %v\n", err)

		if err != nil || !ebr.isBlockValid(blockDetails) {
			blockDetails, err = ebr.getEthBlockFromApi(identifier, identifierType)

			if err != nil || !ebr.isBlockValid(blockDetails) {
				return nil, errors.New("block not found")
			}
		}
	}

	return blockDetails, nil
}

func (ebr *ethBlockRepository) getLatestEthBlocksFromRedis() ([]*models.BlockDetails, error) {
	fmt.Printf("eth-api::EthBlockRepository::getLatestEthBlocksFromRedis")

	const key = "eth_blocks_latest"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure the context is cancelled to release resources

	// Retrieve all elements from the list
	values, err := ebr.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		log.Printf("eth-api::ERROR::getLatestEthBlocksFromRedis::LRange::err: %v\n", err)
		return nil, err
	}

	var blocks []*models.BlockDetails
	for _, val := range values {
		var blockDetails models.BlockDetails
		err = json.Unmarshal([]byte(val), &blockDetails) // Assuming each list item is a JSON serialized BlockDetails
		if err != nil {
			log.Printf("eth-api::ERROR::getLatestEthBlocksFromRedis::Unmarshal::err: %v\n", err)
			return nil, err
		}

		blocks = append(blocks, &blockDetails)
	}

	log.Printf("eth-api::getLatestEthBlocksFromRedis::blocks: %v\n", blocks)
	return blocks, nil
}

func (ebr *ethBlockRepository) getEthBlockFromRedis(identifier string, collectionName string) (*models.BlockDetails, error) {
	var key strings.Builder
	key.WriteString(collectionName)
	key.WriteString(":")
	key.WriteString(identifier)

	log.Printf("eth-api::getEthBlockFromRedis::key: %v\n", key)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Make sure to cancel the context when you are done to release resources

	res, err := ebr.redisClient.HGetAll(ctx, key.String()).Result()

	if err != nil {
		log.Printf("eth-api::ERROR::getEthBlockFromRedis::res::err: %v\n", err)
		return nil, err
	}

	log.Printf("eth-api::getEthBlockFromRedis::res: %v\n", res)

	blockDetails := &models.BlockDetails{}
	err = mapstructure.Decode(res, blockDetails)

	if err != nil {
		log.Printf("eth-api::ERROR::getEthBlockFromRedis::blockDetails::err: %v\n", err)
		return nil, err
	}

	log.Printf("eth-api::getEthBlockFromRedis::blockDetails: %v\n", blockDetails)

	return blockDetails, nil
}

func (ebr *ethBlockRepository) getLatestEthBlockFromMongo(dbName string, collectionName string) ([]*models.BlockDetails, error) {
	var res []*models.BlockDetails
	collection := ebr.mongoClient.Database(dbName).Collection(collectionName)

	// Adjusting the filter to find the latest 50 blocks.
	filter := bson.D{}
	sort := bson.D{{"number", -1}} // assuming "number" is the block number
	limit := int64(50)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(sort).SetLimit(limit))
	if err != nil {
		log.Printf("eth-api::ERROR::getLatestEthBlockFromMongo::res::err: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var block models.BlockDetails
		if err := cursor.Decode(&block); err != nil {
			log.Printf("eth-api::ERROR::getLatestEthBlockFromMongo::Decode::err: %v\n", err)
			return nil, err
		}
		res = append(res, &block)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("eth-api::ERROR::getLatestEthBlockFromMongo::Cursor::err: %v\n", err)
		return nil, err
	}

	log.Printf("eth-api::getLatestEthBlockFromMongo::res: %v\n", res)

	return res, nil
}

func (ebr *ethBlockRepository) getEthBlockFromMongo(identifier string, identifierType string, dbName string, collectionName string) (*models.BlockDetails, error) {

	var res models.BlockDetails
	collection := ebr.mongoClient.Database(dbName).Collection(collectionName)
	filter := bson.M{identifierType: identifier}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&res)

	if err != nil {
		log.Printf("eth-api::ERROR::getEthBlockFromMongo::res::err: %v\n", err)
		return nil, err
	}

	log.Printf("eth-api::getEthBlockFromMongo::res: %v\n", res)

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
		log.Printf("eth-api::ERROR::getEthBlockFromApi::ch::err: %v\n", err)
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
		log.Printf("eth-api::ERROR::getEthBlockFromApi::replyToQueue::err: %v\n", err)
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
		log.Printf("eth-api::ERROR::getEthBlockFromApi::msgs::err: %v\n", err)
	}

	correlationId := strconv.Itoa(rand.Int())

	blockIdentifier := models.BlockIdentifier{
		Identifier:     identifier,
		IdentifierType: identifierType,
	}
	//blockIdentifierStr := fmt.Sprintf("%+v", blockIdentifier)

	blockIdentifierJson, err := json.MarshalIndent(blockIdentifier, "", "  ")

	if err != nil {
		log.Printf("eth-api::ERROR::getEthBlockFromApi::blockIdentifierJson::err: %v\n", err)
	}

	log.Printf("eth-api::getEthBlockFromApi::blockIdentifierJson: %v\n", blockIdentifierJson)

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
			Body:          blockIdentifierJson,
		})
	if err != nil {
		log.Printf("eth-api::ERROR::getEthBlockFromApi::ch::err: %v\n", err)
	}

	select {
	case d := <-msgs:
		if d.CorrelationId == correlationId {
			log.Printf("eth-api::getEthBlockFromApi::d.Body: %v\n", string(d.Body))
			return nil, nil
		}
	case <-time.After(10 * time.Second):
		log.Printf("eth-api::ERROR:::getEthBlockFromApi:timeout\n")
	}

	return nil, nil
}

func (ebr *ethBlockRepository) isBlockValid(block *models.BlockDetails) bool {
	return block != nil && block.Number != ""
}
