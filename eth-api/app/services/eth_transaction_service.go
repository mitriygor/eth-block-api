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
	log.Printf("API::GetTransactionByHashService:: hash: %v\n", hash)
	et, err := ets.ethTransactionRepo.GetTransactionByHash(hash)

	log.Printf("API::GetTransactionByHashService::err: %v\n", err)
	log.Printf("API::GetTransactionByHashService::et: %v\n", et)

	return et, err
}

func (ets *ethTransactionService) GetTransactionsByAddressService(address string) ([]*models.EthTransaction, error) {
	return nil, nil
}
