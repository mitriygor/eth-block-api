package eth_block

type Service interface {
	AddBlockService(bd BlockDetails)
}

type ethBlockService struct {
	ethBlockRepo Repository
}

func NewEthBlockService(repo Repository) Service {
	return &ethBlockService{
		ethBlockRepo: repo,
	}
}

func (es *ethBlockService) AddBlockService(bd BlockDetails) {
	es.ethBlockRepo.AddBlock(bd)
}
