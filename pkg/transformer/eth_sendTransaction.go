package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/pkg/errors"
)

// ProxyETHSendTransaction implements ETHProxy
type ProxyETHSendTransaction struct {
	*qtum.Qtum
}

func (p *ProxyETHSendTransaction) Method() string {
	return "eth_sendTransaction"
}

func (p *ProxyETHSendTransaction) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.SendTransactionRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	if req.IsCreateContract() {
		return p.requestCreateContract(&req)
	} else if req.IsSendEther() {
		return p.requestSendToAddress(&req)
	} else if req.IsCallContract() {
		return p.requestSendToContract(&req)
	}

	return nil, errors.New("Unknown operation")
}

func (p *ProxyETHSendTransaction) requestSendToContract(ethtx *eth.SendTransactionRequest) (*eth.SendTransactionResponse, error) {
	gasLimit, gasPrice, err := EthGasToQtum(ethtx)
	if err != nil {
		return nil, err
	}

	amount := 0.0
	if ethtx.Value != "" {
		var err error
		amount, err = EthValueToQtumAmount(ethtx.Value)
		if err != nil {
			return nil, errors.Wrap(err, "EthValueToQtumAmount:")
		}
	}

	qtumreq := qtum.SendToContractRequest{
		ContractAddress: utils.RemoveHexPrefix(ethtx.To),
		Datahex:         utils.RemoveHexPrefix(ethtx.Data),
		Amount:          amount,
		GasLimit:        gasLimit,
		GasPrice:        gasPrice,
	}

	if from := ethtx.From; from != "" && utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return nil, err
		}
		qtumreq.SenderAddress = from
	}

	var resp *qtum.SendToContractResponse
	if err := p.Qtum.Request(qtum.MethodSendToContract, &qtumreq, &resp); err != nil {
		return nil, err
	}

	ethresp := eth.SendTransactionResponse(resp.Txid)
	return &ethresp, nil
}

func (p *ProxyETHSendTransaction) requestSendToAddress(req *eth.SendTransactionRequest) (*eth.SendTransactionResponse, error) {
	getQtumWalletAddress := func(addr string) (string, error) {
		if utils.IsEthHexAddress(addr) {
			return p.FromHexAddress(utils.RemoveHexPrefix(addr))
		}
		return addr, nil
	}

	from, err := getQtumWalletAddress(req.From)
	if err != nil {
		return nil, err
	}

	to, err := getQtumWalletAddress(req.To)
	if err != nil {
		return nil, err
	}

	amount, err := EthValueToQtumAmount(req.Value)
	if err != nil {
		return nil, err
	}

	qtumreq := qtum.SendToAddressRequest{
		Address:       from,
		Amount:        amount,
		SenderAddress: to,
	}

	var qtumresp qtum.SendToAddressResponse
	if err := p.Qtum.Request(qtum.MethodSendToAddress, &qtumreq, &qtumresp); err != nil {
		return nil, err
	}

	ethresp := eth.SendTransactionResponse(utils.AddHexPrefix(string(qtumresp)))

	return &ethresp, nil
}

func (p *ProxyETHSendTransaction) requestCreateContract(req *eth.SendTransactionRequest) (*eth.SendTransactionResponse, error) {
	gasLimit, gasPrice, err := EthGasToQtum(req)
	if err != nil {
		return nil, err
	}

	qtumreq := &qtum.CreateContractRequest{
		ByteCode: utils.RemoveHexPrefix(req.Data),
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}

	var resp *qtum.CreateContractResponse
	if err := p.Qtum.Request(qtum.MethodCreateContract, qtumreq, &resp); err != nil {
		return nil, err
	}

	ethresp := eth.SendTransactionResponse(utils.AddHexPrefix(string(resp.Txid)))

	return &ethresp, nil
}
