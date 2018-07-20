package transformer

import (
	"encoding/json"

	"github.com/bitly/go-simplejson"
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
)

func (m *Manager) createcontract(req *rpc.JSONRPCRequest, tx *eth.TransactionReq) (ResponseTransformerFunc, error) {
	if tx.Value != "" && tx.Value != "0x0" {
		return nil, &rpc.JSONRPCError{
			Code:    rpc.ErrInvalid,
			Message: "value must be empty",
		}
	}

	gasLimit, gasPrice, err := EthGasToQtum(tx)
	if err != nil {
		return nil, err
	}
	params := []interface{}{
		EthHexToQtum(tx.Data),
		gasLimit,
		gasPrice,
	}

	if tx.From != "" {
		from := tx.From
		if IsEthHex(from) {
			from, err = m.qtumClient.FromHexAddress(EthHexToQtum(from))
			if err != nil {
				return nil, err
			}
		}

		params = append(params, from)
	}

	newParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = qtum.MethodCreatecontract

	return m.CreatecontractResp, nil
}

func (m *Manager) CreatecontractResp(result json.RawMessage) (interface{}, error) {
	sj, err := simplejson.NewJson(result)
	if err != nil {
		return nil, err
	}
	txid, err := sj.Get("txid").String()
	if err != nil {
		return nil, err
	}

	return QtumHexToEth(txid), nil
}

//  Eth RPC
//  params: [{
//    "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
//    "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
//    "gas": "0x76c0", // 30400
//    "gasPrice": "0x9184e72a000", // 10000000000000
//    "value": "",
//    "data": "0xd46e...675"
//  }]

//Qtum RPC
//  createcontract "bytecode" (gaslimit gasprice "senderaddress" broadcast)
//  Create a contract with bytcode.
//
//Arguments:
//  1. "bytecode"  (string, required) contract bytcode.
//  2. gasLimit  (numeric or string, optional) gasLimit, default: 2500000, max: 40000000
//  3. gasPrice  (numeric or string, optional) gasPrice QTUM price per gas unit, default: 0.0000004, min:0.0000004
//  4. "senderaddress" (string, optional) The quantum address that will be used to create the contract.
//  5. "broadcast" (bool, optional, default=true) Whether to broadcast the transaction or not.
//  6. "changeToSender" (bool, optional, default=true) Return the change to the sender.
//
//Result:
//	[
//	{
//		"txid" : (string) The transaction id.
//		"sender" : (string) QTUM address of the sender.
//		"hash160" : (string) ripemd-160 hash of the sender.
//		"address" : (string) expected contract address.
//	}
//	]
//
//Examples:
//	> qtum-cli createcontract "60606040525b33600060006101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690836c010000000000000000000000009081020402179055506103786001600050819055505b600c80605b6000396000f360606040526008565b600256"
//	> qtum-cli createcontract "60606040525b33600060006101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690836c010000000000000000000000009081020402179055506103786001600050819055505b600c80605b6000396000f360606040526008565b600256" 6000000 0.0000004 "QM72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd" true
