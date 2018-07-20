package transformer

import (
	"encoding/json"
	"errors"

	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/go-kit/kit/log"
)

func (m *Manager) Call(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
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

	l := log.WithPrefix(m.logger, "method", req.Method)

	return func(result *rpc.JSONRPCResult) error {
		return m.CallcontractResp(context{
			logger: l,
			req:    req,
		}, result)
	}, nil
}

func (m *Manager) CallcontractResp(c context, result *rpc.JSONRPCResult) error {
	if result.Error != nil {
		return result.Error
	}

	if result.RawResult != nil {
		sj, err := simplejson.NewJson(result.RawResult)
		if err != nil {
			return err
		}
		output, err := sj.Get("executionResult").Get("output").Bytes()
		if err != nil {
			return err
		}

		outputStr := fmt.Sprintf(`"0x%s"`, output)
		result.RawResult = []byte(outputStr)
		result.JSONRPC = "2.0"
		return nil
	}

	return errors.New("result.RawResult must not be nil")
}
