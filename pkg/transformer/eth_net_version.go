package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
)

// ProxyETHNetVersion implements ETHProxy
type ProxyETHNetVersion struct {
	*qtum.Qtum
}

func (p *ProxyETHNetVersion) Method() string {
	return "net_version"
}

func (p *ProxyETHNetVersion) Request(_ *eth.JSONRPCRequest) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHNetVersion) request() (*eth.NetVersionResponse, error) {
	var qtumresp *qtum.GetBlockChainInfoResponse
	if err := p.Qtum.Request(qtum.MethodGetBlockChainInfo, nil, &qtumresp); err != nil {
		return nil, err
	}

	resp := eth.NetVersionResponse(qtumresp.Chain)
	return &resp, nil
}
