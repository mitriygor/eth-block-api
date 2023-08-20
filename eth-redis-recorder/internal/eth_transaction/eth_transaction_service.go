package eth_transaction

type Service interface {
	AddTransactionService(item EthTransaction)
}

type ethTransactionService struct {
	ethTransactionRepo Repository
}

func NewEthTransactionService(repo Repository) Service {
	return &ethTransactionService{
		ethTransactionRepo: repo,
	}
}

func (es *ethTransactionService) AddTransactionService(et EthTransaction) {
	es.ethTransactionRepo.AddBlock(et)
}
