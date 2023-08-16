package main

import (
	"eth-blocks-cron/internal/eth_block"
	"eth-blocks-cron/pkg/cron_job"
	"eth-blocks-cron/pkg/queue_helper/connector"
	"eth-blocks-cron/pkg/url_helper"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	interval, err := strconv.Atoi(os.Getenv("INTERVAL"))
	if err != nil {
		log.Fatal("Error converting interval to integer")
	}

	queueForStorage := os.Getenv("QUEUE_FOR_STORAGE")
	connForStorage, err := connector.ConnectToQueue(queueForStorage)
	if err != nil {
		log.Printf("EthEmitter::RabbitMQ: %v\n", err)
		os.Exit(1)
	}

	log.Println("EthEmitter::RabbitMQ: Connected to RabbitMQ")

	defer connForStorage.Close()

	redisHost := os.Getenv("REDIS_HOST")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	endpoint := os.Getenv("HTTP_ENDPOINT")
	version := os.Getenv("HTTP_ENDPOINT_VERSION")
	secretKey := os.Getenv("HTTP_KEY")
	url := url_helper.GetUrl(endpoint, version, secretKey)

	jsonRpc := os.Getenv("JSONRPC")

	log.Printf("interval: %d\nendpoint: %s\nversion: %s\nsecretKey: %s\njsonRpc: %s\n", interval, endpoint, version, secretKey, jsonRpc)

	ethBlockRepo := eth_block.NewEthBlockRepository(redisClient, connForStorage)
	ethBlockService := eth_block.NewEthBlockService(ethBlockRepo, url, jsonRpc)

	for {
		interval := cron_job.GetInterval(interval)
		ticker := time.NewTicker(interval)
		ethBlockService.PushLatestBlocks()
		<-ticker.C
		ticker.Stop()
	}
}
