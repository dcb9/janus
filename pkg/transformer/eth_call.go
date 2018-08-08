package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
)

// ProxyETHCall implements ETHProxy
type ProxyETHCall struct {
	*qtum.Qtum
}

func (p *ProxyETHCall) Method() string {
	return "eth_call"
}

func (p *ProxyETHCall) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHCall) request(ethreq *eth.CallRequest) (*eth.CallResponse, error) {
	// eth req -> qtum req
	params, err := p.ToRequest(ethreq)
	if err != nil {
		return nil, err
	}

	var qtumresp *qtum.CallContractResponse
	if err = p.Qtum.Request(qtum.MethodCallContract, params, &qtumresp); err != nil {
		return nil, err
	}

	// qtum res -> eth res
	return p.ToResponse(qtumresp), nil
}

func (p *ProxyETHCall) ToRequest(ethreq *eth.CallRequest) (*qtum.CallContractRequest, error) {
	gasLimit, _, err := EthGasToQtum(ethreq)
	if err != nil {
		return nil, err
	}

	from := ethreq.From

	if utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return nil, err
		}
	}

	return &qtum.CallContractRequest{
		To:       ethreq.To,
		From:     ethreq.From,
		Data:     ethreq.Data,
		GasLimit: gasLimit,
	}, nil
}

func (p *ProxyETHCall) ToResponse(ethresp *qtum.CallContractResponse) *eth.CallResponse {
	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
	qtumresp := eth.CallResponse(data)
	return &qtumresp
}
