package transformer

import (
	"encoding/json"
	"errors"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/rpc"
)

func transformCall(req *rpc.JSONRPCRequest) (*rpc.JSONRPCRequest, error) {
	var params []json.RawMessage
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	var tx eth.TransactionCallReq
	if err := json.Unmarshal(params[0], &tx); err != nil {
		return nil, err
	}
	gasLimit, _, err := EthGasToQtum(&tx)
	if err != nil {
		return nil, err
	}

	newParams, err := json.Marshal([]interface{}{
		EthHexToQtum(tx.To),
		EthHexToQtum(tx.Data),
		tx.From,
		gasLimit,
	})
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = "callcontract"

	//Qtum RPC
	// callcontract "address" "data" ( address )
	//
	// Argument:
	//   1. "address"          (string, required) The account address
	//   2. "data"             (string, required) The data hex string
	//   3. address              (string, optional) The sender address hex string
	//   4. gasLimit             (string, optional) The gas limit for executing the contract

	return req, nil
}
