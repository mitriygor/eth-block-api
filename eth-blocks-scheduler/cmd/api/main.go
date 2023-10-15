package main

import (
	"eth-blocks-scheduler/internal/eth_block"
	"eth-blocks-scheduler/pkg/logger"
	"eth-helpers/cron_job"
	"eth-helpers/queue_helper/connector"
	"eth-helpers/url_helper"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

func main() {
	err := logger.Initialize("info")
	if err != nil {
		panic(err)
	}
	defer func(Log *zap.SugaredLogger) {
		err := Log.Sync()
		if err != nil {
			logger.Error("eth-blocks-scheduler:ERROR:sync", "error", err)
		}
	}(logger.Log)

	err = godotenv.Load()
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:loading .env file", "error", err)
	}

	schedulerInterval, err := strconv.Atoi(os.Getenv("SCHEDULER_INTERVAL"))
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:converting schedulerInterval to integer", "error", err)
		schedulerInterval = 120
	}

	requesterInterval, err := strconv.Atoi(os.Getenv("REQUESTER_INTERVAL"))
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:converting requesterInterval to integer", "error", err)
		schedulerInterval = 10
	}

	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:converting cacheSize to integer", "error", err)
		cacheSize = 50
	}

	// Queue to which the service will send the blocks for further recording to the database
	ethBlocksRecorderQueue := os.Getenv("ETH_BLOCKS_RECORDER_QUEUE")
	ethBlocksRecorderQueueName := os.Getenv("ETH_BLOCKS_RECORDER_QUEUE_NAME")
	ethBlocksRecorderQueueConn, err := connector.ConnectToQueue(ethBlocksRecorderQueue)
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:connecting to ethBlocksRecorderQueue", "error", err)
		os.Exit(1)
	}

	ethBlocksRecorderQueueCh, err := ethBlocksRecorderQueueConn.Channel()
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:opening ethBlocksRecorderQueueCh", "error", err)
	}
	defer func(ethBlocksRecorderQueueCh *amqp.Channel) {
		err := ethBlocksRecorderQueueCh.Close()
		if err != nil {
			logger.Error("eth-blocks-scheduler:ERROR:close ethBlocksRecorderQueueCh", "error", err)
		}
	}(ethBlocksRecorderQueueCh)
	defer func(ethBlocksRecorderQueueConn *amqp.Connection) {
		err := ethBlocksRecorderQueueConn.Close()
		if err != nil {
			logger.Error("eth-blocks-scheduler:ERROR:close ethBlocksRecorderQueueConn", "error", err)
		}
	}(ethBlocksRecorderQueueConn)

	// Queue to which the service will send the blocks for further recording to the Redis
	ethRedisRecorderQueue := os.Getenv("ETH_REDIS_RECORDER_QUEUE")
	ethRedisRecorderQueueName := os.Getenv("ETH_REDIS_RECORDER_QUEUE_NAME")
	ethRedisRecorderQueueConn, err := connector.ConnectToQueue(ethRedisRecorderQueue)
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:connecting to ethRedisRecorderQueue", "error", err)
		os.Exit(1)
	}

	ethRedisRecorderQueueCh, err := ethRedisRecorderQueueConn.Channel()
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:opening ethRedisRecorderQueueCh", "error", err)
	}
	defer func(ethRedisRecorderQueueCh *amqp.Channel) {
		err := ethRedisRecorderQueueCh.Close()
		if err != nil {
			logger.Error("eth-blocks-scheduler:ERROR:close ethRedisRecorderQueueCh", "error", err)
		}
	}(ethRedisRecorderQueueCh)
	defer func(ethRedisRecorderQueueConn *amqp.Connection) {
		err := ethRedisRecorderQueueConn.Close()
		if err != nil {
			logger.Error("eth-blocks-scheduler:ERROR:close ethRedisRecorderQueueConn", "error", err)
		}
	}(ethRedisRecorderQueueConn)

	// Queue to which the service will send the transactions for further scheduling
	ethTransactionsSchedulerQueue := os.Getenv("ETH_TRANSACTIONS_SCHEDULER_QUEUE")
	ethTransactionsSchedulerQueueName := os.Getenv("ETH_TRANSACTIONS_SCHEDULER_QUEUE_NAME")
	ethTransactionsSchedulerQueueConn, err := connector.ConnectToQueue(ethTransactionsSchedulerQueue)
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:connecting to ethTransactionsSchedulerQueue", "error", err)
		os.Exit(1)
	}

	ethTransactionsSchedulerQueueCh, err := ethTransactionsSchedulerQueueConn.Channel()
	if err != nil {
		logger.Error("eth-blocks-scheduler:ERROR:opening ethTransactionsSchedulerQueueCh", "error", err)
	}
	defer func(ethTransactionsSchedulerQueueCh *amqp.Channel) {
		err := ethTransactionsSchedulerQueueCh.Close()
		if err != nil {
			logger.Error("eth-blocks-scheduler:ERROR:close ethTransactionsSchedulerQueueCh", "error", err)
		}
	}(ethTransactionsSchedulerQueueCh)
	defer func(ethTransactionsSchedulerQueueConn *amqp.Connection) {
		err := ethTransactionsSchedulerQueueConn.Close()
		if err != nil {
			logger.Error("eth-blocks-scheduler:ERROR:close ethTransactionsSchedulerQueueConn", "error", err)
		}
	}(ethTransactionsSchedulerQueueConn)

	endpoint := os.Getenv("HTTP_ENDPOINT")
	version := os.Getenv("HTTP_ENDPOINT_VERSION")
	secretKey := os.Getenv("HTTP_KEY")
	url := url_helper.GetUrl(endpoint, version, secretKey)

	jsonRpc := os.Getenv("JSONRPC")

	ethBlockRepo := eth_block.NewEthBlockRepository(ethBlocksRecorderQueueCh, ethBlocksRecorderQueueName, ethRedisRecorderQueueCh, ethRedisRecorderQueueName, ethTransactionsSchedulerQueueCh, ethTransactionsSchedulerQueueName)
	ethBlockService := eth_block.NewEthBlockService(ethBlockRepo, url, jsonRpc, cacheSize, requesterInterval)

	for {
		interval := cron_job.GetInterval(schedulerInterval)
		ticker := time.NewTicker(interval)
		ethBlockService.PushLatestBlocks()
		<-ticker.C
		ticker.Stop()
	}
}
