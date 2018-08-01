package eth

import (
	"encoding/json"
	"fmt"
)

const (
	RPCVersion = "2.0"
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
	return fmt.Sprintf("eth [code: %d] %s", err.Code, err.Message)
}

func NewJSONRPCResult(id json.RawMessage, res interface{}) (*JSONRPCResult, error) {
	rawResult, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return &JSONRPCResult{
		JSONRPC:   RPCVersion,
		ID:        id,
		RawResult: rawResult,
	}, nil
}
