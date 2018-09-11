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

func (m *Method) GetBlockHeader(hash string) (resp *GetBlockHeaderResponse, err error) {
	req := GetBlockHeaderRequest{
		Hash: hash,
	}
	err = m.Request(MethodGetBlockHeader, &req, &resp)
	return
}

func (m *Method) GetBlock(hash string) (resp *GetBlockResponse, err error) {
	req := GetBlockRequest{
		Hash: hash,
	}
	err = m.Request(MethodGetBlock, &req, &resp)
	return
}

func (m *Method) Generate(blockNum int, maxTries *int) (resp GenerateResponse, err error) {
	req := GenerateRequest{
		BlockNum: blockNum,
		MaxTries: maxTries,
	}
	err = m.Request(MethodGenerate, &req, &resp)
	return
}

func (m *Method) SearchLogs(req *SearchLogsRequest) (receipts SearchLogsResponse, err error) {
	if err := m.Request(MethodSearchLogs, req, &receipts); err != nil {
		return nil, err
	}
	return
}
