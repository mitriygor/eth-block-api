package eth_block

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Repository interface {
	AddBlock(item BlockDetails)
}

type EthBlockRepository struct {
	redisClient *redis.Client
	cacheSize   int
	queueName   string
}

func NewEthBlockRepository(redisClient *redis.Client, cacheSize int, queueName string) Repository {
	return &EthBlockRepository{
		redisClient: redisClient,
		cacheSize:   cacheSize,
		queueName:   queueName,
	}
}

func (ebr *EthBlockRepository) AddBlock(item BlockDetails) {
	ebr.redisClient.LPush(context.Background(), ebr.queueName, item)
	ebr.redisClient.LTrim(context.Background(), ebr.queueName, -int64(ebr.cacheSize), -1)
}
