package main

import (
	"eth-blocks-recorder/event"
	"eth-blocks-recorder/internal/eth_block"
	"eth-helpers/mongo_helper"
	"eth-helpers/queue_helper/connector"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("eth-blocks-recorder:ERROR:Error loading .env file")
	}

	ethBlocksMongo := os.Getenv("ETH_BLOCKS_MONGO")
	ethBlocksMongoUser := os.Getenv("ETH_BLOCKS_MONGO_USER")
	ethBlocksMongoPassword := os.Getenv("ETH_BLOCKS_MONGO_PASSWORD")

	ethBlocksMongoClient, err := mongo_helper.ConnectToMongo(ethBlocksMongo, ethBlocksMongoUser, ethBlocksMongoPassword)
	if err != nil {
		log.Panic(err)
	}

	ethBlocksRecorderQueue := os.Getenv("ETH_BLOCKS_RECORDER_QUEUE")
	ethBlocksRecorderQueueConn, err := connector.ConnectToQueue(ethBlocksRecorderQueue)
	if err != nil {
		log.Printf("eth-blocks-recorder:ERROR: RabbitMQ connect error: %v\n", err)
		os.Exit(1)
	}

	log.Println("eth-blocks-recorder::RabbitMQ: Connected to RabbitMQ")
	defer ethBlocksRecorderQueueConn.Close()

	ethBlockRepo := eth_block.NewEthBlockRepository(ethBlocksRecorderQueueConn, ethBlocksMongoClient)
	ethBlockService := eth_block.NewEthBlockService(ethBlockRepo)

	consumer, err := event.NewConsumer(ethBlocksRecorderQueueConn, ethBlockService)

	if err != nil {
		log.Println("eth-blocks-recorder:ERROR: RabbitMQ consumer panic")
		panic(err)
	}

	log.Println("eth-blocks-recorder:RabbitMQ:consumer: Consumer is established")

	err = consumer.Listen([]string{"log.INFO"})
	if err != nil {
		log.Printf("eth-blocks-recorder::ERROR::RabbitMQ:consume:error: %v\n", err)
	}
}
