package event

import (
	"encoding/json"
	"eth-listener/internal/eth_block"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn            *amqp.Connection
	ethBlockService eth_block.Service
}

func NewConsumer(conn *amqp.Connection, ethBlockService eth_block.Service) (Consumer, error) {
	consumer := Consumer{
		conn:            conn,
		ethBlockService: ethBlockService,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		fmt.Printf("Error getting channel: %v\n", err.Error())
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		fmt.Printf("Error declaring queue: %v\n", err.Error())
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			fmt.Printf("Error binding queue: %v\n", err.Error())
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Error consuming queue: %v\n", err.Error())
		return err
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for d := range messages {
				var bd eth_block.BlockDetails
				_ = json.Unmarshal(d.Body, &bd)
				consumer.HandlePayload(bd)
			}
		}()
	}
	wg.Wait()

	return nil
}

func (consumer *Consumer) HandlePayload(bd eth_block.BlockDetails) {
	fmt.Printf("Received payload: %v\n", bd)
}
