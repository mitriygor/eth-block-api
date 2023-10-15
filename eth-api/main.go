package main

import (
	"eth-api/app"
	"eth-api/app/helpers/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

func main() {
	// Initialize the logger
	err := logger.Initialize("info")
	if err != nil {
		panic(err)
	}
	defer func(Log *zap.SugaredLogger) {
		err := Log.Sync()
		if err != nil {
			logger.Error("eth-api::main::ERROR::sync", "error", err)
		}
	}(logger.Log)

	err = godotenv.Load()
	if err != nil {
		logger.Error("eth-api::main::ERROR::env", "error", err)
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
		logger.Error("eth-api::main::ERROR::launch", "error", err)
	}
}
