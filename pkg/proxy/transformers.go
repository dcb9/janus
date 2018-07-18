package proxy

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/dcb9/janus/pkg/qtum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

func transformTransactionByHash(req *qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error) {
	var params []string
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	newParams, err := json.Marshal([]interface{}{
		params[0],
	})
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = "gettransaction"

	//Qtum RPC
	// gettransaction "txid" ( include_watchonly ) (waitconf)
	//
	// Get detailed information about in-wallet transaction <txid>
	//
	// Arguments:
	// 1. "txid"                  (string, required) The transaction id
	// 2. "include_watchonly"     (bool, optional, default=false) Whether to include watch-only addresses in balance calculation and details[]
	// 3. "waitconf"              (int, optional, default=0) Wait for enough confirmations before returning
	//
	// Result:
	// {
	//   "amount" : x.xxx,        (numeric) The transaction amount in QTUM
	//   "fee": x.xxx,            (numeric) The amount of the fee in QTUM. This is negative and only available for the
	//                               'send' category of transactions.
	//   "confirmations" : n,     (numeric) The number of confirmations
	//   "blockhash" : "hash",  (string) The block hash
	//   "blockindex" : xx,       (numeric) The index of the transaction in the block that includes it
	//   "blocktime" : ttt,       (numeric) The time in seconds since epoch (1 Jan 1970 GMT)
	//   "txid" : "transactionid",   (string) The transaction id.
	//   "time" : ttt,            (numeric) The transaction time in seconds since epoch (1 Jan 1970 GMT)
	//   "timereceived" : ttt,    (numeric) The time received in seconds since epoch (1 Jan 1970 GMT)
	//   "bip125-replaceable": "yes|no|unknown",  (string) Whether this transaction could be replaced due to BIP125 (replace-by-fee);
	//                                                    may be unknown for unconfirmed transactions not in the mempool
	//   "details" : [
	//     {
	//       "account" : "accountname",      (string) DEPRECATED. The account name involved in the transaction, can be "" for the default account.
	//       "address" : "address",          (string) The qtum address involved in the transaction
	//       "category" : "send|receive",    (string) The category, either 'send' or 'receive'
	//       "amount" : x.xxx,                 (numeric) The amount in QTUM
	//       "label" : "label",              (string) A comment for the address/transaction, if any
	//       "vout" : n,                       (numeric) the vout value
	//       "fee": x.xxx,                     (numeric) The amount of the fee in QTUM. This is negative and only available for the
	//                                            'send' category of transactions.
	//       "abandoned": xxx                  (bool) 'true' if the transaction has been abandoned (inputs are respendable). Only available for the
	//                                            'send' category of transactions.
	//     }
	//     ,...
	//   ],
	//   "hex" : "data"         (string) Raw data for transaction
	// }
	//
	// Examples:
	// > qtum-cli gettransaction "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"
	// > qtum-cli gettransaction "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d" true
	// > curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "gettransaction", "params": ["1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"] }' -H 'content-type: text/plain;' http://127.0.0.1:3889/
	return req, nil
}
func transformTransactionReceipt(req *qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error) {
	var params []json.RawMessage
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	newParams, err := json.Marshal([]interface{}{
		params[0],
	})
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = "gettransactionreceipt"

	//Qtum RPC
	//gettransactionreceipt "hash"
	//  requires -logevents to be enabled
	//  Argument:
	//  1. "hash"          (string, required) The transaction hash

	return req, nil
}
func transformCall(req *qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error) {
	var params []json.RawMessage
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	var txCall map[string]string
	if err := json.Unmarshal(params[0], &txCall); err != nil {
		return nil, err
	}
	gasLimit, _, err := getGasPriceAndGasLimit(txCall)
	if err != nil {
		return nil, err
	}

	newParams, err := json.Marshal([]interface{}{
		txCall["to"],
		txCall["data"],
		txCall["from"],
		gasLimit,
	})
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = "callcontract"

	//Qtum RPC
	// callcontract "address" "data" ( address )
	//
	// Argument:
	//   1. "address"          (string, required) The account address
	//   2. "data"             (string, required) The data hex string
	//   3. address              (string, optional) The sender address hex string
	//   4. gasLimit             (string, optional) The gas limit for executing the contract

	return req, nil
}

