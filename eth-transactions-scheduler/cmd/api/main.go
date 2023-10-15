package main

import (
	"eth-helpers/queue_helper/connector"
	"eth-helpers/url_helper"
	"eth-transactions-scheduler/event"
	"eth-transactions-scheduler/internal/eth_transaction"
	"eth-transactions-scheduler/pkg/logger"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"os"
	"strconv"
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
		logger.Error("eth-transactions-scheduler:ERROR:loading_env", "error", err)
	}

	requesterInterval, err := strconv.Atoi(os.Getenv("REQUESTER_INTERVAL"))
	if err != nil {
		logger.Error("eth-transactions-scheduler:ERROR:requesterInterval", "error", err)
	}

	// Queue from which the service will consume the transactions
	ethTransactionsSchedulerQueue := os.Getenv("ETH_TRANSACTIONS_SCHEDULER_QUEUE")
	ethTransactionsSchedulerQueueConn, err := connector.ConnectToQueue(ethTransactionsSchedulerQueue)
	if err != nil {
		logger.Error("eth-transactions-scheduler:ERROR:connForApi", "error", err)
		os.Exit(1)
	}
	defer func(ethTransactionsSchedulerQueueConn *amqp.Connection) {
		err := ethTransactionsSchedulerQueueConn.Close()
		if err != nil {
			logger.Error("eth-transactions-scheduler:ERROR:connForApi:close", "error", err)
		}
	}(ethTransactionsSchedulerQueueConn)

	// Queue to which the service will send the transactions for further recording to the database
	ethTransactionsRecorderQueue := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE")
	ethTransactionsRecorderQueueName := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE_NAME")
	ethTransactionsRecorderQueueConn, err := connector.ConnectToQueue(ethTransactionsRecorderQueue)
	if err != nil {
		logger.Error("eth-transactions-scheduler:ERROR:connForApi", "error", err)
		os.Exit(1)
	}
	defer func(ethTransactionsRecorderQueueConn *amqp.Connection) {
		err := ethTransactionsRecorderQueueConn.Close()
		if err != nil {
			logger.Error("eth-transactions-scheduler:ERROR:connForApi:close", "error", err)
		}
	}(ethTransactionsRecorderQueueConn)

	ethTransactionsRecorderQueueCh, err := ethTransactionsRecorderQueueConn.Channel()
	if err != nil {
		logger.Error("eth-transactions-scheduler:ERROR:connForApi:close", "error", err)
		panic(err)
	}
	defer func(ethTransactionsRecorderQueueCh *amqp.Channel) {
		err := ethTransactionsRecorderQueueCh.Close()
		if err != nil {
			logger.Error("eth-transactions-scheduler:ERROR:connForApi:close", "error", err)
		}
	}(ethTransactionsRecorderQueueCh)

	// Queue to which the service will send the transactions for further recording to the cache
	ethRedisRecorderQueue := os.Getenv("ETH_REDIS_RECORDER_QUEUE")
	ethRedisRecorderQueueName := os.Getenv("ETH_REDIS_RECORDER_QUEUE_NAME")
	ethRedisRecorderQueueConn, err := connector.ConnectToQueue(ethRedisRecorderQueue)
	if err != nil {
		logger.Error("eth-transactions-scheduler:ERROR:connForApi", "error", err)
		os.Exit(1)
	}
	defer func(ethRedisRecorderQueueConn *amqp.Connection) {
		err := ethRedisRecorderQueueConn.Close()
		if err != nil {
			logger.Error("eth-transactions-scheduler:ERROR:connForApi:close", "error", err)
		}
	}(ethRedisRecorderQueueConn)

	ethRedisRecorderQueueCh, err := ethRedisRecorderQueueConn.Channel()
	if err != nil {
		logger.Error("eth-transactions-scheduler:ERROR:connForApi:close", "error", err)
		panic(err)
	}
	defer func(ethRedisRecorderQueueCh *amqp.Channel) {
		err := ethRedisRecorderQueueCh.Close()
		if err != nil {
			logger.Error("eth-transactions-scheduler:ERROR:connForApi:close", "error", err)
		}
	}(ethRedisRecorderQueueCh)

	// Declaring the API's service and repository
	endpoint := os.Getenv("HTTP_ENDPOINT")
	version := os.Getenv("HTTP_ENDPOINT_VERSION")
	secretKey := os.Getenv("HTTP_KEY")
	url := url_helper.GetUrl(endpoint, version, secretKey)

	jsonRpc := os.Getenv("JSONRPC")

	ethTransactionRepo := eth_transaction.NewEthTransactionRepository(ethTransactionsRecorderQueueCh, ethTransactionsRecorderQueueName, ethRedisRecorderQueueCh, ethRedisRecorderQueueName)
	ethTransactionService := eth_transaction.NewEthTransactionService(ethTransactionRepo, url, jsonRpc)

	consumer, err := event.NewConsumer(ethTransactionsSchedulerQueueConn, ethTransactionService, requesterInterval)

	if err != nil {
		logger.Error("eth-transactions-scheduler:ERROR: RabbitMQ consumer panic", "error", err)
	}

	err = consumer.Listen([]string{"log.INFO"})
	if err != nil {
		logger.Error("eth-transactions-scheduler:ERROR: RabbitMQ consumer panic", "error", err)
	}
}
