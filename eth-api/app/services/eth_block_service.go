package services

import (
	"eth-api/app/repositories"
	"eth-api/pkg/eth_block_helper"
	"log"
)
import "eth-api/app/models"

type EthBlockService interface {
	GetBlockByIdentifierService(hash string) (*models.BlockDetails, error)
}

type ethBlockService struct {
	ethBlockRepo repositories.EthBlockRepository
}

func NewEthBlockService(repo repositories.EthBlockRepository) EthBlockService {
	return &ethBlockService{
		ethBlockRepo: repo,
	}
}

func (ebs *ethBlockService) GetBlockByIdentifierService(blockIdentifier string) (*models.BlockDetails, error) {

	log.Printf("API::GetBlockByIdentifierHandler: %v", blockIdentifier)

	if eth_block_helper.IsInt(blockIdentifier) {
		num := eth_block_helper.StringToInt(blockIdentifier)
		hex := eth_block_helper.IntToHex(num)

		log.Printf("API::GetBlockByIdentifierHandler::number::hex: %v", hex)

		ebs.ethBlockRepo.GetEthBlockByIdentifier(hex, "number")
		return nil, nil
	} else if eth_block_helper.IsHex(blockIdentifier) {

		log.Printf("API::GetBlockByIdentifierHandler::number::blockIdentifier: %v", blockIdentifier)

		ebs.ethBlockRepo.GetEthBlockByIdentifier(blockIdentifier, "number")
		return nil, nil
	} else {

		log.Printf("API::GetBlockByIdentifierHandler::hash::blockIdentifier: %v", blockIdentifier)

		ebs.ethBlockRepo.GetEthBlockByIdentifier(blockIdentifier, "hash")
		return nil, nil
	}

	return nil, nil
}
