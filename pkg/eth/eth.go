package eth

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/dcb9/janus/pkg/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func NewJSONRPCResult(id, rawResult json.RawMessage, err *rpc.JSONRPCError) *rpc.JSONRPCResult {
	return &rpc.JSONRPCResult{
		JSONRPC:   "2.0",
		ID:        id,
		RawResult: rawResult,
		Error:     err,
	}
}

// eth_sendTransaction
type TransactionReq struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Gas      *EthInt `json:"gas"`      // optional
	GasPrice *EthInt `json:"gasPrice"` // optional
	Value    string  `json:"value"`    // optional
	Data     string  `json:"data"`     // optional
	Nonce    string  `json:"nonce"`    // optional
}

// see: https://ethereum.stackexchange.com/questions/8384/transfer-an-amount-between-two-ethereum-accounts-using-json-rpc
func (t *TransactionReq) IsSendEther() bool {
	// data must be empty
	return t.Value != "" && t.To != "" && t.From != "" && t.Data == ""
}

func (t *TransactionReq) IsCreateContract() bool {
	return t.To == "" && t.Data != ""
}

func (t *TransactionReq) IsCallContract() bool {
	return t.To != "" && t.Data != ""
}

// FIXME: GetGas -> GasHex
func (t *TransactionReq) GetGas() string {
	return t.Gas.Hex()
}

// FIXME: GetGasPrice -> GasPriceHex
func (t *TransactionReq) GetGasPrice() string {
	return t.GasPrice.Hex()
}

// eth_call
type TransactionCallReq struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Gas      *EthInt `json:"gas"`      // optional
	GasPrice *EthInt `json:"gasPrice"` // optional
	Value    string  `json:"value"`    // optional
	Data     string  `json:"data"`     // optional
}

func (t *TransactionCallReq) GetGas() string {
	return t.Gas.Hex()
}

func (t *TransactionCallReq) GetGasPrice() string {
	return t.GasPrice.Hex()
}

type (
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

	TransactionReceipt struct {
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

	TransactionResponse struct {
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

	GetLogsFilter struct {
		FromBlock json.RawMessage `json:"fromBlock"`
		ToBlock   json.RawMessage `json:"toBlock"`
		Address   json.RawMessage `json:"address"` // string or []string
		Topics    []string        `json:"topics"`
		Blockhash string          `json:"blockhash"`
	}
)

// FIXME: ETHInt
type EthInt big.Int

func (i *EthInt) Hex() string {
	return hexutil.EncodeBig(i.ToBigInt())
}

func (i *EthInt) ToBigInt() *big.Int {
	v := *i
	vv := big.Int(v)
	return &vv
}

func (i *EthInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ToBigInt())
}

// FIXME: extract parsers into constructors:
// ETHIntFromNumber
// ETHIntFromIntger

// UnmarshalJSON needs to be able to parse ETHInt from both hex string or number
func (i *EthInt) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("data must not be empty")
	}

	if data[0] != '"' && data[len(data)-1] != '"' {
		var v *big.Int
		if err = json.Unmarshal(data, &v); err != nil {
			return err
		}
		vv := *v
		*i = EthInt(vv)
		return
	}

	// hex
	var val string
	if err = json.Unmarshal(data, &val); err != nil {
		return err
	}
	if !strings.HasPrefix(val, "0x") {
		val = "0x" + val
	}

	v, err := hexutil.DecodeBig(val)
	if err != nil {
		return err
	}
	vv := *v
	*i = EthInt(vv)
	return err
}
