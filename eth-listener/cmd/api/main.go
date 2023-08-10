package main

import (
	"context"

	"eth-listener/config"
	"eth-listener/event"
	"eth-listener/internal/eth_block"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	rabbitConn, err := connect()

	if err != nil {
		log.Printf("ERROR::EthListener::RabbitMQ:connect:error: %v\n", err)
		os.Exit(1)
	}

	log.Println("EthListener::RabbitMQ: Connected to RabbitMQ")

	defer rabbitConn.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "",
		DB:       0,
	})

	ethBlockRepo := eth_block.NewEthBlockRepository(redisClient, rabbitConn, mongoClient)
	ethBlockService := eth_block.NewEthBlockService(ethBlockRepo)

	consumer, err := event.NewConsumer(rabbitConn, ethBlockService)

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

func connectToMongo() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(config.MongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password123",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	return c, nil
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
