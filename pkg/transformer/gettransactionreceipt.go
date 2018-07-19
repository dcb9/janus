package transformer

import (
	"encoding/json"
	"errors"

	"github.com/dcb9/janus/pkg/rpc"
)

func transformTransactionReceipt(req *rpc.JSONRPCRequest) (*rpc.JSONRPCRequest, error) {
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
	req.Method = "gettransactionreceipt"

	//Qtum RPC
	//gettransactionreceipt "hash"
	//  requires -logevents to be enabled
	//  Argument:
	//  1. "hash"          (string, required) The transaction hash

	return req, nil
}
