package main

import (
	"eth-transactions-store/event"
	"eth-transactions-store/internal/eth_transaction"
	"eth-transactions-store/pkg/mongo_helper"
	"eth-transactions-store/pkg/queue_helper/connector"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR::REQ::Error loading .env file")
	}

	mongoUrl := os.Getenv("MONGO_URL")
	mongoUsername := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	mongoClient, err := mongo_helper.ConnectToMongo(mongoUrl, mongoUsername, mongoPassword)
	if err != nil {
		log.Panic(err)
	}

	queueForStorage := os.Getenv("QUEUE_FOR_STORAGE")
	connForStorage, err := connector.ConnectToQueue(queueForStorage)
	if err != nil {
		log.Printf("ERROR::EthListener::RabbitMQ:connect:error: %v\n", err)
		os.Exit(1)
	}

	log.Println("EthListener::RabbitMQ: Connected to RabbitMQ")
	defer connForStorage.Close()

	redisHost := os.Getenv("REDIS_HOST")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	ethTransactionRepo := eth_transaction.NewEthTransactionRepository(redisClient, connForStorage, mongoClient)
	ethTransactionService := eth_transaction.NewEthTransactionService(ethTransactionRepo)

	consumer, err := event.NewConsumer(connForStorage, ethTransactionService)

	if err != nil {
		log.Println("ERROR::EthListener::RabbitMQ:consumer:PANIC")
		panic(err)
	}

	log.Println("EthListener::RabbitMQ:consumer: Consumer is established")

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Printf("ERROR::EthListener::RabbitMQ:consume:error: %v\n", err)
	}
}
