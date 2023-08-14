package eth_block

import (
	"eth-req/pkg/json_helper"
	"eth-req/pkg/url_helper"
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
		log.Fatalf("ERROR::REQ::GetEthBlock: %v", err)
		return nil, err
	}

	ebs.ethBlockRepo.PushEthBlock(result.Result)

	return &result, nil
}
