package eth_block

import (
	"eth-blocks-scheduler/pkg/logger"
	"eth-helpers/json_helper"
	"eth-helpers/url_helper"
	"fmt"
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
	diff := latest - ebs.cacheSize

	if diff < 0 {
		logger.Error("eth-blocks-scheduler::Issue to get blocks numbers", "latest", latest, "cacheSize", ebs.cacheSize)
		return
	}

	for i := diff; i < latest; i++ {
		blockResponse, err := ebs.getBlockByNumber(i)

		if err != nil {
			logger.Error("eth-blocks-scheduler::PushLatestBlocks::blockResponse::error", "error", err)
			continue
		}

		ebs.pushBlock(blockResponse.Result)

		time.Sleep(time.Duration(ebs.requesterInterval) * time.Second)
	}
}

func (ebs *ethBlockService) pushBlock(blockDetails BlockDetails) {
	err := ebs.ethBlockRepo.PushBlocksForRecording(blockDetails)
	if err != nil {
		logger.Error("eth-blocks-scheduler::pushBlock::PushBlocksForRecording::error", "error", err)
		return
	}
	err = ebs.ethBlockRepo.PushBlocksForRedis(blockDetails)
	if err != nil {
		logger.Error("eth-blocks-scheduler::pushBlock::PushBlocksForRedis::error", "error", err)
		return
	}
	err = ebs.ethBlockRepo.PushTransactionsForScheduling(blockDetails.Transactions)
	if err != nil {
		logger.Error("eth-blocks-scheduler::pushBlock::PushTransactionsForScheduling::error", "error", err)
		return
	}
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
		logger.Error("eth-blocks-scheduler::Failed to post request", "error", err)
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
		logger.Error("eth-blocks-scheduler::Failed to post request", "error", err)
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
		logger.Error("eth-blocks-scheduler::ERROR::HexToInt::error", "error", err)
		return -1
	}

	return int(intValue)
}

func (ebs *ethBlockService) IntToHex(num int) string {
	n := int64(num)
	return fmt.Sprintf("0x%x", n)
}
