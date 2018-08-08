package qtum

import (
	"math/big"

	"github.com/dcb9/janus/pkg/utils"
)

type Method struct {
	*Client
}

func (m *Method) Base58AddressToHex(addr string) (string, error) {
	var response GetHexAddressResponse
	err := m.Request(MethodGetHexAddress, GetHexAddressRequest(addr), &response)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func (m *Method) FromHexAddress(addr string) (string, error) {
	addr = utils.RemoveHexPrefix(addr)

	var response FromHexAddressResponse
	err := m.Request(MethodFromHexAddress, FromHexAddressRequest(addr), &response)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func (m *Method) GetTransactionReceipt(txHash string) (*GetTransactionReceiptResponse, error) {
	var resp *GetTransactionReceiptResponse
	err := m.Request(MethodGetTransactionReceipt, GetTransactionReceiptRequest(txHash), &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *Method) DecodeRawTransaction(hex string) (*DecodedRawTransactionResponse, error) {
	var resp *DecodedRawTransactionResponse
	err := m.Request(MethodDecodeRawTransaction, DecodeRawTransactionRequest(hex), &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *Method) GetBlockCount() (resp *GetBlockCountResponse, err error) {
	err = m.Request(MethodGetBlockCount, nil, &resp)
	return
}

func (m *Method) GetBlockHash(b *big.Int) (resp GetBlockHashResponse, err error) {
	req := GetBlockHashRequest{
		Int: b,
	}
	err = m.Request(MethodGetBlockHash, &req, &resp)
	return
}
