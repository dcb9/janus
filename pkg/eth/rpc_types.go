package eth

import (
	"encoding/json"
	"errors"
)

type (
	SendTransactionResponse string

	// SendTransactionRequest eth_sendTransaction
	SendTransactionRequest struct {
		From     string  `json:"from"`
		To       string  `json:"to"`
		Gas      *ETHInt `json:"gas"`      // optional
		GasPrice *ETHInt `json:"gasPrice"` // optional
		Value    string  `json:"value"`    // optional
		Data     string  `json:"data"`     // optional
		Nonce    string  `json:"nonce"`    // optional
	}
)

func (r *SendTransactionRequest) UnmarshalJSON(data []byte) error {
	type Request SendTransactionRequest

	var params []Request
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}

	*r = SendTransactionRequest(params[0])

	return nil
}

// see: https://ethereum.stackexchange.com/questions/8384/transfer-an-amount-between-two-ethereum-accounts-using-json-rpc
func (t *SendTransactionRequest) IsSendEther() bool {
	// data must be empty
	return t.Value != "" && t.To != "" && t.From != "" && t.Data == ""
}

func (t *SendTransactionRequest) IsCreateContract() bool {
	return t.To == "" && t.Data != ""
}

func (t *SendTransactionRequest) IsCallContract() bool {
	return t.To != "" && t.Data != ""
}

func (t *SendTransactionRequest) GasHex() string {
	if t.Gas == nil {
		return ""
	}

	return t.Gas.Hex()
}

func (t *SendTransactionRequest) GasPriceHex() string {
	if t.GasPrice == nil {
		return ""
	}
	return t.GasPrice.Hex()
}

// CallResponse
type CallResponse string

// CallRequest eth_call
type CallRequest struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Gas      *ETHInt `json:"gas"`      // optional
	GasPrice *ETHInt `json:"gasPrice"` // optional
	Value    string  `json:"value"`    // optional
	Data     string  `json:"data"`     // optional
}

func (t *CallRequest) GasHex() string {
	if t.Gas == nil {
		return ""
	}
	return t.Gas.Hex()
}

func (t *CallRequest) GasPriceHex() string {
	if t.GasPrice == nil {
		return ""
	}
	return t.GasPrice.Hex()
}

func (t *CallRequest) UnmarshalJSON(data []byte) error {
	var err error
	var params []json.RawMessage
	if err = json.Unmarshal(data, &params); err != nil {
		return err
	}

	if len(params) == 0 {
		return errors.New("params must be set")
	}

	type txCallObject CallRequest
	var obj txCallObject
	if err = json.Unmarshal(params[0], &obj); err != nil {
		return err
	}

	cr := CallRequest(obj)
	*t = cr
	return nil
}

type (
	PersonalUnlockAccountResponse bool
	BlockNumberResponse           string
	NetVersionResponse            string
)

// ========== GetLogs ============= //

type (
	GetLogsRequest struct {
		FromBlock json.RawMessage `json:"fromBlock"`
		ToBlock   json.RawMessage `json:"toBlock"`
		Address   json.RawMessage `json:"address"` // string or []string
		Topics    []string        `json:"topics"`
		Blockhash string          `json:"blockhash"`
	}
	GetLogsResponse []Log
)

func (r *GetLogsRequest) UnmarshalJSON(data []byte) error {
	type Request GetLogsRequest
	var params []Request
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}

	if len(params) == 0 {
		return errors.New("params must be set")
	}

	*r = GetLogsRequest(params[0])

	return nil
}

