package transformer

import (
	"encoding/json"
	"errors"

	"github.com/bitly/go-simplejson"
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
)

// method: eth_call
//
// eth request -> qtum request
// call qtum rpc
// qtum reponse -> eth response

func (m *Manager) proxyRPC(req *rpc.JSONRPCRequest) (interface{}, error) {
	switch req.Method {
	case "eth_call":
		var req eth.CallRequest
		err := unmarshalRequest(req.Params, &req)
		if err != nil {
			return nil, err
		}
		return m.proxyETHCall(req)
	default:
		return errors.New("unsupported method")
	}
}

type Proxy interface {
	request(ethreq interface{}) (interface{}, error)
}

// FIXME: rename file to eth_call.go
// proxy eth_call

type ProxyETHCall struct{}

func (p *ProxyETHCall) request(rawreq *rpc.JSONRPCRequest) (interface{}, error) {
	var req eth.CallRequest
	err := unmarshalRequest(rawreq.Params, &req)
	if err != nil {
		return nil, err
	}
	return p.proxyInternal(req)
}

func (p *ProxyETHCall) requestInternal(ethreq *eth.CallRequest) (*eth.CallResponse, error) {
	// eth req -> qtum req
	qtumreq, err = p.ToRequest(ethreq)

	var qtumres qtum.CallContract
	err = p.rpc.Request(qtumreq, &qtumres)

	// qtum res -> eth res
	ethres, err = p.ToResponse(ethreq)

	return ethres, err
}

func (p *ProxyETHCall) ToRequest(ethreq *eth.CallRequest) (*qtum.CallContract, error) {
}

func (p *ProxyETHCall) ToResponse(res *qtum.CallResponse) (*eth.CallResponse, error) {
}

func (m *Manager) Call(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	var params []json.RawMessage
	if err := unmarshalRequest(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	// FIXME: rename to eth.CallRequest
	var tx eth.TransactionCallReq
	if err := unmarshalRequest(params[0], &tx); err != nil {
		return nil, err
	}

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

	return AddHexPrefix(output), nil
}
