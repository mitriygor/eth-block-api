package main

import (
	"eth-blocks-scheduler/internal/eth_block"
	"eth-helpers/cron_job"
	"eth-helpers/queue_helper/connector"
	"eth-helpers/url_helper"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("eth-blocks-scheduler::Error loading .env file")
	}

	schedulerInterval, err := strconv.Atoi(os.Getenv("SCHEDULER_INTERVAL"))
	if err != nil {
		log.Println("eth-blocks-scheduler::ERROR::Error converting schedulerInterval to integer")
		schedulerInterval = 120
	}

	requesterInterval, err := strconv.Atoi(os.Getenv("REQUESTER_INTERVAL"))
	if err != nil {
		log.Println("eth-blocks-scheduler::ERROR::Error converting interval to integer")
		schedulerInterval = 10
	}

	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		log.Println("eth-blocks-scheduler::ERROR::Error converting cacheSize to integer")
		cacheSize = 50
	}

	// Queue to which the service will send the blocks for further recording to the database
	ethBlocksRecorderQueue := os.Getenv("ETH_BLOCKS_RECORDER_QUEUE")
	ethBlocksRecorderQueueName := os.Getenv("ETH_BLOCKS_RECORDER_QUEUE_NAME")
	ethBlocksRecorderQueueConn, err := connector.ConnectToQueue(ethBlocksRecorderQueue)
	if err != nil {
		log.Printf("eth-blocks-scheduler::ERROR::RabbitMQ: %v\n", err)
		os.Exit(1)
	}

	ethBlocksRecorderQueueCh, err := ethBlocksRecorderQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-blocks-scheduler::ERROR::failed to open ethBlocksRecorderQueueCh:", err)
	}
	defer ethBlocksRecorderQueueCh.Close()
	defer ethBlocksRecorderQueueConn.Close()

	// Queue to which the service will send the blocks for further recording to the Redis
	ethRedisRecorderQueue := os.Getenv("ETH_REDIS_RECORDER_QUEUE")
	ethRedisRecorderQueueName := os.Getenv("ETH_REDIS_RECORDER_QUEUE_NAME")
	ethRedisRecorderQueueConn, err := connector.ConnectToQueue(ethRedisRecorderQueue)
	if err != nil {
		log.Printf("eth-blocks-scheduler::ERROR::RabbitMQ: %v\n", err)
		os.Exit(1)
	}

	ethRedisRecorderQueueCh, err := ethRedisRecorderQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-blocks-scheduler::ERROR::failed to open ethRedisRecorderQueueCh:", err)
	}
	defer ethRedisRecorderQueueCh.Close()
	defer ethRedisRecorderQueueConn.Close()

	// Queue to which the service will send the transactions for further scheduling
	ethTransactionsSchedulerQueue := os.Getenv("ETH_TRANSACTIONS_SCHEDULER_QUEUE")
	ethTransactionsSchedulerQueueName := os.Getenv("ETH_TRANSACTIONS_SCHEDULER_QUEUE_NAME")
	ethTransactionsSchedulerQueueConn, err := connector.ConnectToQueue(ethTransactionsSchedulerQueue)
	if err != nil {
		log.Printf("eth-blocks-scheduler::ERROR::RabbitMQ: %v\n", err)
		os.Exit(1)
	}

	ethTransactionsSchedulerQueueCh, err := ethTransactionsSchedulerQueueConn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "eth-blocks-scheduler::ERROR::failed to open ethTransactionsSchedulerQueueCh:", err)
	}
	defer ethTransactionsSchedulerQueueCh.Close()
	defer ethTransactionsSchedulerQueueConn.Close()

	endpoint := os.Getenv("HTTP_ENDPOINT")
	version := os.Getenv("HTTP_ENDPOINT_VERSION")
	secretKey := os.Getenv("HTTP_KEY")
	url := url_helper.GetUrl(endpoint, version, secretKey)

	jsonRpc := os.Getenv("JSONRPC")

	log.Printf("schedulerInterval: %d\nrequesterInterval: %v\nendpoint: %s\nversion: %s\nsecretKey: %s\njsonRpc: %s\n", schedulerInterval, requesterInterval, endpoint, version, secretKey, jsonRpc)

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
