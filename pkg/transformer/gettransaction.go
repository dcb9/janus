package transformer

import (
	"encoding/json"
	"errors"

	"github.com/dcb9/janus/pkg/rpc"
)

func transformTransactionByHash(req *rpc.JSONRPCRequest) (*rpc.JSONRPCRequest, error) {
	var params []string
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}
	if len(params) == 0 {
		return nil, errors.New("params must be set")
	}

	txid := EthHexToQtum(params[0])
	newParams, err := json.Marshal([]interface{}{
		txid,
	})
	if err != nil {
		return nil, err
	}

	req.Params = newParams
	req.Method = "gettransaction"

	return req, nil
}

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
