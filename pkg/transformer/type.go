package transformer

import (
	"errors"

	"github.com/dcb9/janus/pkg/eth"
)

var UnmarshalRequestErr = errors.New("Input is invalid")

type Option func(*Transformer) error

type ETHProxy interface {
	Request(*eth.JSONRPCRequest) (interface{}, error)
	Method() string
}
