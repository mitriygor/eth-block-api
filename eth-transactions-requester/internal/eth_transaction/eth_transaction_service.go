package eth_transaction

import (
	"eth-helpers/json_helper"
	"eth-helpers/url_helper"
	"eth-transactions-requester/pkg/logger"
)

type Service interface {
	GetEthTransaction(hash string) (*EthTransaction, error)
}

type ethTransactionService struct {
	ethTransactionRepo Repository
	url                string
	jsonRpc            string
}

func NewEthTransactionService(repo Repository, url string, jsonRpc string) Service {
	return &ethTransactionService{
		ethTransactionRepo: repo,
		url:                url,
		jsonRpc:            jsonRpc,
	}
}

func (ebs *ethTransactionService) GetEthTransaction(hash string) (*EthTransaction, error) {
	id := url_helper.GetRandId()
	params := []interface{}{hash}
	method := "eth_getTransactionByHash"

	body := JsonBody{
		Jsonrpc: ebs.jsonRpc,
		Method:  method,
		Params:  params,
		Id:      id,
	}

	var result EthTransactionResponse

	if err := json_helper.PostRequest(ebs.url, body, &result); err != nil {
		logger.Error("eth-transactions-requester::ERROR::GetEthTransaction::err", "error", err)
		return nil, err
	}

	et := result.Result

	err := ebs.ethTransactionRepo.PushEthTransaction(result.Result)
	if err != nil {
		logger.Error("eth-transactions-requester::ERROR::GetEthTransaction::err", "error", err)
		return nil, err
	}

	return &et, nil
}
