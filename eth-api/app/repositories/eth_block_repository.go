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
	"go.mongodb.org/mongo-driver/mongo/options"
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

	const key = "eth_blocks_latest"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure the context is cancelled to release resources

	// Retrieve all elements from the list
	values, err := ebr.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		logger.Error("eth-api::ERROR::getLatestEthBlocksFromRedis::LRange::err", err)
		return nil, err
	}

	var blocks []*models.BlockDetails
	for _, val := range values {
		var blockDetails models.BlockDetails
		err = json.Unmarshal([]byte(val), &blockDetails) // Assuming each list item is a JSON serialized BlockDetails
		if err != nil {
			logger.Error("eth-api::ERROR::getLatestEthBlocksFromRedis::Unmarshal::err", err)
			return nil, err
		}

		blocks = append(blocks, &blockDetails)
	}

	return blocks, nil
}

func (ebr *ethBlockRepository) getEthBlockFromRedis(identifier string, collectionName string) (*models.BlockDetails, error) {
	var key strings.Builder
	key.WriteString(collectionName)
	key.WriteString(":")
	key.WriteString(identifier)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Make sure to cancel the context when you are done to release resources

	res, err := ebr.redisClient.HGetAll(ctx, key.String()).Result()

	if err != nil {
		logger.Error("eth-api::ERROR::getEthBlockFromRedis::HGetAll::err", err)
		return nil, err
	}

	blockDetails := &models.BlockDetails{}
	err = mapstructure.Decode(res, blockDetails)

	if err != nil {
		logger.Error("eth-api::ERROR::getEthBlockFromRedis::Decode::err", err)
		return nil, err
	}

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
		logger.Error("eth-api::ERROR::getLatestEthBlockFromMongo::Find::err", err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			logger.Error("eth-api::ERROR::getLatestEthBlockFromMongo::Close", err)
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var block models.BlockDetails
		if err := cursor.Decode(&block); err != nil {
			logger.Error("eth-api::ERROR::getLatestEthBlockFromMongo::Decode::err", err)
			return nil, err
		}
		res = append(res, &block)
	}

	if err := cursor.Err(); err != nil {
		logger.Error("eth-api::ERROR::getLatestEthBlockFromMongo::Cursor::err", err)
		return nil, err
	}

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
		logger.Error("eth-api::ERROR::getEthBlockFromMongo::res::err", err)
		return nil, err
	}

	return &res, nil
}

func (ebr *ethBlockRepository) getEthBlockFromApi(identifier string, identifierType string) (*models.BlockDetails, error) {
	conn, err := connector.ConnectToQueue(ebr.queueHost)
	if err != nil {
		logger.Error("eth-api::ERROR::getEthBlockFromApi::conn", err)
		return nil, err
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			logger.Error("eth-api::ERROR::getEthBlockFromApi::conn", err)
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("eth-api::ERROR::getEthBlockFromApi::ch", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			logger.Error("eth-api::ERROR::getEthBlockFromApi::ch", err)
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
		logger.Error("eth-api::ERROR::getEthBlockFromApi::replyToQueue", err)
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
		logger.Error("eth-api::ERROR::getEthBlockFromApi::msgs", err)
	}

	correlationId := strconv.Itoa(rand.Int())
	blockIdentifier := models.BlockIdentifier{
		Identifier:     identifier,
		IdentifierType: identifierType,
	}

	blockIdentifierJson, err := json.MarshalIndent(blockIdentifier, "", "  ")
	if err != nil {
		logger.Error("eth-api::ERROR::getEthBlockFromApi::blockIdentifierJson", err)
	}

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
		logger.Error("eth-api::ERROR::getEthBlockFromApi::ch", err)
	}

	select {
	case d := <-msgs:
		if d.CorrelationId == correlationId {
			logger.Error("eth-api::getEthBlockFromApi::d.Body", string(d.Body))
			return nil, nil
		}
	case <-time.After(10 * time.Second):
		logger.Error("eth-api::ERROR::getEthBlockFromApi::timeout", err)
	}

	return nil, nil
}

func (ebr *ethBlockRepository) isBlockValid(block *models.BlockDetails) bool {
	return block != nil && block.Number != ""
}
