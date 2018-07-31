package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProxyETHGetTransactionReceipt implements ETHProxy
type ProxyETHGetTransactionReceipt struct {
	*qtum.Qtum
}

func (p *ProxyETHGetTransactionReceipt) Method() string {
	return "eth_getTransactionReceipt"
}

func (p *ProxyETHGetTransactionReceipt) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req *eth.GetTransactionReceiptRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	qtumreq, err := p.ToRequest(req)
	if err != nil {
		return nil, err
	}

	return p.request(qtumreq)
}

func (p *ProxyETHGetTransactionReceipt) request(req *qtum.GetTransactionReceiptRequest) (*eth.GetTransactionReceiptResponse, error) {
	var receipt qtum.GetTransactionReceiptResponse
	if err := p.Qtum.Request(qtum.MethodGetTransactionReceipt, req, &receipt); err != nil {
		if err == qtum.EmptyResponseErr {
			return nil, nil
		}
		return nil, err
	}

	status := "0x0"
	if receipt.Excepted == "None" {
		status = "0x1"
	}

	r := qtum.TransactionReceiptStruct(receipt)
	logs := getEthLogs(&r)

	ethTxReceipt := eth.GetTransactionReceiptResponse{
		TransactionHash:   utils.AddHexPrefix(receipt.TransactionHash),
		TransactionIndex:  hexutil.EncodeUint64(receipt.TransactionIndex),
		BlockHash:         utils.AddHexPrefix(receipt.BlockHash),
		BlockNumber:       hexutil.EncodeUint64(receipt.BlockNumber),
		ContractAddress:   utils.AddHexPrefix(receipt.ContractAddress),
		CumulativeGasUsed: hexutil.EncodeUint64(receipt.CumulativeGasUsed),
		GasUsed:           hexutil.EncodeUint64(receipt.GasUsed),
		Logs:              logs,
		Status:            status,

		// see Known issues
		LogsBloom: "",
	}

	return &ethTxReceipt, nil
}

func (p *ProxyETHGetTransactionReceipt) ToRequest(ethreq *eth.GetTransactionReceiptRequest) (*qtum.GetTransactionReceiptRequest, error) {
	qtumreq := qtum.GetTransactionReceiptRequest(utils.RemoveHexPrefix(string(*ethreq)))
	return &qtumreq, nil
}
