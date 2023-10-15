package eth_transaction

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
	return ebs.ethTransactionRepo.InsertEthTransaction(et)
}
