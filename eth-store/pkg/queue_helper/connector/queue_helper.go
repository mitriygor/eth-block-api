package connector

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"time"
)

func ConnectToQueue(host string) (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial(host)
		if err != nil {
			fmt.Printf("ERROR::REQ::connectToQueue::err: %v\n", err)
			counts++
		} else {
			fmt.Println("REQ::connectToQueue::CONNECTED")
			connection = c
			break
		}

		if counts > 10 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("REQ::connectToQueue::backOff: %v\n", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
