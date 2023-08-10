package eth_block

type Service interface {
	CreateEthBlockService(bd BlockDetails) error
}

type ethBlockService struct {
	ethBlockRepo Repository
}

func NewEthBlockService(repo Repository) Service {
	return &ethBlockService{
		ethBlockRepo: repo,
	}
}

func (ebs *ethBlockService) CreateEthBlockService(bd BlockDetails) error {
	return ebs.ethBlockRepo.CreateEthBlock(bd)
}
