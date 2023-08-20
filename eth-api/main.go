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
		log.Printf("eth-api::ERROR::Error loading .env file: %v\n", err)
	}

	redisHost := os.Getenv("ETH_REDIS")

	ethBlocksMongo := os.Getenv("ETH_BLOCKS_MONGO")
	ethBlocksMongoUser := os.Getenv("ETH_BLOCKS_MONGO_USER")
	ethBlocksMongoPassword := os.Getenv("ETH_BLOCKS_MONGO_PASSWORD")
	ethBlocksMongoCredentials := app.MongoCredentials{
		Url:      ethBlocksMongo,
		User:     ethBlocksMongoUser,
		Password: ethBlocksMongoPassword,
	}

	ethBlocksRequesterQueue := os.Getenv("ETH_BLOCKS_REQUESTER_QUEUE")
	ethBlocksRequesterQueueName := os.Getenv("ETH_BLOCKS_REQUESTER_QUEUE_NAME")
	ethBlocksQueueCredentials := app.QueueCredentials{
		Host: ethBlocksRequesterQueue,
		Name: ethBlocksRequesterQueueName,
	}

	ethTransactionsMongo := os.Getenv("ETH_TRANSACTIONS_MONGO")
	ethTransactionsMongoUser := os.Getenv("ETH_TRANSACTIONS_MONGO_USER")
	ethTransactionsMongoPassword := os.Getenv("ETH_TRANSACTIONS_MONGO_PASSWORD")
	ethTransactionMongo := app.MongoCredentials{
		Url:      ethTransactionsMongo,
		User:     ethTransactionsMongoUser,
		Password: ethTransactionsMongoPassword,
	}

	ethTransactionsRequesterQueue := os.Getenv("ETH_TRANSACTIONS_REQUESTER_QUEUE")
	ethTransactionsRequesterQueueName := os.Getenv("ETH_TRANSACTIONS_REQUESTER_QUEUE_NAME")
	ethTransactionsRequesterCredentials := app.QueueCredentials{
		Host: ethTransactionsRequesterQueue,
		Name: ethTransactionsRequesterQueueName,
	}

	port := os.Getenv("PORT")

	app := app.NewApp(redisHost, ethBlocksMongoCredentials, ethTransactionMongo, ethBlocksQueueCredentials, ethTransactionsRequesterCredentials)

	err = app.Listen(":" + port)
	if err != nil {
		log.Printf("eth-api::ERROR::Error launch server: %v\n", err)
	}
}
