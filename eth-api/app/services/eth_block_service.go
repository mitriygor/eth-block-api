package services

import "eth-api/app/repositories"
import "eth-api/app/models"

type EthBlockService interface {
	GetEthBlockByNumberService(blockNumber string) (*models.EthBlock, error)
	CreateEthBlockService(dto models.CreateEthBlockDto) (*models.EthBlock, error)
}

type ethBlockService struct {
	ethBlockRepo repositories.EthBlockRepository
}

func NewEthBlockService(repo repositories.EthBlockRepository) EthBlockService {
	return &ethBlockService{
		ethBlockRepo: repo,
	}
}

func (ebs *ethBlockService) GetEthBlockByNumberService(blockNumber string) (*models.EthBlock, error) {
	return ebs.ethBlockRepo.GetEthBlockByNumber(blockNumber)
}

func (ebs *ethBlockService) CreateEthBlockService(dto models.CreateEthBlockDto) (*models.EthBlock, error) {
	return ebs.ethBlockRepo.CreateEthBlock(dto)
}
