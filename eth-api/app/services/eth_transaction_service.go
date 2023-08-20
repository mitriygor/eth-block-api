package services

import (
	"eth-api/app/models"
	"eth-api/app/repositories"
	"log"
)

type EthTransactionService interface {
	GetTransactionByHashService(hash string) (*models.EthTransaction, error)
	GetTransactionsByAddressService(address string) ([]*models.EthTransaction, error)
}

type ethTransactionService struct {
	ethTransactionRepo repositories.EthTransactionRepository
}

func NewEthTransactionService(repo repositories.EthTransactionRepository) EthTransactionService {
	return &ethTransactionService{
		ethTransactionRepo: repo,
	}
}

func (ets *ethTransactionService) GetTransactionByHashService(hash string) (*models.EthTransaction, error) {
	log.Printf("eth-api::GetTransactionByHashService:: hash: %v\n", hash)
	et, err := ets.ethTransactionRepo.GetTransactionByHash(hash)

	log.Printf("eth-api::GetTransactionByHashService::err: %v\n", err)
	log.Printf("eth-api::GetTransactionByHashService::et: %v\n", et)

	return et, err
}

func (ets *ethTransactionService) GetTransactionsByAddressService(address string) ([]*models.EthTransaction, error) {
	log.Printf("eth-api::GetTransactionsByAddressService:: address: %v\n", address)

	transactions, err := ets.ethTransactionRepo.GetTransactionsByAddress(address)

	log.Printf("eth-api::GetTransactionsByAddressService:: err: %v\n", err)
	log.Printf("eth-api::GetTransactionsByAddressService::transactions: %v\n", transactions)

	return transactions, err
}
