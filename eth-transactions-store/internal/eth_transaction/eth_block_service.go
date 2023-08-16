package eth_transaction

import "log"

type Service interface {
	InsertEthTransactionService(bd EthTransaction) error
}

type ethTransactionService struct {
	ethTransactionRepo Repository
}

func NewEthTransactionService(repo Repository) Service {
	return &ethTransactionService{
		ethTransactionRepo: repo,
	}
}

func (ebs *ethTransactionService) InsertEthTransactionService(et EthTransaction) error {
	log.Printf("ERROR::EthTransactionsStore::InsertEthTransactionService::et: %v\n", et)

	return ebs.ethTransactionRepo.InsertEthTransaction(et)
}
