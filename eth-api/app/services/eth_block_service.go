package services

import (
	"eth-api/app/repositories"
	"eth-helpers/eth_block_helper"
	"log"
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

	log.Printf("eth-api:GetLatestEthBlocks")

	ebs.ethBlockRepo.GetLatestEthBlocks()
	return nil, nil
}

func (ebs *ethBlockService) GetBlockByIdentifierService(blockIdentifier string) (*models.BlockDetails, error) {

	log.Printf("eth-api::GetBlockByIdentifierService: %v\n", blockIdentifier)

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

	log.Printf("eth-api::GetBlockByIdentifierService:::identifierType: %v\n", identifierType)

	bd, err = ebs.ethBlockRepo.GetEthBlockByIdentifier(blockIdentifier, identifierType)

	log.Printf("eth-api::GetBlockByIdentifierService:::bd: %v\n", bd)
	log.Printf("eth-api::GetBlockByIdentifierService:::err: %v\n", err)

	return bd, err
}
