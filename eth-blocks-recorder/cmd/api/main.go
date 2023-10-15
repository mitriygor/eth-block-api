package main

import (
	"eth-blocks-recorder/event"
	"eth-blocks-recorder/internal/eth_block"
	"eth-blocks-recorder/pkg/logger"
	"eth-helpers/mongo_helper"
	"eth-helpers/queue_helper/connector"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"os"
)

func main() {
	// Initialize the logger
	err := logger.Initialize("info")
	if err != nil {
		panic(err)
	}
	defer func(Log *zap.SugaredLogger) {
		err := Log.Sync()
		if err != nil {
			logger.Error("eth-blocks-recorder:ERROR:sync", "error", err)
		}
	}(logger.Log)

	err = godotenv.Load()
	if err != nil {
		logger.Error("eth-blocks-recorder:ERROR:Error loading .env file", "error", err)
	}

	ethBlocksMongo := os.Getenv("ETH_BLOCKS_MONGO")
	ethBlocksMongoUser := os.Getenv("ETH_BLOCKS_MONGO_USER")
	ethBlocksMongoPassword := os.Getenv("ETH_BLOCKS_MONGO_PASSWORD")

	ethBlocksMongoClient, err := mongo_helper.ConnectToMongo(ethBlocksMongo, ethBlocksMongoUser, ethBlocksMongoPassword)
	if err != nil {
		logger.Error("eth-blocks-recorder:ERROR:MongoDB connect error", "error", err)
	}

	ethBlocksRecorderQueue := os.Getenv("ETH_BLOCKS_RECORDER_QUEUE")
	ethBlocksRecorderQueueConn, err := connector.ConnectToQueue(ethBlocksRecorderQueue)
	if err != nil {
		logger.Error("eth-blocks-recorder:ERROR: RabbitMQ connect error", "error", err)
		os.Exit(1)
	}

	defer func(ethBlocksRecorderQueueConn *amqp.Connection) {
		err := ethBlocksRecorderQueueConn.Close()
		if err != nil {
			logger.Error("eth-blocks-recorder:ERROR: RabbitMQ connection close error", "error", err)
		}
	}(ethBlocksRecorderQueueConn)

	ethBlockRepo := eth_block.NewEthBlockRepository(ethBlocksRecorderQueueConn, ethBlocksMongoClient)
	ethBlockService := eth_block.NewEthBlockService(ethBlockRepo)

	consumer, err := event.NewConsumer(ethBlocksRecorderQueueConn, ethBlockService)

	if err != nil {
		logger.Error("eth-blocks-recorder:ERROR: RabbitMQ consumer panic", "error", err)
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO"})
	if err != nil {
		logger.Error("eth-blocks-recorder:ERROR: RabbitMQ consume error", "error", err)
	}
}
