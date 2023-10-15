package main

import (
	"eth-helpers/mongo_helper"
	"eth-helpers/queue_helper/connector"
	"eth-transactions-recorder/event"
	"eth-transactions-recorder/internal/eth_transaction"
	"eth-transactions-recorder/pkg/logger"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
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
		logger.Error("eth-transactions-recorder::ERROR::Error loading .env file", "error", err)
	}

	ethTransactionsMongo := os.Getenv("ETH_TRANSACTIONS_MONGO")
	ethTransactionsMongoUser := os.Getenv("ETH_TRANSACTIONS_MONGO_USER")
	ethTransactionsMongoPassword := os.Getenv("ETH_TRANSACTIONS_MONGO_PASSWORD")

	mongoClient, err := mongo_helper.ConnectToMongo(ethTransactionsMongo, ethTransactionsMongoUser, ethTransactionsMongoPassword)
	if err != nil {
		logger.Error("eth-transactions-recorder::ERROR::MongoDB:connect:error", "error", err)
	}

	ethTransactionsRecorderQueue := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE")
	ethTransactionsRecorderQueueConn, err := connector.ConnectToQueue(ethTransactionsRecorderQueue)
	if err != nil {
		logger.Error("eth-transactions-recorder::ERROR::RabbitMQ:connect:error", "error", err)
		os.Exit(1)
	}

	defer func(ethTransactionsRecorderQueueConn *amqp.Connection) {
		err := ethTransactionsRecorderQueueConn.Close()
		if err != nil {
			logger.Error("eth-transactions-recorder::ERROR::RabbitMQ:connection:close:error", "error", err)
		}
	}(ethTransactionsRecorderQueueConn)

	ethTransactionRepo := eth_transaction.NewEthTransactionRepository(ethTransactionsRecorderQueueConn, mongoClient)
	ethTransactionService := eth_transaction.NewEthTransactionService(ethTransactionRepo)

	consumer, err := event.NewConsumer(ethTransactionsRecorderQueueConn, ethTransactionService)

	if err != nil {
		logger.Error("eth-transactions-recorder::ERRORRabbitMQ:consumer:PANIC", "error", err)
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO"})
	if err != nil {
		logger.Error("eth-transactions-recorder::ERRORRabbitMQ:consume:error", "error", err)
	}
}
