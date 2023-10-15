package eth_block

import (
	"context"
	"eth-blocks-recorder/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	InsertEthBlock(bd BlockDetails) error
}

type EthBlockRepository struct {
	rabbitConn           *amqp.Connection
	ethBlocksMongoClient *mongo.Client
}

func NewEthBlockRepository(rabbitConn *amqp.Connection, ethBlocksMongoClient *mongo.Client) Repository {

	return &EthBlockRepository{
		rabbitConn:           rabbitConn,
		ethBlocksMongoClient: ethBlocksMongoClient,
	}
}

func (ebr *EthBlockRepository) InsertEthBlock(bd BlockDetails) error {
	collection := ebr.ethBlocksMongoClient.Database("eth_blocks").Collection("eth_blocks")

	_, err := collection.InsertOne(context.TODO(), bd)
	if err != nil {
		logger.Error("eth-blocks-recorder::InsertEthBlock::eth_blocks::error", "error", err)
		return err
	}

	return nil
}
