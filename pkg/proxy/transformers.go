package proxy

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/dcb9/janus/pkg/qtum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

func transformSendTransaction(req *qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error) {
	params := make([]map[string]string, 0)
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be empty")
	}

	param := params[0]
	if param["to"] == "" {
		return createcontract(req, param)
	}

	return sendtocontract(param)
}

type transformerFunc func(*qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error)

func (fn transformerFunc) transform(req *qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error) {
	return fn(req)
}

type transformer interface {
	transform(*qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error)
}

func createcontract(req *qtum.JSONRPCRequest, param map[string]string) (*qtum.JSONRPCRequest, error) {
	if v, ok := param["value"]; ok {
		if v != "" && v != "0x0" {
			return nil, &qtum.JSONRPCError{
				Code:    ErrInvalid,
				Message: "value must be empty",
			}
		}
	}

	//  eth rpc
	//  params: [{
	//    "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
	//    "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
	//    "gas": "0x76c0", // 30400
	//    "gasPrice": "0x9184e72a000", // 10000000000000
	//    "value": "0x9184e72a", // 2441406250
	//    "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
	//  }]

	// qtum rpc
	//  createcontract "bytecode" (gaslimit gasprice "senderaddress" broadcast)
	//  Create a contract with bytcode.
	//  Arguments:
	//  1. "bytecode"  (string, required) contract bytcode.
	//  2. gasLimit  (numeric or string, optional) gasLimit, default: 2500000, max: 40000000
	//  3. gasPrice  (numeric or string, optional) gasPrice QTUM price per gas unit, default: 0.0000004, min:0.0000004
	//  4. "senderaddress" (string, optional) The quantum address that will be used to create the contract.
	//  5. "broadcast" (bool, optional, default=true) Whether to broadcast the transaction or not.
	//  6. "changeToSender" (bool, optional, default=true) Return the change to the sender.

	var err error
	gasLimit := big.NewInt(2500000)
	if v, ok := param["gas"]; ok {
		gasLimit, err = hexutil.DecodeBig(v)
		if err != nil {
			return nil, errors.Wrap(err, "decode gas")
		}
	}
	gasPrice := "0.0000004"
	if _, ok := param["gasPrice"]; ok {
		// fixme parse gas price
	}

	_, _ = gasLimit, gasPrice
	//params := []interface{}{param["data"], gasLimit, gasPrice}
	params := []interface{}{param["data"], gasLimit, gasPrice}

	if f, ok := param["from"]; ok {
		sender := f
		if strings.HasPrefix(f, "0x") {
			// todo convert hexaddress
		}

		params = append(params, sender)
	}

	newParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = "createcontract"
	return req, nil
}

func sendtocontract(param map[string]string) (*qtum.JSONRPCRequest, error) {
	return nil, &qtum.JSONRPCError{
		Code:    ErrInvalid,
		Message: "unsupport",
	}
}
