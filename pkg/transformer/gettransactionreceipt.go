package transformer

import (
	"encoding/json"
	"errors"

	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/go-kit/kit/log"
)

func (m *Manager) GetTransactionReceipt(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	var params []json.RawMessage
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	newParams, err := json.Marshal([]interface{}{
		EthHexToQtum(string(params[0])),
	})
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = qtum.MethodGettransactionreceipt

	//Qtum RPC
	//gettransactionreceipt "hash"
	//  requires -logevents to be enabled
	//  Argument:
	//  1. "hash"          (string, required) The transaction hash

	l := log.WithPrefix(m.logger, "method", req.Method)
	return func(result *rpc.JSONRPCResult) error {
		return m.GettransactionreceiptResp(context{
			logger: l,
			req:    req,
		}, result)
	}, nil
}

func (m *Manager) GettransactionreceiptResp(c context, result *rpc.JSONRPCResult) error {
	return nil
}
