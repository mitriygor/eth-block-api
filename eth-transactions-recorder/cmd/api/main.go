package main

import (
	"eth-helpers/mongo_helper"
	"eth-helpers/queue_helper/connector"
	"eth-transactions-recorder/event"
	"eth-transactions-recorder/internal/eth_transaction"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("eth-transactions-recorder::ERROR::Error loading .env file")
	}

	ethTransactionsMongo := os.Getenv("ETH_TRANSACTIONS_MONGO")
	ethTransactionsMongoUser := os.Getenv("ETH_TRANSACTIONS_MONGO_USER")
	ethTransactionsMongoPassword := os.Getenv("ETH_TRANSACTIONS_MONGO_PASSWORD")

	mongoClient, err := mongo_helper.ConnectToMongo(ethTransactionsMongo, ethTransactionsMongoUser, ethTransactionsMongoPassword)
	if err != nil {
		log.Panic(err)
	}

	ethTransactionsRecorderQueue := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE")
	ethTransactionsRecorderQueueConn, err := connector.ConnectToQueue(ethTransactionsRecorderQueue)
	if err != nil {
		log.Printf("eth-transactions-recorder::ERRORRabbitMQ:connect:error: %v\n", err)
		os.Exit(1)
	}

	log.Println("EthListener::RabbitMQ: Connected to RabbitMQ")
	defer ethTransactionsRecorderQueueConn.Close()

	ethTransactionRepo := eth_transaction.NewEthTransactionRepository(ethTransactionsRecorderQueueConn, mongoClient)
	ethTransactionService := eth_transaction.NewEthTransactionService(ethTransactionRepo)

	consumer, err := event.NewConsumer(ethTransactionsRecorderQueueConn, ethTransactionService)

	if err != nil {
		log.Println("eth-transactions-recorder::ERRORRabbitMQ:consumer:PANIC")
		panic(err)
	}

	log.Println("EthListener::RabbitMQ:consumer: Consumer is established")

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Printf("eth-transactions-recorder::ERRORRabbitMQ:consume:error: %v\n", err)
	}
}
