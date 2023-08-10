package main

import (
	"eth-listener/config"
	"eth-listener/event"
	"eth-listener/internal/eth_block"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	//connWriter := db.ConnectToDB("DSN_WRITER")
	//
	//if connWriter == nil {
	//	log.Panic("Can't connect to Postgres Writer!")
	//}

	//defer connWriter.Close()

	rabbitConn, err := connect()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: "",
		DB:       0,
	})

	log.Println("Listening for and consuming RabbitMQ messages...")

	ethBlockRepo := eth_block.NewEthBlockRepository(nil, redisClient, rabbitConn)
	ethBlockService := eth_block.NewEthBlockService(ethBlockRepo)

	consumer, err := event.NewConsumer(rabbitConn, ethBlockService)

	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
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
			log.Println("Connected to RabbitMQ!")
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
