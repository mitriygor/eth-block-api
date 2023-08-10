package eth_block

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

type BlockNumberResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Result  string        `json:"result"`
	Error   ErrorResponse `json:"error,omitempty"`
}

type BlockDetails struct {
	BaseFeePerGas    string   `json:"baseFeePerGas" bson:"baseFeePerGas"`
	Difficulty       string   `json:"difficulty" bson:"difficulty"`
	ExtraData        string   `json:"extraData" bson:"extraData"`
	GasLimit         string   `json:"gasLimit" bson:"gasLimit"`
	GasUsed          string   `json:"gasUsed" bson:"gasUsed"`
	Hash             string   `json:"hash" bson:"hash"`
	LogsBloom        string   `json:"logsBloom" bson:"logsBloom"`
	Miner            string   `json:"miner" bson:"miner"`
	MixHash          string   `json:"mixHash" bson:"mixHash"`
	Nonce            string   `json:"nonce" bson:"nonce"`
	Number           string   `json:"number" bson:"number"`
	ParentHash       string   `json:"parentHash" bson:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot" bson:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles" bson:"sha3Uncles"`
	Size             string   `json:"size" bson:"size"`
	StateRoot        string   `json:"stateRoot" bson:"stateRoot"`
	Timestamp        string   `json:"timestamp" bson:"timestamp"`
	TotalDifficulty  string   `json:"totalDifficulty" bson:"totalDifficulty"`
	Transactions     []string `json:"transactions" bson:"transactions"`
	TransactionsRoot string   `json:"transactionsRoot" bson:"transactionsRoot"`
	Uncles           []string `json:"uncles" bson:"uncles"`
	Withdrawals      []struct {
		Address        string `json:"address" bson:"address"`
		Amount         string `json:"amount" bson:"amount"`
		Index          string `json:"index" bson:"index"`
		ValidatorIndex string `json:"validatorIndex" bson:"validatorIndex"`
	} `json:"withdrawals" bson:"withdrawals"`
	WithdrawalsRoot string `json:"withdrawalsRoot" bson:"withdrawalsRoot"`
}

type BlockResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Result  BlockDetails  `json:"result"`
	Error   ErrorResponse `json:"error,omitempty"`
}
