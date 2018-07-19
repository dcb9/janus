package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/rpc"
)

func sendtoaddress(req *rpc.JSONRPCRequest, tx *eth.TransactionReq) (*rpc.JSONRPCRequest, error) {
	// fixme
	return req, nil
}
