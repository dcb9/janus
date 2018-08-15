package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// ProxyETHGetTransactionCount implements ETHProxy
type ProxyETHGetTransactionCount struct {
	*qtum.Qtum
}

func (p *ProxyETHGetTransactionCount) Method() string {
	return "eth_getTransactionCount"
}

func (p *ProxyETHGetTransactionCount) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetTransactionCountRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHGetTransactionCount) request(ethreq *eth.GetTransactionCountRequest) (eth.GetTransactionCountResponse, error) {
	// eth req -> qtum req
	var err error
	addr := ethreq.Address
	if utils.IsEthHexAddress(addr) {
		addr, err = p.FromHexAddress(addr)
		if err != nil {
			return "", errors.Wrap(err, "FromHexAddress")
		}
	}

	summary, err := p.Insight.GetAddressSummary(addr)
	if err != nil {
		return "", errors.Wrap(err, "Insight#GetAddressSummary")
	}

	count := len(summary.Transactions)
	countHex := hexutil.EncodeUint64(uint64(count))

	return eth.GetTransactionCountResponse(countHex), nil
}

func (p *ProxyETHGetTransactionCount) ToRequest(ethreq *eth.CallRequest) (*qtum.CallContractRequest, error) {
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

func (p *ProxyETHGetTransactionCount) ToResponse(ethresp *qtum.CallContractResponse) *eth.CallResponse {
	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
	qtumresp := eth.CallResponse(data)
	return &qtumresp
}
