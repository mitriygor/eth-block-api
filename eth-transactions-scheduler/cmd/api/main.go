package main

import (
	"eth-transactions-scheduler/event"
	"eth-transactions-scheduler/internal/eth_transaction"
	"eth-transactions-scheduler/pkg/queue_helper/connector"
	"eth-transactions-scheduler/pkg/url_helper"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("eth-transactions-scheduler::ERROR::Error loading .env file")
	}

	requesterInterval, err := strconv.Atoi(os.Getenv("REQUESTER_INTERVAL"))
	if err != nil {
		log.Println("eth-transactions-scheduler::ERROR::requesterInterval: Error converting interval to integer")
	}

	// Queue from which the service will consume the transactions
	ethTransactionsSchedulerQueue := os.Getenv("ETH_TRANSACTIONS_SCHEDULER_QUEUE")
	ethTransactionsSchedulerQueueConn, err := connector.ConnectToQueue(ethTransactionsSchedulerQueue)
	if err != nil {
		log.Printf("eth-transactions-scheduler::ERROR::connForApi::err: %v\n", err)
		os.Exit(1)
	}
	defer ethTransactionsSchedulerQueueConn.Close()

	// Queue to which the service will send the transactions for further recording to the database
	ethTransactionsRecorderQueue := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE")
	ethTransactionsRecorderQueueName := os.Getenv("ETH_TRANSACTIONS_RECORDER_QUEUE_NAME")
	ethTransactionsRecorderQueueConn, err := connector.ConnectToQueue(ethTransactionsRecorderQueue)
	if err != nil {
		log.Printf("eth-transactions-scheduler::ERROR::connForApi::err: %v\n", err)
		os.Exit(1)
	}
	defer ethTransactionsRecorderQueueConn.Close()

	ethTransactionsRecorderQueueCh, err := ethTransactionsRecorderQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-transactions-scheduler::ERROR::failed to open ethTransactionsRecorderQueueCh:", err)
	}
	defer ethTransactionsRecorderQueueCh.Close()

	// Queue to which the service will send the transactions for further recording to the cache
	ethRedisRecorderQueue := os.Getenv("ETH_REDIS_RECORDER_QUEUE")
	ethRedisRecorderQueueName := os.Getenv("ETH_REDIS_RECORDER_QUEUE_NAME")
	ethRedisRecorderQueueConn, err := connector.ConnectToQueue(ethRedisRecorderQueue)
	if err != nil {
		log.Printf("eth-transactions-scheduler::ERROR::connForApi::err: %v\n", err)
		os.Exit(1)
	}
	defer ethRedisRecorderQueueConn.Close()

	ethRedisRecorderQueueCh, err := ethRedisRecorderQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-transactions-scheduler::ERROR::failed to open ethRedisRecorderQueueCh:", err)
	}
	defer ethRedisRecorderQueueCh.Close()

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
		log.Printf("eth-transactions-scheduler:ERROR: RabbitMQ consumer panic: %v\n", err)
	}

	log.Println("eth-blocks-recorder:RabbitMQ:consumer: Consumer is established")

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Printf("eth-blocks-recorder::ERROR::RabbitMQ:consume:error: %v\n", err)
	}
}
