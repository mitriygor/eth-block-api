package app

import (
	"eth-api/app/handlers"
	"eth-api/app/middleware"
	"eth-api/app/repositories"
	"eth-api/app/routes"
	"eth-api/app/services"
	"eth-api/pkg/mongo_helper"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func NewApp(redisHost string, mongoUrl string, mongoUser string, mongoPassword string, queueForApi string, queueForApiName string) *fiber.App {

	mongoClient, err := mongo_helper.ConnectToMongo(mongoUrl, mongoUser, mongoPassword)
	if err != nil {
		log.Panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	handler := initEthBlockHandler(mongoClient, redisClient, queueForApi, queueForApiName)

	app := fiber.New()

	app.Use(middleware.LoggingMiddleware)

	routes.SetupRoutes(app, handler)

	return app
}

func initEthBlockHandler(mongoClient *mongo.Client, redisClient *redis.Client, queueForApi string, queueForApiName string) *handlers.EthBlockHandler {
	ethBlockRepo := repositories.NewEthBlockRepository(mongoClient, redisClient, queueForApi, queueForApiName)
	ethBlockService := services.NewEthBlockService(ethBlockRepo)
	ethBlockHandler := handlers.NewEthBlockHandler(ethBlockService)
	return ethBlockHandler
}
