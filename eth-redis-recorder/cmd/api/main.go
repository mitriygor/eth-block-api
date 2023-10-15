package main

import (
	"eth-helpers/queue_helper/connector"
	"eth-redis-recorder/event"
	"eth-redis-recorder/internal/eth_block"
	"eth-redis-recorder/internal/eth_transaction"
	"eth-redis-recorder/pkg/logger"
	"github.com/go-redis/redis/v8"
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
			logger.Error("eth-blocks-recorder:ERROR:sync", "error", err)
		}
	}(logger.Log)

	err = godotenv.Load()
	if err != nil {
		logger.Error("eth-redis-recorder::ERROR::Error loading .env file", "error", err)
	}

	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		logger.Error("eth-redis-recorder::ERROR::Error converting interval to integer", "error", err)
	}

	maxTransactionsPerEthBlock, err := strconv.Atoi(os.Getenv("MAX_TRANSACTIONS_PER_BLOCK"))
	if err != nil {
		logger.Error("eth-redis-recorder::ERROR::Error converting interval to integer", "error", err)
	}

	redisEthBlocksQueueName := os.Getenv("ETH_REDIS_ETH_BLOCKS_RECORDER_QUEUE_NAME")
	redisEthTransactionsQueueName := os.Getenv("ETH_REDIS_ETH_TRANSACTIONS_RECORDER_QUEUE_NAME")

	ethRedisRecorderQueue := os.Getenv("ETH_REDIS_RECORDER_QUEUE")
	ethRedisRecorderQueueConn, err := connector.ConnectToQueue(ethRedisRecorderQueue)
	if err != nil {
		logger.Error("eth-redis-recorder::ERROR::RabbitMQ:connect:error", "error", err)
		os.Exit(1)
	}

	defer func(ethRedisRecorderQueueConn *amqp.Connection) {
		err := ethRedisRecorderQueueConn.Close()
		if err != nil {
			logger.Error("eth-redis-recorder::ERROR::RabbitMQ:connection:close:error", "error", err)
		}
	}(ethRedisRecorderQueueConn)

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
		logger.Error("eth-redis-recorder::ERROR::RabbitMQ:consumer:PANIC", "error", err)
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO"})
	if err != nil {
		logger.Error("eth-redis-recorder::ERROR::RabbitMQ:consume:error", "error", err)
	}
}
