package models

type EthTransaction struct {
	AccessList           []AccessListEntry `json:"accessList" bson:"accessList"`
	BlockHash            string            `json:"blockHash" bson:"blockHash"`
	BlockNumber          string            `json:"blockNumber" bson:"blockNumber"`
	ChainId              string            `json:"chainId" bson:"chainId"`
	From                 string            `json:"from" bson:"from"`
	Gas                  string            `json:"gas" bson:"gas"`
	GasPrice             string            `json:"gasPrice" bson:"gasPrice"`
	Hash                 string            `json:"hash" bson:"hash"`
	Input                string            `json:"input" bson:"input"`
	MaxFeePerGas         string            `json:"maxFeePerGas" bson:"maxFeePerGas"`
	MaxPriorityFeePerGas string            `json:"maxPriorityFeePerGas" bson:"maxPriorityFeePerGas"`
	Nonce                string            `json:"nonce" bson:"nonce"`
	R                    string            `json:"r" bson:"r"`
	S                    string            `json:"s" bson:"s"`
	To                   string            `json:"to" bson:"to"`
	TransactionIndex     string            `json:"transactionIndex" bson:"transactionIndex"`
	Type                 string            `json:"type" bson:"type"`
	V                    string            `json:"v" bson:"v"`
	Value                string            `json:"value" bson:"value"`
}

type AccessListEntry struct {
	Address     string   `json:"address" bson:"address"`
	StorageKeys []string `json:"storageKeys" bson:"storageKeys"`
}

type EthTransactionResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      int            `json:"id"`
	Result  EthTransaction `json:"result"`
	Error   ErrorResponse  `json:"error,omitempty"`
}
