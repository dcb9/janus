package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHEstimateGas struct {
	*ProxyETHCall
}

func (p *ProxyETHEstimateGas) Method() string {
	return "eth_estimateGas"
}

func (p *ProxyETHEstimateGas) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var ethreq eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &ethreq); err != nil {
		return nil, err
	}

	// eth req -> qtum req
	qtumreq, err := p.ToRequest(&ethreq)
	if err != nil {
		return nil, err
	}

	qtumresp, err := p.CallContract(qtumreq)
	if err != nil {
		return nil, err
	}

	return p.toResp(qtumresp)
}

func (p *ProxyETHEstimateGas) toResp(qtumresp *qtum.CallContractResponse) (*eth.EstimateGasResponse, error) {
	gas := eth.EstimateGasResponse(hexutil.EncodeUint64(uint64(qtumresp.ExecutionResult.GasUsed)))
	return &gas, nil
}