func transformSendTransaction(req *qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error) {
	var params []map[string]string
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	transaction := params[0]
	if transaction["to"] == "" {
		return createcontract(req, transaction)
	}

	return sendtocontract(req, transaction)
}

type transformerFunc func(*qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error)

func (fn transformerFunc) transform(req *qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error) {
	return fn(req)
}

type transformer interface {
	transform(*qtum.JSONRPCRequest) (*qtum.JSONRPCRequest, error)
}

func createcontract(req *qtum.JSONRPCRequest, transaction map[string]string) (*qtum.JSONRPCRequest, error) {
	if v, ok := transaction["value"]; ok {
		if v != "" && v != "0x0" {
			return nil, &qtum.JSONRPCError{
				Code:    ErrInvalid,
				Message: "value must be empty",
			}
		}
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
	gasLimit, gasPrice, err := getGasPriceAndGasLimit(transaction)
	if err != nil {
		return nil, err
	}
	params := []interface{}{transaction["data"], gasLimit, gasPrice}

	if f, ok := transaction["from"]; ok {
		sender := f
		if strings.HasPrefix(f, "0x") {
			// todo convert hexaddress
		}

		params = append(params, sender)
	}

	newParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = "createcontract"
	return req, nil
}

func sendtocontract(req *qtum.JSONRPCRequest, transaction map[string]string) (*qtum.JSONRPCRequest, error) {
	//  Eth RPC
	//  params: [{
	//    "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
	//    "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
	//    "gas": "0x76c0", // 30400
	//    "gasPrice": "0x9184e72a000", // 10000000000000
	//    "value": "0x9184e72a", // 2441406250
	//    "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
	//  }]

	//Qtum RPC
	//  sendtocontract "contractaddress" "data" (amount gaslimit gasprice senderaddress broadcast)
	//  Send funds and data to a contract.
	//
	//Arguments:
	//  1. "contractaddress" (string, required) The contract address that will receive the funds and data.
	//  2. "datahex"  (string, required) data to send.
	//  3. "amount"      (numeric or string, optional) The amount in QTUM to send. eg 0.1, default: 0
	//  4. gasLimit  (numeric or string, optional) gasLimit, default: 250000, max: 40000000
	//  5. gasPrice  (numeric or string, optional) gasPrice Qtum price per gas unit, default: 0.0000004, min:0.0000004
	//  6. "senderaddress" (string, optional) The quantum address that will be used as sender.
	//  7. "broadcast" (bool, optional, default=true) Whether to broadcast the transaction or not.
	//  8. "changeToSender" (bool, optional, default=true) Return the change to the sender.
	//
	//Result:
	//  [
	//  {
	//  "txid" : (string) The transaction id.
	//  "sender" : (string) QTUM address of the sender.
	//  "hash160" : (string) ripemd-160 hash of the sender.
	//  }
	//  ]
	//
	//Examples:
	//  > qtum-cli sendtocontract "c6ca2697719d00446d4ea51f6fac8fd1e9310214" "54f6127f"
	//  > qtum-cli sendtocontract "c6ca2697719d00446d4ea51f6fac8fd1e9310214" "54f6127f" 12.0015 6000000 0.0000004 "QM72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd"
	gasLimit, gasPrice, err := getGasPriceAndGasLimit(transaction)
	if err != nil {
		return nil, err
	}

	amount := 0
	if v, ok := transaction["value"]; ok {
		_ = v
		// FIXME
		// amount = v
	}
	params := []interface{}{transaction["to"], transaction["data"], amount, gasLimit, gasPrice}

	if f, ok := transaction["from"]; ok {
		sender := f
		if strings.HasPrefix(f, "0x") {
			// todo convert hexaddress
		}

		params = append(params, sender)
	}

	newParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = "sendtocontract"

	return req, nil
}

func getGasPriceAndGasLimit(tx map[string]string) (gasLimit *big.Int, gasPrice string, err error) {
	gasLimit = big.NewInt(2500000)
	if v, ok := tx["gas"]; ok {
		gasLimit, err = hexutil.DecodeBig(v)
		if err != nil {
			err = errors.Wrap(err, "decode gas")
			return
		}
	}
	gasPrice = "0.0000004"
	if _, ok := tx["gasPrice"]; ok {
		// fixme parse gas price
	}

	return
}
