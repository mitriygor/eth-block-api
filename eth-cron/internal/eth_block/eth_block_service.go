package eth_block

import (
	"eth-cron/pkg/json_helper"
	"eth-cron/pkg/url_helper"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Service interface {
	PushLatestBlocks()
	PushBlock(blockDetails BlockDetails) error
	GetCurrentBlockNumber() int
	GetLatestBlockNumber() int
	GetBlockByNumber(blockNumber int) (BlockResponse, error)
	SetCurrentBlockNumber(blockNumber int)
	HexToInt(hexStr string) int
	IntToHex(num int) string
}

type ethBlockService struct {
	ethBlockRepo Repository
	url          string
	jsonRpc      string
}

func NewEthBlockService(repo Repository, url string, jsonRpc string) Service {
	return &ethBlockService{
		ethBlockRepo: repo,
		url:          url,
		jsonRpc:      jsonRpc,
	}
}

func (ebs *ethBlockService) PushLatestBlocks() {
	curr := ebs.GetCurrentBlockNumber()
	latest := ebs.GetLatestBlockNumber()

	if latest == -1 {
		log.Println("Issue to get blocks numbers")
		return
	}

	if curr == latest || curr > latest {
		log.Println("No blocks to push")
		return
	}

	if curr == -1 {
		curr = latest - 1
	}

	for i := curr; i < latest; i++ {
		blockResponse, err := ebs.GetBlockByNumber(i)

		if err != nil {
			log.Printf("PushLatestBlocks::blockResponse::error:%v\n", err)
			continue
		}

		err = ebs.PushBlock(blockResponse.Result)

		if err != nil {
			log.Printf("%v\n", err)
			continue
		}

		num := ebs.HexToInt(blockResponse.Result.Number)

		if num > curr {
			curr = num
			ebs.SetCurrentBlockNumber(curr)
		}

		return
	}
}

func (ebs *ethBlockService) PushBlock(blockDetails BlockDetails) error {
	err := ebs.ethBlockRepo.PushEthBlock(blockDetails)

	if err != nil {
		log.Printf("PushBlock::err:%v\n", err)
		return err
	}

	return nil
}

func (ebs *ethBlockService) GetCurrentBlockNumber() int {
	return ebs.ethBlockRepo.GetCurrentBlockNumber()
}

func (ebs *ethBlockService) GetLatestBlockNumber() int {
	body := JsonBody{
		Jsonrpc: ebs.jsonRpc,
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		Id:      url_helper.GetRandId(),
	}

	var result BlockNumberResponse

	if err := json_helper.PostRequest(ebs.url, body, &result); err != nil {
		log.Fatalf("Failed to post request: %v", err)
		return -1
	}

	return ebs.HexToInt(result.Result)
}

func (ebs *ethBlockService) GetBlockByNumber(blockNumber int) (BlockResponse, error) {

	hexStr := ebs.IntToHex(blockNumber)
	body := JsonBody{
		Jsonrpc: ebs.jsonRpc,
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{hexStr, false},
		Id:      url_helper.GetRandId(),
	}

	var result BlockResponse

	if err := json_helper.PostRequest(ebs.url, body, &result); err != nil {
		log.Fatalf("Failed to post request: %v", err)
		return result, err
	}

	return result, nil
}

func (ebs *ethBlockService) SetCurrentBlockNumber(blockNumber int) {
	ebs.ethBlockRepo.SetCurrentBlockNumber(blockNumber)
}

func (ebs *ethBlockService) HexToInt(hexStr string) int {

	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}

	intValue, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		log.Printf("SetCurrentBlockNumber::error: %v\n", err)
		return -1
	}

	return int(intValue)
}

func (ebs *ethBlockService) IntToHex(num int) string {
	n := int64(num)
	return fmt.Sprintf("0x%x", n)
}
