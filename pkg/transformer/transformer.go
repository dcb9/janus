package transformer

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type RequestTransformerFunc func(*rpc.JSONRPCRequest) (*rpc.JSONRPCRequest, error)

func (fn RequestTransformerFunc) Transform(req *rpc.JSONRPCRequest) (*rpc.JSONRPCRequest, error) {
	return fn(req)
}

type RequestTransformer interface {
	Transform(*rpc.JSONRPCRequest) (*rpc.JSONRPCRequest, error)
}

var RequestTransformers = map[string]RequestTransformer{
	"eth_sendTransaction":       RequestTransformerFunc(transformSendTransaction),
	"eth_call":                  RequestTransformerFunc(transformCall),
	"eth_getTransactionByHash":  RequestTransformerFunc(transformTransactionByHash),
	"eth_getTransactionReceipt": RequestTransformerFunc(transformTransactionReceipt),
}

type ResponseTransformerFunc func(*rpc.JSONRPCResult) (*rpc.JSONRPCResult, error)

func (fn ResponseTransformerFunc) Transform(r *rpc.JSONRPCResult) (*rpc.JSONRPCResult, error) {
	return fn(r)
}

type ResponseTransformer interface {
	Transform(*rpc.JSONRPCResult) (*rpc.JSONRPCResult, error)
}

var ResponseTransformers = map[string]ResponseTransformer{}

func transformSendTransaction(req *rpc.JSONRPCRequest) (*rpc.JSONRPCRequest, error) {
	var txs []*eth.TransactionReq
	if err := json.Unmarshal(req.Params, &txs); err != nil {
		return nil, err
	}
	if len(txs) == 0 {
		return nil, errors.New("params must be set")
	}

	t := txs[0]
	if t.IsCreateContract() {
		return createcontract(req, t)
	} else if t.IsSendEther() {
		return sendtoaddress(req, t)
	} else if t.IsCallContract() {
		return sendtocontract(req, t)
	}

	return nil, &rpc.JSONRPCError{
		Code:    rpc.ErrUnknownOperation,
		Message: "unknown operation",
	}
}

type EthGas interface {
	GetGas() string
	GetGasPrice() string
}

func EthGasToQtum(g EthGas) (gasLimit *big.Int, gasPrice string, err error) {
	gasLimit = big.NewInt(2500000)
	if gas := g.GetGas(); gas != "" {
		gasLimit, err = hexutil.DecodeBig(gas)
		if err != nil {
			err = errors.Wrap(err, "decode gas")
			return
		}
	}
	gasPrice = "0.0000004"
	// fixme parse gas price

	return
}

func EthHexToQtum(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return hex[2:]
	}
	return hex
}
