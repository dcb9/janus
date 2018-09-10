package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
)

// ProxyETHAccounts implements ETHProxy
type ProxyETHAccounts struct {
	*qtum.Qtum
}

func (p *ProxyETHAccounts) Method() string {
	return "eth_accounts"
}

func (p *ProxyETHAccounts) Request(_ *eth.JSONRPCRequest) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHAccounts) request() (ethresp eth.AccountsResponse, err error) {
	ethresp = make(eth.AccountsResponse, 0)

	// eth req -> qtum req
	qtumreq := qtum.GetAddressesByAccountRequest("")

	var qtumresp qtum.GetAddressesByAccountResponse
	if err = p.Qtum.Request(qtum.MethodGetAddressesByAccount, &qtumreq, &qtumresp); err != nil {
		return
	}

	// qtum res -> eth res
	var addr string
	for _, base58Addr := range qtumresp {
		addr, err = p.Base58AddressToHex(base58Addr)
		if err != nil {
			return
		}

		ethresp = append(ethresp, utils.AddHexPrefix(addr))
	}

	return
}

func (p *ProxyETHAccounts) ToResponse(ethresp *qtum.CallContractResponse) *eth.CallResponse {
	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
	qtumresp := eth.CallResponse(data)
	return &qtumresp
}
