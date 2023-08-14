package models

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

type BlockIdentifier struct {
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifierType"`
}

type BlockDetails struct {
	BaseFeePerGas    string   `json:"baseFeePerGas" bson:"baseFeePerGas" redis:"baseFeePerGas"`
	Difficulty       string   `json:"difficulty" bson:"difficulty" redis:"difficulty"`
	ExtraData        string   `json:"extraData" bson:"extraData" redis:"extraData"`
	GasLimit         string   `json:"gasLimit" bson:"gasLimit" redis:"gasLimit"`
	GasUsed          string   `json:"gasUsed" bson:"gasUsed" redis:"gasUsed"`
	Hash             string   `json:"hash" bson:"hash" redis:"hash"`
	LogsBloom        string   `json:"logsBloom" bson:"logsBloom" redis:"logsBloom"`
	Miner            string   `json:"miner" bson:"miner" redis:"miner"`
	MixHash          string   `json:"mixHash" bson:"mixHash" redis:"mixHash"`
	Nonce            string   `json:"nonce" bson:"nonce" redis:"nonce"`
	Number           string   `json:"number" bson:"number" redis:"number"`
	ParentHash       string   `json:"parentHash" bson:"parentHash" redis:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot" bson:"receiptsRoot" redis:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles" bson:"sha3Uncles" redis:"sha3Uncles"`
	Size             string   `json:"size" bson:"size" redis:"size"`
	StateRoot        string   `json:"stateRoot" bson:"stateRoot" redis:"stateRoot"`
	Timestamp        string   `json:"timestamp" bson:"timestamp" redis:"timestamp"`
	TotalDifficulty  string   `json:"totalDifficulty" bson:"totalDifficulty" redis:"totalDifficulty"`
	Transactions     []string `json:"transactions" bson:"transactions" redis:"transactions"`
	TransactionsRoot string   `json:"transactionsRoot" bson:"transactionsRoot" redis:"transactionsRoot"`
	Uncles           []string `json:"uncles" bson:"uncles" redis:"uncles"`
	Withdrawals      []struct {
		Address        string `json:"address" bson:"address" redis:"address"`
		Amount         string `json:"amount" bson:"amount" redis:"amount"`
		Index          string `json:"index" bson:"index" redis:"index"`
		ValidatorIndex string `json:"validatorIndex" bson:"validatorIndex" redis:"validatorIndex"`
	} `json:"withdrawals" bson:"withdrawals" redis:"withdrawals"`
	WithdrawalsRoot string `json:"withdrawalsRoot" bson:"withdrawalsRoot" redis:"withdrawalsRoot"`
}

type BlockResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Result  BlockDetails  `json:"result"`
	Error   ErrorResponse `json:"error,omitempty"`
}
