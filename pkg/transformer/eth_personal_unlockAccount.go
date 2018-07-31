package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
)

// ProxyETHPersonalUnlockAccount implements ETHProxy
type ProxyETHPersonalUnlockAccount struct{}

func (p *ProxyETHPersonalUnlockAccount) Method() string {
	return "personal_unlockAccount"
}

func (p *ProxyETHPersonalUnlockAccount) Request(req *eth.JSONRPCRequest) (interface{}, error) {
	return eth.PersonalUnlockAccountResponse(true), nil
}
