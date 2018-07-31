package transformer

import (
	"math/big"

	"fmt"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// ProxyETHGetTransactionByHash implements ETHProxy
type ProxyETHGetTransactionByHash struct {
	*qtum.Qtum
}

func (p *ProxyETHGetTransactionByHash) Method() string {
	return "eth_getTransactionByHash"
}

func (p *ProxyETHGetTransactionByHash) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetTransactionByHashRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	qtumreq := p.ToRequest(&req)

	return p.request(qtumreq)
}

func (p *ProxyETHGetTransactionByHash) request(req *qtum.GetTransactionRequest) (*eth.GetTransactionByHashResponse, error) {
	var tx *qtum.GetTransactionResponse
	if err := p.Qtum.Request(qtum.MethodGetTransaction, req, &tx); err != nil {
		fmt.Println("err", err.Error())
		if err == qtum.EmptyResponseErr {
			return nil, nil
		}

		return nil, err
	}

	ethVal, err := QtumAmountToEthValue(tx.Amount)
	if err != nil {
		return nil, err
	}

	decodedRawTx, err := p.Qtum.DecodeRawTransaction(tx.Hex)
	if err != nil {
		return nil, errors.Wrap(err, "Qtum#DecodeRawTransaction")
	}

	var gas, gasPrice, input string
	type asmWithGasGasPriceEncodedABI interface {
		CallData() string
		GasPrice() (*big.Int, error)
		GasLimit() (*big.Int, error)
	}

	var asm asmWithGasGasPriceEncodedABI
	for _, out := range decodedRawTx.Vout {
		switch out.ScriptPubKey.Type {
		case "call":
			if asm, err = qtum.ParseCallASM(out.ScriptPubKey.Asm); err != nil {
				return nil, err
			}
		case "create":
			if asm, err = qtum.ParseCreateASM(out.ScriptPubKey.Asm); err != nil {
				return nil, err
			}
		default:
			continue
		}
		break
	}

	if asm != nil {
		input = utils.AddHexPrefix(asm.CallData())
		gasLimitBigInt, err := asm.GasLimit()
		if err != nil {
			return nil, err
		}
		gasPriceBigInt, err := asm.GasPrice()
		if err != nil {
			return nil, err
		}
		gas = hexutil.EncodeBig(gasLimitBigInt)
		gasPrice = hexutil.EncodeBig(gasPriceBigInt)
	}

	ethTxResp := eth.GetTransactionByHashResponse{
		Hash:      utils.AddHexPrefix(tx.Txid),
		BlockHash: utils.AddHexPrefix(tx.Blockhash),
		Nonce:     "",
		Value:     ethVal,
		Input:     input,
		Gas:       gas,
		GasPrice:  gasPrice,
	}

	if asm != nil {
		receipt, err := p.Qtum.GetTransactionReceipt(tx.Txid)
		if err != nil && err != qtum.EmptyResponseErr {
			return nil, err
		}
		if receipt != nil {
			ethTxResp.BlockNumber = hexutil.EncodeUint64(receipt.BlockNumber)
			ethTxResp.TransactionIndex = hexutil.EncodeUint64(receipt.TransactionIndex)
			ethTxResp.From = utils.AddHexPrefix(receipt.From)
			ethTxResp.To = utils.AddHexPrefix(receipt.ContractAddress)
		}
	}

	return &ethTxResp, nil
}

func (p *ProxyETHGetTransactionByHash) ToRequest(ethreq *eth.GetTransactionByHashRequest) *qtum.GetTransactionRequest {
	return &qtum.GetTransactionRequest{
		Txid: utils.RemoveHexPrefix(string(*ethreq)),
	}
}
