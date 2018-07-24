package transformer

import (
	"encoding/json"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
)

func (m *Manager) sendtoaddress(req *rpc.JSONRPCRequest, tx *eth.TransactionReq) (ResponseTransformerFunc, error) {
	req.Method = qtum.MethodSendtoaddress

	from, err := m.getQtumWalletAddress(tx.From)
	if err != nil {
		return nil, err
	}
	to, err := m.getQtumWalletAddress(tx.To)
	if err != nil {
		return nil, err
	}
	amount, err := EthValueToQtumAmount(tx.Value)
	if err != nil {
		return nil, err
	}

	params := []interface{}{
		to,
		amount,
		"", // comment
		"", // comment_to
		false,
		nil,
		nil,
		nil,
		from,
		true,
	}

	req.Params, err = json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return m.SendtoaddressResp, nil
}
func (m *Manager) SendtoaddressResp(result json.RawMessage) (interface{}, error) {
	var txid string
	err := json.Unmarshal(result, &txid)
	if err != nil {
		return nil, err
	}

	return AddHexPrefix(txid), nil
}

//  $ qcli help sendtoaddress
//  sendtoaddress "address" amount ( "comment" "comment_to" subtractfeefromamount replaceable conf_target "estimate_mode")
//
//  Send an amount to a given address.
//
//  Arguments:
//    1. "address"            (string, required) The qtum address to send to.
//    2. "amount"             (numeric or string, required) The amount in QTUM to send. eg 0.1
//    3. "comment"            (string, optional) A comment used to store what the transaction is for.
//    This is not part of the transaction, just kept in your wallet.
//    4. "comment_to"         (string, optional) A comment to store the name of the person or organization
//    to which you're sending the transaction. This is not part of the
//    transaction, just kept in your wallet.
//    5. subtractfeefromamount  (boolean, optional, default=false) The fee will be deducted from the amount being sent.
//    The recipient will receive less qtums than you enter in the amount field.
//    6. replaceable            (boolean, optional) Allow this transaction to be replaced by a transaction with higher fees via BIP 125
//    7. conf_target            (numeric, optional) Confirmation target (in blocks)
//    8. "estimate_mode"      (string, optional, default=UNSET) The fee estimate mode, must be one of:
//    "UNSET"
//    "ECONOMICAL"
//    "CONSERVATIVE"
//    9. "senderaddress"      (string, optional) The quantum address that will be used to send money from.
//    10."changeToSender"     (bool, optional, default=false) Return the change to the sender.
//
//  Result:
//    "txid"                  (string) The transaction id.
//
//  Examples:
//  > qtum-cli sendtoaddress "QM72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd" 0.1
//  > qtum-cli sendtoaddress "QM72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd" 0.1 "donation" "seans outpost"
//  > qtum-cli sendtoaddress "QM72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd" 0.1 "" "" true
//  > qtum-cli sendtoaddress "QM72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd", 0.1, "donation", "seans outpost", false, null, null, "", "QX1GkJdye9WoUnrE2v6ZQhQ72EUVDtGXQX", true
//  > curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "sendtoaddress", "params": ["QM72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd", 0.1, "donation", "seans outpost"] }' -H 'content-type: text/plain;' http://127.0.0.1:3889/
//  > curl --user myusername --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "sendtoaddress", "params": ["QM72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd", 0.1, "donation", "seans outpost", false, null, null, "", "QX1GkJdye9WoUnrE2v6ZQhQ72EUVDtGXQX", true] }' -H 'content-type: text/plain;' http://127.0.0.1:3889/
