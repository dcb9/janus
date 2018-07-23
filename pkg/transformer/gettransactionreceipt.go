package transformer

import (
	"encoding/json"
	"errors"

	"github.com/bitly/go-simplejson"
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (m *Manager) GetTransactionReceipt(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	var params []string
	if err := unmarshalRequest(req.Params, &params); err != nil {
		return nil, &rpc.JSONRPCError{
			Code:    rpc.ErrInvalid,
			Message: "invalid input",
		}
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	newParams, err := json.Marshal([]string{
		RemoveHexPrefix(params[0]),
	})
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = qtum.MethodGettransactionreceipt

	//Qtum RPC
	//gettransactionreceipt "hash"
	//  requires -logevents to be enabled
	//  Argument:
	//  1. "hash"          (string, required) The transaction hash

	return m.GettransactionreceiptResp, nil
}

func (m *Manager) GettransactionreceiptResp(result json.RawMessage) (interface{}, error) {
	if string(result) == "[]" {
		return nil, nil
	}

	sj, err := simplejson.NewJson(result)
	if err != nil {
		return nil, err
	}

	receiptBytes, err := sj.GetIndex(0).Encode()
	if err != nil {
		return nil, err
	}
	var receipt *qtum.TransactionReceipt
	if err = json.Unmarshal(receiptBytes, &receipt); err != nil {
		return nil, err
	}

	status := "0x0"
	if receipt.Excepted == "None" {
		status = "0x1"
	}

	logs := make([]eth.Log, 0, len(receipt.Log))
	for index, log := range receipt.Log {
		topics := make([]string, 0, len(log.Topics))
		for _, topic := range log.Topics {
			topics = append(topics, AddHexPrefix(topic))
		}
		logs = append(logs, eth.Log{
			TransactionHash:  AddHexPrefix(receipt.TransactionHash),
			TransactionIndex: hexutil.EncodeUint64(receipt.TransactionIndex),
			BlockHash:        AddHexPrefix(receipt.BlockHash),
			BlockNumber:      hexutil.EncodeUint64(receipt.BlockNumber),
			Data:             AddHexPrefix(log.Data),
			Address:          AddHexPrefix(log.Address),
			Topics:           topics,
			LogIndex:         hexutil.EncodeUint64(uint64(index)),
		})
	}

	ethTxReceipt := eth.TransactionReceipt{
		TransactionHash:   AddHexPrefix(receipt.TransactionHash),
		TransactionIndex:  hexutil.EncodeUint64(receipt.TransactionIndex),
		BlockHash:         AddHexPrefix(receipt.BlockHash),
		BlockNumber:       hexutil.EncodeUint64(receipt.BlockNumber),
		ContractAddress:   AddHexPrefix(receipt.ContractAddress),
		CumulativeGasUsed: hexutil.EncodeUint64(receipt.CumulativeGasUsed),
		GasUsed:           hexutil.EncodeUint64(receipt.GasUsed),
		Logs:              logs,
		Status:            status,

		// see Known issues
		LogsBloom: "",
	}

	return &ethTxReceipt, nil
}
