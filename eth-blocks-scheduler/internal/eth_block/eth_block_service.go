package eth_block

import (
	"eth-blocks-scheduler/pkg/json_helper"
	"eth-blocks-scheduler/pkg/url_helper"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	PushLatestBlocks()
}

type ethBlockService struct {
	ethBlockRepo      Repository
	url               string
	jsonRpc           string
	cacheSize         int
	requesterInterval int
}

func NewEthBlockService(repo Repository, url string, jsonRpc string, cacheSize int, requesterInterval int) Service {
	return &ethBlockService{
		ethBlockRepo:      repo,
		url:               url,
		jsonRpc:           jsonRpc,
		cacheSize:         cacheSize,
		requesterInterval: requesterInterval,
	}
}

func (ebs *ethBlockService) PushLatestBlocks() {
	latest := ebs.getLatestBlockNumber()

	log.Printf("eth-blocks-scheduler::PushLatestBlocks::latest::%v\n", latest)
	diff := latest - ebs.cacheSize
	log.Printf("eth-blocks-scheduler::PushLatestBlocks::diff::%v\n", diff)

	if diff < 0 {
		log.Println("eth-blocks-scheduler::Issue to get blocks numbers")
		return
	}

	for i := diff; i < latest; i++ {
		blockResponse, err := ebs.getBlockByNumber(i)

		if err != nil {
			log.Printf("eth-blocks-scheduler::PushLatestBlocks::blockResponse::error:%v\n", err)
			continue
		}

		log.Printf("eth-blocks-scheduler::PushLatestBlocks::blockResponse::%v\n", blockResponse.Result)

		ebs.pushBlock(blockResponse.Result)

		time.Sleep(time.Duration(ebs.requesterInterval) * time.Second)
	}
}

func (ebs *ethBlockService) pushBlock(blockDetails BlockDetails) {

	log.Printf("eth-blocks-scheduler::pushBlock::blockDetails::%v\n", blockDetails)

	ebs.ethBlockRepo.PushBlocksForRecording(blockDetails)
	ebs.ethBlockRepo.PushBlocksForRedis(blockDetails)
	ebs.ethBlockRepo.PushTransactionsForScheduling(blockDetails.Transactions)

}

func (ebs *ethBlockService) getLatestBlockNumber() int {
	body := JsonBody{
		Jsonrpc: ebs.jsonRpc,
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		Id:      url_helper.GetRandId(),
	}

	var result BlockNumberResponse

	if err := json_helper.PostRequest(ebs.url, body, &result); err != nil {
		log.Printf("eth-blocks-scheduler::Failed to post request: %v\n", err)
		return -1
	}

	return ebs.HexToInt(result.Result)
}

func (ebs *ethBlockService) getBlockByNumber(blockNumber int) (BlockResponse, error) {

	hexStr := ebs.IntToHex(blockNumber)
	body := JsonBody{
		Jsonrpc: ebs.jsonRpc,
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{hexStr, false},
		Id:      url_helper.GetRandId(),
	}

	var result BlockResponse

	if err := json_helper.PostRequest(ebs.url, body, &result); err != nil {
		log.Printf("eth-blocks-scheduler::Failed to post request: %v\n", err)
		return result, err
	}

	return result, nil
}

func (ebs *ethBlockService) HexToInt(hexStr string) int {

	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}

	intValue, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		log.Printf("eth-blocks-scheduler::ERROR::HexToInt::error: %v\n", err)
		return -1
	}

	return int(intValue)
}

func (ebs *ethBlockService) IntToHex(num int) string {
	n := int64(num)
	return fmt.Sprintf("0x%x", n)
}
