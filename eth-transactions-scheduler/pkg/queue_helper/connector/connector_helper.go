package connector

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"time"
)

func ConnectToQueue(host string) (*amqp.Connection, error) {
	log.Printf("eth-transactions-scheduler::ConnectToQueue::host: %v\n", host)

	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial(host)
		if err != nil {
			log.Printf("eth-transactions-scheduler::ERROR::ConnectToQueue::err: %v\n", err)
			counts++
		} else {
			log.Println("eth-transactions-scheduler::ConnectToQueue::CONNECTED")
			connection = c
			break
		}

		if counts > 10 {
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("eth-transactions-scheduler::ConnectToQueue::backOff: %v\n", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
