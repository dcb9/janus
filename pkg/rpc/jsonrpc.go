package rpc

import (
	"encoding/json"
	"fmt"
)

const (
	ErrInvalid          = 150
	ErrUnknownOperation = 151
)

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	ID      json.RawMessage `json:"id"`
	Params  json.RawMessage `json:"params"`
}

type JSONRPCResult struct {
	JSONRPC   string          `json:"jsonrpc"`
	RawResult json.RawMessage `json:"result"`
	Error     *JSONRPCError   `json:"error,omitempty"`
	ID        json.RawMessage `json:"id"`
}

type SuccessJSONRPCResult struct {
	JSONRPC   string          `json:"jsonrpc"`
	RawResult json.RawMessage `json:"result"`
	ID        json.RawMessage `json:"id"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *JSONRPCError) Error() string {
	return fmt.Sprintf("[code: %d] %s", err.Code, err.Message)
}
