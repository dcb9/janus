package transformer

import (
	"encoding/json"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
)

func (m *Manager) sendtoaddress(req *rpc.JSONRPCRequest, tx *eth.TransactionReq) (ResponseTransformerFunc, error) {
	// todo
	req.Method = qtum.MethodSendtoaddress

	return m.SendtoaddressResp, nil
}

func (m *Manager) SendtoaddressResp(result json.RawMessage) (interface{}, error) {
	var i interface{}
	err := json.Unmarshal(result, &i)
	if err != nil {
		return nil, err
	}

	return i, nil
}
