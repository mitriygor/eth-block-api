package main

import (
	"eth-blocks-requester/internal/eth_block"
	"eth-blocks-requester/pkg/queue_helper/consumer"
	"eth-helpers/queue_helper/connector"
	"eth-helpers/url_helper"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("eth-blocks-requester::ERROR::Error loading .env file")
	}

	// Connect to RabbitMQ which works with the API in order to request data from the blockchain and return it to the API
	ethBlocksRequesterQueue := os.Getenv("ETH_BLOCKS_REQUESTER_QUEUE")
	ethBlocksRequesterQueueName := os.Getenv("ETH_BLOCKS_REQUESTER_QUEUE_NAME")
	ethBlocksRequesterQueueConn, err := connector.ConnectToQueue(ethBlocksRequesterQueue)
	if err != nil {
		log.Printf("eth-blocks-requester::ERROR::ethBlocksRequesterQueueConn::err: %v\n", err)
		os.Exit(1)
	}
	defer ethBlocksRequesterQueueConn.Close()

	ethBlocksRequesterQueueCh, err := ethBlocksRequesterQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-blocks-requester::ERROR::failed to open ethBlocksRequesterQueueCh:", err)
	}
	defer ethBlocksRequesterQueueCh.Close()

	// Connect to RabbitMQ which works with the storage in order to store data from the blockchain
	ethBlocksRecorderQueue := os.Getenv("ETH_BLOCKS_RECORDER_QUEUE")
	ethBlocksRecorderQueueName := os.Getenv("ETH_BLOCKS_RECORDER_QUEUE_NAME")
	ethBlocksRecorderQueueConn, err := connector.ConnectToQueue(ethBlocksRecorderQueue)
	if err != nil {
		log.Printf("eth-blocks-requester::ERROR::ethBlocksRecorderQueueConn::err: %v\n", err)
		os.Exit(1)
	}
	defer ethBlocksRecorderQueueConn.Close()

	ethBlocksRecorderQueueCh, err := ethBlocksRecorderQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-blocks-requester::ERROR::failed to open ethBlocksRecorderQueueCh:", err)
	}
	defer ethBlocksRecorderQueueCh.Close()

	// Declaring the API's service and repository
	endpoint := os.Getenv("HTTP_ENDPOINT")
	version := os.Getenv("HTTP_ENDPOINT_VERSION")
	secretKey := os.Getenv("HTTP_KEY")
	url := url_helper.GetUrl(endpoint, version, secretKey)

	jsonRpc := os.Getenv("JSONRPC")

	ethBlockRepo := eth_block.NewEthBlockRepository(ethBlocksRecorderQueueCh, ethBlocksRecorderQueueName)
	ethBlockService := eth_block.NewEthBlockService(ethBlockRepo, url, jsonRpc)

	consumer.Consume(ethBlocksRequesterQueueCh, ethBlocksRequesterQueueName, ethBlockService)
}
