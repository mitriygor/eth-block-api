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

	mongoUrl := os.Getenv("MONGO_URL")
	mongoUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	queueForApi := os.Getenv("QUEUE_FOR_API")
	queueForApiName := os.Getenv("QUEUE_FOR_API_NAME")

	port := os.Getenv("PORT")

	app := app.NewApp(redisHost, mongoUrl, mongoUser, mongoPassword, queueForApi, queueForApiName)

	err = app.Listen(":" + port)
	if err != nil {
		log.Printf("ERROR::API::Error launch server: %v\n", err)
	}
}