// ========== GetTransactionByHash ============= //
type (
	GetTransactionByHashRequest  string
	GetTransactionByHashResponse struct {
		Hash             string `json:"hash"`             // DATA, 32 Bytes - hash of the transaction.
		Nonce            string `json:"nonce"`            // QUANTITY - the number of transactions made by the sender prior to this one.
		BlockHash        string `json:"blockHash"`        // DATA, 32 Bytes - hash of the block where this transaction was in. null when its pending.
		BlockNumber      string `json:"blockNumber"`      // QUANTITY - block number where this transaction was in. null when its pending.
		TransactionIndex string `json:"transactionIndex"` // QUANTITY - integer of the transactions index position in the block. null when its pending.
		From             string `json:"from"`             // DATA, 20 Bytes - address of the sender.
		To               string `json:"to"`               // DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
		Value            string `json:"value"`            // QUANTITY - value transferred in Wei.
		GasPrice         string `json:"gasPrice"`         // QUANTITY - gas price provided by the sender in Wei.
		Gas              string `json:"gas"`              // QUANTITY - gas provided by the sender.
		Input            string `json:"input"`            // DATA - the data send along with the transaction.
	}
)

func (r *GetTransactionByHashRequest) UnmarshalJSON(data []byte) error {
	var params []string
	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	if len(params) == 0 {
		return errors.New("params must be set")
	}

	*r = GetTransactionByHashRequest(params[0])
	return nil
}

// ========== GetTransactionReceipt ============= //

type (
	GetTransactionReceiptRequest  string
	GetTransactionReceiptResponse struct {
		TransactionHash   string `json:"transactionHash"`   // DATA, 32 Bytes - hash of the transaction.
		TransactionIndex  string `json:"transactionIndex"`  // QUANTITY - integer of the transactions index position in the block.
		BlockHash         string `json:"blockHash"`         // DATA, 32 Bytes - hash of the block where this transaction was in.
		BlockNumber       string `json:"blockNumber"`       // QUANTITY - block number where this transaction was in.
		From              string `json:"from,omitempty"`    // DATA, 20 Bytes - address of the sender.
		To                string `json:"to,omitempty"`      // DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
		CumulativeGasUsed string `json:"cumulativeGasUsed"` // QUANTITY - The total amount of gas used when this transaction was executed in the block.
		GasUsed           string `json:"gasUsed"`           // QUANTITY - The amount of gas used by this specific transaction alone.
		ContractAddress   string `json:"contractAddress"`   // DATA, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise null.
		Logs              []Log  `json:"logs"`              // Array - Array of log objects, which this transaction generated.
		LogsBloom         string `json:"logsBloom"`         // DATA, 256 Bytes - Bloom filter for light clients to quickly retrieve related logs.
		Root              string `json:"root,omitempty"`    // DATA 32 bytes of post-transaction stateroot (pre Byzantium)
		Status            string `json:"status"`            // QUANTITY either 1 (success) or 0 (failure)
	}

	Log struct {
		Removed          string   `json:"removed,omitempty"` // TAG - true when the log was removed, due to a chain reorganization. false if its a valid log.
		LogIndex         string   `json:"logIndex"`          // QUANTITY - integer of the log index position in the block. null when its pending log.
		TransactionIndex string   `json:"transactionIndex"`  // QUANTITY - integer of the transactions index position log was created from. null when its pending log.
		TransactionHash  string   `json:"transactionHash"`   // DATA, 32 Bytes - hash of the transactions this log was created from. null when its pending log.
		BlockHash        string   `json:"blockHash"`         // DATA, 32 Bytes - hash of the block where this log was in. null when its pending. null when its pending log.
		BlockNumber      string   `json:"blockNumber"`       // QUANTITY - the block number where this log was in. null when its pending. null when its pending log.
		Address          string   `json:"address"`           // DATA, 20 Bytes - address from which this log originated.
		Data             string   `json:"data"`              // DATA - contains one or more 32 Bytes non-indexed arguments of the log.
		Topics           []string `json:"topics"`            // Array of DATA - Array of 0 to 4 32 Bytes DATA of indexed log arguments.
		Type             string   `json:"type,omitempty"`
	}
)

func (r *GetTransactionReceiptRequest) UnmarshalJSON(data []byte) error {
	var params []string
	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	if len(params) == 0 {
		return errors.New("params must be set")
	}

	*r = GetTransactionReceiptRequest(params[0])
	return nil
}

