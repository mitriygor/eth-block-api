package main

import (
	"eth-helpers/queue_helper/connector"
	"eth-helpers/url_helper"
	"eth-transactions-requester/internal/eth_transaction"
	"eth-transactions-requester/pkg/logger"
	"eth-transactions-requester/pkg/queue_helper/consumer"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	err := logger.Initialize("info")
	if err != nil {
		panic(err)
	}
	defer func(Log *zap.SugaredLogger) {
		err := Log.Sync()
		if err != nil {
			logger.Error("eth-transactions-recorder:ERROR:sync", "error", err)
		}
	}(logger.Log)

	err = godotenv.Load()
	if err != nil {
		logger.Error("eth-transactions-requester::ERROR::Error loading .env file", "error", err)
	}

	// Connect to RabbitMQ which works with the API in order to request data from the transactionchain and return it to the API
	ethTransactionsRequesterQueue := os.Getenv("ETH_TRANSACTIONS_REQUESTER_QUEUE")
	ethTransactionsRequesterQueueName := os.Getenv("ETH_TRANSACTIONS_REQUESTER_QUEUE_NAME")
	ethTransactionsRequesterQueueConn, err := connector.ConnectToQueue(ethTransactionsRequesterQueue)
	if err != nil {
		logger.Error("eth-transactions-requester::ERROR::ethTransactionsRequesterQueueConn::err", "error", err)
		os.Exit(1)
	}
	defer func(ethTransactionsRequesterQueueConn *amqp.Connection) {
		err := ethTransactionsRequesterQueueConn.Close()
		if err != nil {
			logger.Error("eth-transactions-requester::ERROR::RabbitMQ:connection:close:error", "error", err)
		}
	}(ethTransactionsRequesterQueueConn)

	ethTransactionsRequesterQueueCh, err := ethTransactionsRequesterQueueConn.Channel()
	if err != nil {
		logger.Error("eth-transactions-requester::ERROR::failed to open ethTransactionsRequesterQueueCh", "error", err)
	}
	defer func(ethTransactionsRequesterQueueCh *amqp.Channel) {
		err := ethTransactionsRequesterQueueCh.Close()
		if err != nil {
			logger.Error("eth-transactions-requester::ERROR::RabbitMQ:channel:close:error", "error", err)
		}
	}(ethTransactionsRequesterQueueCh)

	// Connect to RabbitMQ which works with the storage in order to store data from the transactionchain
	ethTransactionsRecorderQueue := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE")
	ethTransactionsRecorderQueueName := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE_NAME")
	ethTransactionsRecorderQueueConn, err := connector.ConnectToQueue(ethTransactionsRecorderQueue)
	if err != nil {
		logger.Error("eth-transactions-requester::ERROR::ethTransactionsRecorderQueueConn::err", "error", err)
		os.Exit(1)
	}
	defer func(ethTransactionsRecorderQueueConn *amqp.Connection) {
		err := ethTransactionsRecorderQueueConn.Close()
		if err != nil {
			logger.Error("eth-transactions-requester::ERROR::RabbitMQ:connection:close:error", "error", err)
		}
	}(ethTransactionsRecorderQueueConn)

	ethTransactionsRecorderQueueCh, err := ethTransactionsRecorderQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-transactions-requester::ERROR::failed to open ethTransactionsRecorderQueueCh:", err)
	}
	defer func(ethTransactionsRecorderQueueCh *amqp.Channel) {
		err := ethTransactionsRecorderQueueCh.Close()
		if err != nil {
			logger.Error("eth-transactions-requester::ERROR::RabbitMQ:channel:close:error", "error", err)
		}
	}(ethTransactionsRecorderQueueCh)

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
