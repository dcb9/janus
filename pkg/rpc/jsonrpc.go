package rpc

import (
	"encoding/json"
	"fmt"
)

// FIXME: move to qtum package
const (
	ErrInvalid          = 150
	ErrUnknownOperation = 151
)

// FIXME: rename package to jsonrpc
// FIXME: remove JSONRPC prefix for JSONRPCRequest, JSONRPCResult
// FIXME: this package seems kinda pointless

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

// FIXME: move this to qtum package
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *JSONRPCError) Error() string {
	return fmt.Sprintf("[code: %d] %s", err.Code, err.Message)
}
