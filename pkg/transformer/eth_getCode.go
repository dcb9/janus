package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
)

// ProxyETHGetCode implements ETHProxy
type ProxyETHGetCode struct {
	*qtum.Qtum
}

func (p *ProxyETHGetCode) Method() string {
	return "eth_getCode"
}

func (p *ProxyETHGetCode) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetCodeRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHGetCode) request(ethreq *eth.GetCodeRequest) (eth.GetCodeResponse, error) {
	qtumreq := qtum.GetAccountInfoRequest(utils.RemoveHexPrefix(ethreq.Address))

	qtumresp, err := p.GetAccountInfo(&qtumreq)
	if err != nil {
		return "", err
	}

	// qtum res -> eth res
	return eth.GetCodeResponse(utils.AddHexPrefix(qtumresp.Code)), nil
}
