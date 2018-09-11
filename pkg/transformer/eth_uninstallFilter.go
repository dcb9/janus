package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProxyETHUninstallFilter implements ETHProxy
type ProxyETHUninstallFilter struct {
	*qtum.Qtum
	filter *eth.FilterSimulator
}

func (p *ProxyETHUninstallFilter) Method() string {
	return "eth_uninstallFilter"
}

func (p *ProxyETHUninstallFilter) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.UninstallFilterRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHUninstallFilter) request(ethreq *eth.UninstallFilterRequest) (eth.UninstallFilterResponse, error) {
	id, err := hexutil.DecodeUint64(string(*ethreq))
	if err != nil {
		return false, err
	}

	// uninstall
	p.filter.Uninstall(id)

	return true, nil
}
