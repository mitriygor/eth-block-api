package main

import (
	"eth-helpers/queue_helper/connector"
	"eth-redis-recorder/event"
	"eth-redis-recorder/internal/eth_block"
	"eth-redis-recorder/internal/eth_transaction"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("eth-redis-recorder::ERROR::Error loading .env file")
	}

	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		log.Println("eth-blocks-scheduler::ERROR::Error converting interval to integer")
	}

	maxTransactionsPerEthBlock, err := strconv.Atoi(os.Getenv("MAX_TRANSACTIONS_PER_BLOCK"))
	if err != nil {
		log.Println("eth-blocks-scheduler::ERROR::Error converting interval to integer")
	}

	redisEthBlocksQueueName := os.Getenv("ETH_REDIS_ETH_BLOCKS_RECORDER_QUEUE_NAME")
	redisEthTransactionsQueueName := os.Getenv("ETH_REDIS_ETH_TRANSACTIONS_RECORDER_QUEUE_NAME")

	ethRedisRecorderQueue := os.Getenv("ETH_REDIS_RECORDER_QUEUE")
	ethRedisRecorderQueueConn, err := connector.ConnectToQueue(ethRedisRecorderQueue)
	if err != nil {
		log.Printf("eth-redis-recorder::ERROR::RabbitMQ:connect:error: %v\n", err)
		os.Exit(1)
	}

	log.Println("eth-redis-recorder::RabbitMQ: Connected to RabbitMQ")
	defer ethRedisRecorderQueueConn.Close()

	redisHost := os.Getenv("ETH_REDIS")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	ethBlockRepo := eth_block.NewEthBlockRepository(redisClient, cacheSize, redisEthBlocksQueueName)
	ethBlockService := eth_block.NewEthBlockService(ethBlockRepo)

	ethTransactionRepo := eth_transaction.NewEthTransactionRepository(redisClient, cacheSize, maxTransactionsPerEthBlock, redisEthTransactionsQueueName)
	ethTransactionService := eth_transaction.NewEthTransactionService(ethTransactionRepo)

	consumer, err := event.NewConsumer(ethRedisRecorderQueueConn, ethBlockService, ethTransactionService)

	if err != nil {
		log.Println("eth-redis-recorder::ERROR::RabbitMQ:consumer:PANIC")
		panic(err)
	}

	log.Println("eth-redis-recorder::RabbitMQ:consumer: Consumer is established")

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Printf("eth-redis-recorder::ERROR::RabbitMQ:consume:error: %v\n", err)
	}
}
