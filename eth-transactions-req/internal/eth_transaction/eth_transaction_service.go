package eth_transaction

import (
	"eth-transactions-req/pkg/json_helper"
	"eth-transactions-req/pkg/url_helper"
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

	log.Printf("TRANSACTIONS-REQ::GetEthTransaction::hash: %v\n", hash)

	id := url_helper.GetRandId()
	params := []interface{}{hash}
	method := "eth_getTransactionByHash"

	log.Printf("TRANSACTIONS-REQ::GetEthTransaction::id: %v\n", id)
	log.Printf("TRANSACTIONS-REQ::GetEthTransaction::params: %v\n", params)
	log.Printf("TRANSACTIONS-REQ::GetEthTransaction::method: %v\n", method)

	body := JsonBody{
		Jsonrpc: ebs.jsonRpc,
		Method:  method,
		Params:  params,
		Id:      id,
	}

	log.Printf("TRANSACTIONS-REQ::GetEthTransaction::body: %v\n", body)

	var result EthTransactionResponse

	if err := json_helper.PostRequest(ebs.url, body, &result); err != nil {
		log.Fatalf("ERROR::TRANSACTIONS-REQ::GetEthTransaction: %v", err)
		return nil, err
	}

	log.Printf("TRANSACTIONS-REQ::GetEthTransaction::result: %v\n", result)

	et := result.Result

	log.Printf("TRANSACTIONS-REQ::GetEthTransaction::et: %v\n", et)

	ebs.ethTransactionRepo.PushEthTransaction(result.Result)

	return &et, nil
}
