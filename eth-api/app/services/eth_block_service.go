package services

import (
	"eth-api/app/helpers/logger"
	"eth-api/app/repositories"
	"eth-helpers/eth_block_helper"
)
import "eth-api/app/models"

type EthBlockService interface {
	GetBlockByIdentifierService(hash string) (*models.BlockDetails, error)
	GetLatestEthBlocks() ([]*models.BlockDetails, error)
}

type ethBlockService struct {
	ethBlockRepo repositories.EthBlockRepository
}

func NewEthBlockService(repo repositories.EthBlockRepository) EthBlockService {
	return &ethBlockService{
		ethBlockRepo: repo,
	}
}

func (ebs *ethBlockService) GetLatestEthBlocks() ([]*models.BlockDetails, error) {

	latestBlockDetails, err := ebs.ethBlockRepo.GetLatestEthBlocks()

	if err != nil || len(latestBlockDetails) == 0 {
		logger.Error("eth-api::ERROR::GetLatestEthBlocks", "error", err)
		return nil, err
	}

	return latestBlockDetails, nil
}

func (ebs *ethBlockService) GetBlockByIdentifierService(blockIdentifier string) (*models.BlockDetails, error) {
	var bd *models.BlockDetails
	var err error

	identifierType := "hash"

	if eth_block_helper.IsInt(blockIdentifier) {
		num := eth_block_helper.StringToInt(blockIdentifier)
		blockIdentifier = eth_block_helper.IntToHex(num)
		identifierType = "number"
	} else if eth_block_helper.IsHex(blockIdentifier) {
		identifierType = "number"
	}

	bd, err = ebs.ethBlockRepo.GetEthBlockByIdentifier(blockIdentifier, identifierType)

	if err != nil {
		logger.Error("eth-api::ERROR::GetBlockByIdentifierService", "error", err)
	}

	return bd, err
}
