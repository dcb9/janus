package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/go-kit/kit/log"
)

func (m *Manager) sendtoaddress(req *rpc.JSONRPCRequest, tx *eth.TransactionReq) (ResponseTransformerFunc, error) {
	// todo
	req.Method = qtum.MethodSendtoaddress

	l := log.WithPrefix(m.logger, "method", req.Method)
	return func(result *rpc.JSONRPCResult) error {
		return m.SendtoaddressResp(context{
			logger: l,
			req:    req,
		}, result)
	}, nil
}

func (m *Manager) SendtoaddressResp(c context, result *rpc.JSONRPCResult) error {
	result.JSONRPC = "2.0"
	return nil
}
