package main

import (
	"eth-helpers/queue_helper/connector"
	"eth-helpers/url_helper"
	"eth-transactions-requester/internal/eth_transaction"
	"eth-transactions-requester/pkg/queue_helper/consumer"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("eth-transactions-requester::ERROR::Error loading .env file")
	}

	// Connect to RabbitMQ which works with the API in order to request data from the transactionchain and return it to the API
	ethTransactionsRequesterQueue := os.Getenv("ETH_TRANSACTIONS_REQUESTER_QUEUE")
	ethTransactionsRequesterQueueName := os.Getenv("ETH_TRANSACTIONS_REQUESTER_QUEUE_NAME")
	ethTransactionsRequesterQueueConn, err := connector.ConnectToQueue(ethTransactionsRequesterQueue)
	if err != nil {
		log.Printf("eth-transactions-requester::ERROR::ethTransactionsRequesterQueueConn::err: %v\n", err)
		os.Exit(1)
	}
	defer ethTransactionsRequesterQueueConn.Close()

	ethTransactionsRequesterQueueCh, err := ethTransactionsRequesterQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-transactions-requester::ERROR::failed to open ethTransactionsRequesterQueueCh:", err)
	}
	defer ethTransactionsRequesterQueueCh.Close()

	// Connect to RabbitMQ which works with the storage in order to store data from the transactionchain
	ethTransactionsRecorderQueue := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE")
	ethTransactionsRecorderQueueName := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE_NAME")
	ethTransactionsRecorderQueueConn, err := connector.ConnectToQueue(ethTransactionsRecorderQueue)
	if err != nil {
		log.Printf("eth-transactions-requester::ERROR::ethTransactionsRecorderQueueConn::err: %v\n", err)
		os.Exit(1)
	}
	defer ethTransactionsRecorderQueueConn.Close()

	ethTransactionsRecorderQueueCh, err := ethTransactionsRecorderQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-transactions-requester::ERROR::failed to open ethTransactionsRecorderQueueCh:", err)
	}
	defer ethTransactionsRecorderQueueCh.Close()

	// Declaring the API's service and repository
	endpoint := os.Getenv("HTTP_ENDPOINT")
	version := os.Getenv("HTTP_ENDPOINT_VERSION")
	secretKey := os.Getenv("HTTP_KEY")
	url := url_helper.GetUrl(endpoint, version, secretKey)

	jsonRpc := os.Getenv("JSONRPC")

	ethTransactionRepo := eth_transaction.NewEthTransactionRepository(ethTransactionsRecorderQueueCh, ethTransactionsRecorderQueueName)
	ethTransactionService := eth_transaction.NewEthTransactionService(ethTransactionRepo, url, jsonRpc)

	consumer.Consume(ethTransactionsRequesterQueueCh, ethTransactionsRequesterQueueName, ethTransactionService)
}
