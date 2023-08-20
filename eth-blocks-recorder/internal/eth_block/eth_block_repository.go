package eth_block

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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
		log.Printf("eth-blocks-recorder::InsertEthBlock::eth_blocks::error: %v\n", err)
		return err
	}

	return nil
}
