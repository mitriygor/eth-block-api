package eth_transaction

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Repository interface {
	AddBlock(item EthTransaction)
}

type EthTransactionRepository struct {
	redisClient                *redis.Client
	cacheSize                  int
	maxTransactionsPerEthBlock int
	queueName                  string
}

func NewEthTransactionRepository(redisClient *redis.Client, cacheSize int, maxTransactionsPerEthBlock int, queueName string) Repository {
	return &EthTransactionRepository{
		redisClient:                redisClient,
		cacheSize:                  cacheSize,
		maxTransactionsPerEthBlock: maxTransactionsPerEthBlock,
		queueName:                  queueName,
	}
}

func (etr *EthTransactionRepository) AddBlock(item EthTransaction) {
	queueSize := int64(etr.cacheSize * etr.maxTransactionsPerEthBlock)
	etr.redisClient.LPush(context.Background(), etr.queueName, item)
	etr.redisClient.LTrim(context.Background(), etr.queueName, -queueSize, -1)
}
