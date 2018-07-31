package qtum

import (
	"encoding/json"
	"fmt"
)

const (
	RPCVersion = "1.0"
)

const (
	MethodGetHexAddress         = "gethexaddress"
	MethodFromHexAddress        = "fromhexaddress"
	MethodSendToContract        = "sendtocontract"
	MethodGetTransactionReceipt = "gettransactionreceipt"
	MethodGetTransaction        = "gettransaction"
	MethodCreateContract        = "createcontract"
	MethodSendToAddress         = "sendtoaddress"
	MethodCallContract          = "callcontract"
	MethodDecodeRawTransaction  = "decoderawtransaction"
	MethodGetBlockCount         = "getblockcount"
	MethodGetBlockChainInfo     = "getblockchaininfo"
	MethodSearchLogs            = "searchlogs"
	MethodWaitForLogs           = "waitforlogs"
)

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	ID      json.RawMessage `json:"id"`
	Params  json.RawMessage `json:"params"`
}

type JSONRPCResult struct {
	JSONRPC   string          `json:"jsonrpc"`
	RawResult json.RawMessage `json:"result,omitempty"`
	Error     *JSONRPCError   `json:"error,omitempty"`
	ID        json.RawMessage `json:"id"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *JSONRPCError) Error() string {
	return fmt.Sprintf("qtum [code: %d] %s", err.Code, err.Message)
}

type SuccessJSONRPCResult struct {
	JSONRPC   string          `json:"jsonrpc"`
	RawResult json.RawMessage `json:"result"`
	ID        json.RawMessage `json:"id"`
}
