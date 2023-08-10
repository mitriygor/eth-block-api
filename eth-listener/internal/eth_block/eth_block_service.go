package eth_block

type Service interface {
	InsertEthBlockService(bd BlockDetails) error
}

type ethBlockService struct {
	ethBlockRepo Repository
}

func NewEthBlockService(repo Repository) Service {
	return &ethBlockService{
		ethBlockRepo: repo,
	}
}

func (ebs *ethBlockService) InsertEthBlockService(bd BlockDetails) error {
	return ebs.ethBlockRepo.InsertEthBlock(bd)
}
