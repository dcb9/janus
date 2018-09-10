package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProxyETHBlockNumber implements ETHProxy
type ProxyETHBlockNumber struct {
	*qtum.Qtum
}

func (p *ProxyETHBlockNumber) Method() string {
	return "eth_blockNumber"
}

func (p *ProxyETHBlockNumber) Request(_ *eth.JSONRPCRequest) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHBlockNumber) request() (*eth.BlockNumberResponse, error) {
	qtumresp, err := p.Qtum.GetBlockCount()
	if err != nil {
		return nil, err
	}

	// qtum res -> eth res
	return p.ToResponse(qtumresp), nil
}

func (p *ProxyETHBlockNumber) ToResponse(qtumresp *qtum.GetBlockCountResponse) *eth.BlockNumberResponse {
	hexVal := hexutil.EncodeBig(qtumresp.Int)
	ethresp := eth.BlockNumberResponse(hexVal)
	return &ethresp
}
