package transformer

import (
	"encoding/json"
	"errors"

	"github.com/bitly/go-simplejson"
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
)

func (m *Manager) Call(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	var params []json.RawMessage
	if err := unmarshalRequest(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	var tx eth.TransactionCallReq
	if err := unmarshalRequest(params[0], &tx); err != nil {
		return nil, err
	}
	gasLimit, _, err := EthGasToQtum(&tx)
	if err != nil {
		return nil, err
	}

	from := tx.From

	if IsEthHex(from) {
		from, err = m.qtumClient.FromHexAddress(EthHexToQtum(from))
		if err != nil {
			return nil, err
		}
	}

	newParams, err := json.Marshal([]interface{}{
		EthHexToQtum(tx.To),
		EthHexToQtum(tx.Data),
		from,
		gasLimit,
	})
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = qtum.MethodCallcontract

	//Qtum RPC
	// callcontract "address" "data" ( address )
	//
	// Argument:
	//   1. "address"          (string, required) The account address
	//   2. "data"             (string, required) The data hex string
	//   3. address              (string, optional) The sender address hex string
	//   4. gasLimit             (string, optional) The gas limit for executing the contract

	return m.CallcontractResp, nil
}

func (m *Manager) CallcontractResp(result json.RawMessage) (interface{}, error) {
	sj, err := simplejson.NewJson(result)
	if err != nil {
		return nil, err
	}
	output, err := sj.Get("executionResult").Get("output").String()
	if err != nil {
		return nil, err
	}

	return QtumHexToEth(output), nil
}
