package services

import (
	"eth-api/app/models"
	"eth-api/app/repositories"
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
	et, err := ets.ethTransactionRepo.GetTransactionByHash(hash)

	return et, err
}

func (ets *ethTransactionService) GetTransactionsByAddressService(address string) ([]*models.EthTransaction, error) {
	transactions, err := ets.ethTransactionRepo.GetTransactionsByAddress(address)

	return transactions, err
}
