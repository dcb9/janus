package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHEstimateGas struct {
	*qtum.Qtum
}

func (p *ProxyETHEstimateGas) Method() string {
	return "eth_estimateGas"
}

func (p *ProxyETHEstimateGas) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHEstimateGas) request() (*eth.EstimateGasResponse, error) {
	gas := eth.EstimateGasResponse("0x500000")
	return &gas, nil
}
