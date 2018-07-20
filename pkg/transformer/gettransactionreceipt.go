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
		EthHexToQtum(params[0]),
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
	sj, err := simplejson.NewJson(result)
	if err != nil {
		return nil, err
	}
	sj = sj.GetIndex(0)
	transactionHash, err := sj.Get("transactionHash").String()
	if err != nil {
		return nil, err
	}
	blockHash, err := sj.Get("blockHash").String()
	if err != nil {
		return nil, err
	}
	contractAddress, err := sj.Get("contractAddress").String()
	if err != nil {
		return nil, err
	}

	transactionIndex, err := sj.Get("transactionIndex").Uint64()
	if err != nil {
		return nil, err
	}
	cumulativeGasUsed, err := sj.Get("cumulativeGasUsed").Uint64()
	if err != nil {
		return nil, err
	}
	gasUsed, err := sj.Get("gasUsed").Uint64()
	if err != nil {
		return nil, err
	}
	blockNumber, err := sj.Get("blockNumber").Uint64()
	if err != nil {
		return nil, err
	}

	var qtumLogs []qtum.Log
	qtumRawLog, err := sj.Get("log").Encode()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(qtumRawLog, &qtumLogs)
	if err != nil {
		return nil, err
	}
	logs := make([]eth.Log, 0, len(qtumLogs))
	for _, log := range qtumLogs {
		topics := make([]string, 0, len(log.Topics))
		for _, topic := range log.Topics {
			topics = append(topics, QtumHexToEth(topic))
		}
		logs = append(logs, eth.Log{
			Data:    QtumHexToEth(log.Data),
			Address: QtumHexToEth(log.Address),
			Topics:  topics,
		})
	}

	ethTxReceipt := eth.TransactionReceipt{
		TransactionHash:   QtumHexToEth(transactionHash),
		BlockHash:         QtumHexToEth(blockHash),
		ContractAddress:   QtumHexToEth(contractAddress),
		TransactionIndex:  hexutil.EncodeUint64(transactionIndex),
		CumulativeGasUsed: hexutil.EncodeUint64(cumulativeGasUsed),
		GasUsed:           hexutil.EncodeUint64(gasUsed),
		BlockNumber:       hexutil.EncodeUint64(blockNumber),
		Logs:              logs,

		// todo there must be a way to know if the transaction is succeeded
		Status: "0x1",

		// fixme
		LogsBloom: "",
	}

	return &ethTxReceipt, nil
}
