package repositories

import (
	"eth-api/app/models"
	"gorm.io/gorm"
)

type EthBlockRepository interface {
	GetEthBlockByNumber(blockNumber string) (*models.EthBlock, error)
	CreateEthBlock(createEthBlockDto models.CreateEthBlockDto) (*models.EthBlock, error)
}

type ethBlockRepository struct {
	db *gorm.DB
}

func NewEthBlockRepository(db *gorm.DB) EthBlockRepository {
	return &ethBlockRepository{
		db: db,
	}
}

func (ebr *ethBlockRepository) GetEthBlockByNumber(blockNumber string) (*models.EthBlock, error) {
	var ethBlock models.EthBlock
	err := ebr.db.Where("block_number = ?", blockNumber).First(&ethBlock).Error
	if err != nil {
		//if gorm.IsRecordNotFoundError(err) {
		//	return nil, err
		//}
		return nil, err
	}
	return &ethBlock, nil
}

func (ebr *ethBlockRepository) CreateEthBlock(createEthBlockDto models.CreateEthBlockDto) (*models.EthBlock, error) {
	ethBlock := models.EthBlock{
		BlockNumber: createEthBlockDto.BlockNumber,
		BlockHash:   createEthBlockDto.BlockHash,
	}

	err := ebr.db.Create(&ethBlock).Error
	if err != nil {
		return nil, err
	}

	return &ethBlock, nil
}
