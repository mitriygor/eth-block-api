package main

import (
	"eth-transactions-req/internal/eth_transaction"
	"eth-transactions-req/pkg/queue_helper/connector"
	"eth-transactions-req/pkg/queue_helper/consumer"
	"eth-transactions-req/pkg/url_helper"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR::REQ::Error loading .env file")
	}

	// Connect to RabbitMQ which works with the API in order to request data from the transactionchain and return it to the API
	queueForApi := os.Getenv("QUEUE_FOR_API")
	queueForApiName := os.Getenv("QUEUE_FOR_API_NAME")
	connForApi, err := connector.ConnectToQueue(queueForApi)
	if err != nil {
		log.Printf("ERROR::REQ::connForApi::err: %v\n", err)
		os.Exit(1)
	}
	defer connForApi.Close()

	chForApi, err := connForApi.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "ERROR::REQ::failed to open chForApi:", err)
	}
	defer chForApi.Close()

	// Connect to RabbitMQ which works with the storage in order to store data from the transactionchain
	queueForStorage := os.Getenv("QUEUE_FOR_STORAGE")
	connForStorage, err := connector.ConnectToQueue(queueForStorage)
	if err != nil {
		log.Printf("ERROR::REQ::connForStorage::err: %v\n", err)
		os.Exit(1)
	}
	defer connForStorage.Close()

	chForStorage, err := connForStorage.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "ERROR::REQ::failed to open chForStorage:", err)
	}
	defer chForStorage.Close()

	// Declaring the API's service and repository
	endpoint := os.Getenv("HTTP_ENDPOINT")
	version := os.Getenv("HTTP_ENDPOINT_VERSION")
	secretKey := os.Getenv("HTTP_KEY")
	url := url_helper.GetUrl(endpoint, version, secretKey)

	jsonRpc := os.Getenv("JSONRPC")

	queueForStorageName := os.Getenv("QUEUE_FOR_STORAGE_NAME")

	ethTransactionRepo := eth_transaction.NewEthTransactionRepository(chForStorage, queueForStorageName)
	ethTransactionService := eth_transaction.NewEthTransactionService(ethTransactionRepo, url, jsonRpc)

	consumer.Consume(chForApi, queueForApiName, ethTransactionService)
}
