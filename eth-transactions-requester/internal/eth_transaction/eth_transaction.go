package eth_transaction

type JsonBody struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type EthTransactionResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      int            `json:"id"`
	Result  EthTransaction `json:"result"`
	Error   ErrorResponse  `json:"error,omitempty"`
}

type EthTransactionIdentifier struct {
	Hash string `json:"hash"`
}

type EthTransaction struct {
	AccessList           []AccessListEntry `json:"accessList" bson:"accessList"`
	TransactionHash      string            `json:"transactionHash" bson:"transactionHash"`
	TransactionNumber    string            `json:"transactionNumber" bson:"transactionNumber"`
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
