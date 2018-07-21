package transformer

import (
	"encoding/json"

	"github.com/dcb9/janus/pkg/rpc"
	"github.com/pkg/errors"
)

type (
	RequestTransformerFunc  func(*rpc.JSONRPCRequest) (ResponseTransformerFunc, error)
	ResponseTransformerFunc func(result json.RawMessage) (newResult interface{}, err error)
)

var (
	UnmarshalRequestErr = errors.New("unmarshal request error")
)
