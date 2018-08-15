package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// ProxyETHGetBalance implements ETHProxy
type ProxyETHGetBalance struct {
	*qtum.Qtum
}

func (p *ProxyETHGetBalance) Method() string {
	return "eth_getBalance"
}

func (p *ProxyETHGetBalance) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetBalanceRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHGetBalance) request(ethreq *eth.GetBalanceRequest) (eth.GetBalanceResponse, error) {
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

	balanceInQtumSatoshi := uint64(summary.Balance * 1e8)
	// Ether Wei = Qtum Satoshi * 10^10
	balanceInEtherWei := uint64(balanceInQtumSatoshi * 1e10)

	hexBalance := hexutil.EncodeUint64(balanceInEtherWei)

	return eth.GetBalanceResponse(hexBalance), nil
}
