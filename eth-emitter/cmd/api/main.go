package main

import (
	"eth-emitter/config"
	"eth-emitter/internal/eth_block"
	"eth-emitter/pkg/cron_job"
	"eth-emitter/pkg/url_helper"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	id, err := strconv.Atoi(os.Getenv("INTERVAL"))
	if err != nil {
		log.Fatal("Error converting id to integer")
	}

	interval, err := strconv.Atoi(os.Getenv("INTERVAL"))
	if err != nil {
		log.Fatal("Error converting interval to integer")
	}

	rabbitConn, err := connect()
	if err != nil {
		log.Printf("EthEmitter::RabbitMQ: %v\n", err)
		os.Exit(1)
	}

	log.Println("EthEmitter::RabbitMQ: Connected to RabbitMQ")

	defer rabbitConn.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "",
		DB:       0,
	})

	endpoint := os.Getenv("HTTP_ENDPOINT")
	version := os.Getenv("HTTP_ENDPOINT_VERSION")
	secretKey := os.Getenv("HTTP_KEY")
	url := url_helper.GetUrl(endpoint, version, secretKey)

	jsonRpc := os.Getenv("JSONRPC")

	log.Printf("id: %d\ninterval: %d\nendpoint: %s\nversion: %s\nsecretKey: %s\njsonRpc: %s\n", id, interval, endpoint, version, secretKey, jsonRpc)

	ethBlockRepo := eth_block.NewEthBlockRepository(redisClient, rabbitConn)
	ethBlockService := eth_block.NewEthBlockService(ethBlockRepo, url, id, jsonRpc)

	for {
		interval := cron_job.GetInterval(interval)
		ticker := time.NewTicker(interval)
		ethBlockService.PushLatestBlocks()
		<-ticker.C
		ticker.Stop()
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			log.Printf("Failed to connect to RabbitMQ after %d retries\n", counts)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
