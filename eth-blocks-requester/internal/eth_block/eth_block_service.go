package eth_block

import (
	"eth-helpers/json_helper"
	"eth-helpers/url_helper"
	"log"
)

type Service interface {
	GetEthBlock(identifier string, identifierType string) (*BlockResponse, error)
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

func (ebs *ethBlockService) GetEthBlock(identifier string, identifierType string) (*BlockResponse, error) {

	id := url_helper.GetRandId()
	params := []interface{}{identifier, false}
	method := "eth_getBlockByNumber"

	if identifierType == "hash" {
		method = "eth_getBlockByHash"
	}

	body := JsonBody{
		Jsonrpc: ebs.jsonRpc,
		Method:  method,
		Params:  params,
		Id:      id,
	}

	var result BlockResponse

	if err := json_helper.PostRequest(ebs.url, body, &result); err != nil {
		log.Fatalf("eth-blocks-requester::ERROR::GetEthBlock: %v\n", err)
		return nil, err
	}

	ebs.ethBlockRepo.PushEthBlock(result.Result)

	return &result, nil
}
