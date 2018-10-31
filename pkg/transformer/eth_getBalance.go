package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

	addr := utils.RemoveHexPrefix(req.Address)
	{
		// is address a contract or an account?
		qtumreq := qtum.GetAccountInfoRequest(addr)
		qtumresp, err := p.GetAccountInfo(&qtumreq)

		// the address is a contract
		if err == nil {
			// the unit of the balance Satoshi
			return hexutil.EncodeUint64(uint64(qtumresp.Balance)), nil
		}
	}

	{
		// try account
		base58Addr, err := p.FromHexAddress(addr)
		if err != nil {
			return nil, err
		}

		qtumreq := qtum.NewListUnspentRequest(base58Addr)
		qtumresp, err := p.ListUnspent(qtumreq)
		if err != nil {
			return nil, err
		}

		balance := float64(0)
		for _, utxo := range *qtumresp {
			balance += utxo.Amount
		}

		// 1 QTUM = 10 ^ 8 Satoshi
		return hexutil.EncodeUint64(uint64(balance * 1e8)), nil
	}
}
