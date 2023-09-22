package app

import (
	"eth-api/app/handlers"
	"eth-api/app/middleware"
	"eth-api/app/repositories"
	"eth-api/app/routes"
	"eth-api/app/services"
	"eth-helpers/mongo_helper"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type MongoCredentials struct {
	Url      string
	User     string
	Password string
}

type QueueCredentials struct {
	Host string
	Name string
}

func NewApp(redisHost string, ethBlocksMongoCredentials MongoCredentials, ethTransactionMongo MongoCredentials, ethBlocksQueueCredentials QueueCredentials, ethTransactionsRequesterCredentials QueueCredentials) *fiber.App {

	ethBlockMongoClient, err := mongo_helper.ConnectToMongo(ethBlocksMongoCredentials.Url, ethBlocksMongoCredentials.User, ethBlocksMongoCredentials.Password)
	if err != nil {
		log.Panic(err)
	}

	ethTransactionMongoClient, err := mongo_helper.ConnectToMongo(ethTransactionMongo.Url, ethTransactionMongo.User, ethTransactionMongo.Password)
	if err != nil {
		log.Panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	ethBlockHandler := initEthBlockHandler(ethBlockMongoClient, redisClient, ethBlocksQueueCredentials.Host, ethBlocksQueueCredentials.Name)
	ethTransactionHandler := initEthTransactionHandler(ethTransactionMongoClient, redisClient, ethTransactionsRequesterCredentials.Host, ethTransactionsRequesterCredentials.Name)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "*",
	}))

	app.Use(middleware.LoggingMiddleware)

	routes.SetupRoutes(app, ethBlockHandler, ethTransactionHandler)

	return app
}

func initEthBlockHandler(mongoClient *mongo.Client, redisClient *redis.Client, queueHost string, queueName string) *handlers.EthBlockHandler {
	ethBlockRepo := repositories.NewEthBlockRepository(mongoClient, redisClient, queueHost, queueName)
	ethBlockService := services.NewEthBlockService(ethBlockRepo)
	ethBlockHandler := handlers.NewEthBlockHandler(ethBlockService)
	return ethBlockHandler
}

func initEthTransactionHandler(mongoClient *mongo.Client, redisClient *redis.Client, queueHost string, queueName string) *handlers.EthTransactionHandler {
	ethTransactionRepo := repositories.NewEthTransactionRepository(mongoClient, redisClient, queueHost, queueName)
	ethTransactionService := services.NewEthTransactionService(ethTransactionRepo)
	ethTransactionHandler := handlers.NewEthTransactionHandler(ethTransactionService)
	return ethTransactionHandler
}
