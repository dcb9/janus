package qtum

import (
	"encoding/json"
	"fmt"
)

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	ID      json.RawMessage `json:"id"`
	Params  json.RawMessage `json:"params"`
}

type JSONRPCRersult struct {
	RawResult json.RawMessage `json:"result"`
	Error     *JSONRPCError   `json:"error"`
	ID        json.RawMessage `json:"id"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *JSONRPCError) Error() string {
	return fmt.Sprintf("[code: %d] %s", err.Code, err.Message)
}
