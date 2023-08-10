package models

import "gorm.io/gorm"

type EthBlock struct {
	gorm.Model
	Number           uint64 `gorm:"uniqueIndex;not null"`
	Hash             string `gorm:"uniqueIndex;not null"`
	ParentHash       string
	Nonce            string
	Sha3Uncles       string
	LogsBloom        string
	TransactionsRoot string
	StateRoot        string
	Miner            string
	Difficulty       uint64
	TotalDifficulty  uint64
	ExtraData        string
	Size             uint64
	GasLimit         uint64
	GasUsed          uint64
	Timestamp        uint64
	Transactions     []string `gorm:"-"`
	Uncles           []string `gorm:"-"`
}

type CreateEthBlockDto struct {
	Number     uint64 `json:"number" binding:"required"`
	Hash       string `json:"hash" binding:"required"`
	ParentHash string `json:"parentHash"`
	Nonce      string `json:"nonce"`
}