// ========== eth_accounts ============= //
type AccountsResponse []string

// ========== eth_getCode ============= //
type (
	GetCodeRequest struct {
		Address     string
		BlockNumber string
	}
	// the code from the given address.
	GetCodeResponse string
)

func (r *GetCodeRequest) UnmarshalJSON(data []byte) error {
	var params []string
	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	if len(params) == 0 {
		return errors.New("params must be set")
	}

	r.Address = params[0]
	if len(params) > 1 {
		r.BlockNumber = params[1]
	}

	return nil
}

// ========== eth_newBlockFilter ============= //
// a filter id
type NewBlockFilterResponse string

// ========== eth_uninstallFilter ============= //
// the filter id
type UninstallFilterRequest string

// true if the filter was successfully uninstalled, otherwise false.
type UninstallFilterResponse bool

func (r *UninstallFilterRequest) UnmarshalJSON(data []byte) error {
	var params []string
	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	if len(params) == 0 {
		return errors.New("params must be set")
	}

	*r = UninstallFilterRequest(params[0])

	return nil
}

// ========== eth_getFilterChanges ============= //
// the filter id
type GetFilterChangesRequest string

//For filters created with eth_newBlockFilter the return are block hashes (DATA, 32 Bytes),
// e.g. ["0x3454645634534..."].
type GetFilterChangesResponse []string

func (r *GetFilterChangesRequest) UnmarshalJSON(data []byte) error {
	var params []string
	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	if len(params) == 0 {
		return errors.New("params must be set")
	}

	*r = GetFilterChangesRequest(params[0])

	return nil
}

// ========== eth_getTransactionCount ============= //
type GetTransactionCountRequest struct {
	Address     string
	BlockNumber json.RawMessage
}

// integer of the number of transactions send from this address.
type GetTransactionCountResponse string

func (r *GetTransactionCountRequest) UnmarshalJSON(data []byte) error {
	var params []json.RawMessage
	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	if len(params) < 2 {
		return errors.New("params must be set")
	}

	var addr string
	if err = json.Unmarshal(params[0], &addr); err != nil {
		return err
	}

	r.Address = addr
	r.BlockNumber = params[1]

	return nil
}

// ========== eth_getBalance ============= //
type GetBalanceRequest struct {
	Address     string
	BlockNumber json.RawMessage
}

// integer of the current balance in wei.
type GetBalanceResponse string

func (r *GetBalanceRequest) UnmarshalJSON(data []byte) error {
	var params []json.RawMessage
	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	if len(params) < 2 {
		return errors.New("params must be set")
	}

	var addr string
	if err = json.Unmarshal(params[0], &addr); err != nil {
		return err
	}

	r.Address = addr
	r.BlockNumber = params[1]

	return nil
}

// ========== eth_getBlockByNumber ============= //
type GetBlockByNumberRequest struct {
	BlockNumber json.RawMessage
	// If `true` it returns the full transaction objects,
	// if `false` only the hashes of the transactions.
	FullTransactions bool
}

type GetBlockByNumberResponse struct {
	Number           string        `json:"number"`
	Hash             string        `json:"hash"`
	ParentHash       string        `json:"parentHash"`
	Nonce            string        `json:"nonce"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	StateRoot        string        `json:"stateRoot"`
	Miner            string        `json:"miner"`
	Difficulty       string        `json:"difficulty"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             string        `json:"size"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []interface{} `json:"transactions"`
	Uncles           []string      `json:"uncles"`
}

func (r *GetBlockByNumberRequest) UnmarshalJSON(data []byte) error {
	var params []json.RawMessage
	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	if len(params) < 2 {
		return errors.New("params must be set")
	}

	var fullTransactions bool
	if err = json.Unmarshal(params[1], &fullTransactions); err != nil {
		return err
	}

	r.FullTransactions = fullTransactions
	r.BlockNumber = params[0]

	return nil
}
