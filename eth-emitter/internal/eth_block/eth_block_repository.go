package eth_block

import (
	"encoding/json"
	"eth-emitter/event"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Repository interface {
	PushEthBlock(bd BlockDetails) error
	GetCurrentBlockNumber() int
	SetCurrentBlockNumber(blockNumber int)
}

type EthBlockRepository struct {
	rabbitConn  *amqp.Connection
	redisClient *redis.Client
	emitter     *event.Emitter
}

func NewEthBlockRepository(redisClient *redis.Client, rabbitConn *amqp.Connection) Repository {

	emitter, err := event.NewEmitter(rabbitConn)

	if err != nil {
		return nil
	}

	return &EthBlockRepository{
		rabbitConn:  rabbitConn,
		redisClient: redisClient,
		emitter:     emitter,
	}
}

func (ebr *EthBlockRepository) PushEthBlock(bd BlockDetails) error {
	j, _ := json.MarshalIndent(&bd, "", "\t")
	err := ebr.emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

func (ebr *EthBlockRepository) GetCurrentBlockNumber() int {
	blockNumber, err := ebr.redisClient.Get(ebr.redisClient.Context(), "current_block_number").Int()
	if err != nil {
		return -1
	}

	return blockNumber
}

func (ebr *EthBlockRepository) SetCurrentBlockNumber(blockNumber int) {
	ebr.redisClient.Set(ebr.redisClient.Context(), "current_block_number", blockNumber, 0)
}
