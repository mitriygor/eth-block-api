package eth_transaction

import (
	"eth-helpers/json_helper"
	"eth-helpers/url_helper"
	"log"
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

	log.Printf("eth-transactions-requester::GetEthTransaction::hash: %v\n", hash)

	id := url_helper.GetRandId()
	params := []interface{}{hash}
	method := "eth_getTransactionByHash"

	log.Printf("eth-transactions-requester::GetEthTransaction::id: %v\n", id)
	log.Printf("eth-transactions-requester::GetEthTransaction::params: %v\n", params)
	log.Printf("eth-transactions-requester::GetEthTransaction::method: %v\n", method)

	body := JsonBody{
		Jsonrpc: ebs.jsonRpc,
		Method:  method,
		Params:  params,
		Id:      id,
	}

	log.Printf("eth-transactions-requester::GetEthTransaction::body: %v\n", body)

	var result EthTransactionResponse

	if err := json_helper.PostRequest(ebs.url, body, &result); err != nil {
		log.Fatalf("eth-transactions-requester::ERROR::GetEthTransaction: %v\n", err)
		return nil, err
	}

	log.Printf("eth-transactions-requester::GetEthTransaction::result: %v\n", result)

	et := result.Result

	log.Printf("eth-transactions-requester::GetEthTransaction::et: %v\n", et)

	ebs.ethTransactionRepo.PushEthTransaction(result.Result)

	return &et, nil
}
