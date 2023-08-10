package eth_block

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Repository interface {
	CreateEthBlock(bd BlockDetails) error
}

type EthBlockRepository struct {
	DBW         *sql.DB
	rabbitConn  *amqp.Connection
	redisClient *redis.Client
}

func NewEthBlockRepository(dbw *sql.DB, redisClient *redis.Client, rabbitConn *amqp.Connection) Repository {

	return &EthBlockRepository{
		DBW:         dbw,
		rabbitConn:  rabbitConn,
		redisClient: redisClient,
	}
}

func (ebr *EthBlockRepository) CreateEthBlock(bd BlockDetails) error {
	fmt.Printf("CreateEthBlock::bd: %v\n", bd)
	return nil
}
