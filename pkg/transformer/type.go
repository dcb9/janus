package transformer

import (
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/go-kit/kit/log"
)

type RequestTransformerFunc func(*rpc.JSONRPCRequest) (ResponseTransformerFunc, error)
type ResponseTransformerFunc func(*rpc.JSONRPCResult) error

type context struct {
	logger log.Logger
	req    *rpc.JSONRPCRequest
}

func NopRequestTransformer(_ *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	return func(result *rpc.JSONRPCResult) error {
		return nil
	}, nil
}

/*

func (fn RequestTransformerFunc) Transform(req *rpc.JSONRPCRequest) (*rpc.JSONRPCRequest, error) {
	return fn(req)
}

type RequestTransformer interface {
	Transform(*rpc.JSONRPCRequest) (*rpc.JSONRPCRequest, error)
}

type ResponseTransformerFunc func(*rpc.JSONRPCResult) (*rpc.JSONRPCResult, error)

func (fn ResponseTransformerFunc) Transform(r *rpc.JSONRPCResult) (*rpc.JSONRPCResult, error) {
	return fn(r)
}

type ResponseTransformer interface {
	Transform(*rpc.JSONRPCResult) (*rpc.JSONRPCResult, error)
}
*/
