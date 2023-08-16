package main

import (
	"eth-api/app"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("ERROR::API::Error loading .env file: %v\n", err)
	}

	redisHost := os.Getenv("REDIS_HOST")

	ethBlockMongoUrl := os.Getenv("ETH_BLOCK_MONGO_URL")
	ethBlockMongoUser := os.Getenv("ETH_BLOCK_MONGO_USER")
	ethBlockMongoPassword := os.Getenv("ETH_BLOCK_MONGO_PASSWORD")
	ethBlockMongo := app.MongoCredentials{
		Url:      ethBlockMongoUrl,
		User:     ethBlockMongoUser,
		Password: ethBlockMongoPassword,
	}

	ethBlockQueueHost := os.Getenv("ETH_BLOCK_QUEUE_HOST")
	ethBlockQueueName := os.Getenv("ETH_BLOCK_QUEUE_NAME")
	ethBlockQueue := app.QueueCredentials{
		Host: ethBlockQueueHost,
		Name: ethBlockQueueName,
	}

	ethTransactionMongoUrl := os.Getenv("ETH_TRANSACTION_MONGO_URL")
	ethTransactionMongoUser := os.Getenv("ETH_TRANSACTION_MONGO_USER")
	ethTransactionMongoPassword := os.Getenv("ETH_TRANSACTION_MONGO_PASSWORD")
	ethTransactionMongo := app.MongoCredentials{
		Url:      ethTransactionMongoUrl,
		User:     ethTransactionMongoUser,
		Password: ethTransactionMongoPassword,
	}

	ethTransactionQueueHost := os.Getenv("ETH_TRANSACTION_QUEUE_HOST")
	ethTransactionQueueName := os.Getenv("ETH_TRANSACTION_QUEUE_NAME")
	ethTransactionQueue := app.QueueCredentials{
		Host: ethTransactionQueueHost,
		Name: ethTransactionQueueName,
	}

	port := os.Getenv("PORT")

	app := app.NewApp(redisHost, ethBlockMongo, ethTransactionMongo, ethBlockQueue, ethTransactionQueue)

	err = app.Listen(":" + port)
	if err != nil {
		log.Printf("ERROR::API::Error launch server: %v\n", err)
	}
}
