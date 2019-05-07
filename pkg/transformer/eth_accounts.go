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

func (p *ProxyETHAccounts) request() (eth.AccountsResponse, error) {
	// eth req -> qtum req
	qtumreq := qtum.GetAddressesByAccountRequest("")

	var qtumresp qtum.GetAddressesByAccountResponse
	if err := p.Qtum.Request(qtum.MethodGetAddressesByAccount, &qtumreq, &qtumresp); err != nil {
		return nil, err
	}

	// qtum res -> eth res
	var accounts eth.AccountsResponse
	for _, base58Addr := range qtumresp {
		addr, err := p.Base58AddressToHex(base58Addr)

		// discard addresses that cannot be converted to hex format (i.e. multisig, segwit)
		if err != nil {
			continue
		}

		accounts = append(accounts, utils.AddHexPrefix(addr))
	}

	return accounts, nil
}

func (p *ProxyETHAccounts) ToResponse(ethresp *qtum.CallContractResponse) *eth.CallResponse {
	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
	qtumresp := eth.CallResponse(data)
	return &qtumresp
}
