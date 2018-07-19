package transformer

import (
	"encoding/json"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

type Manager struct {
	qtumClient          *qtum.Client
	debug               bool
	logger              log.Logger
	requestTransformers map[string]RequestTransformerFunc
}

func NewManager(qtumClient *qtum.Client, opts ...func(*Manager) error) (*Manager, error) {
	if qtumClient == nil {
		return nil, errors.New("qtumClient cannot be nil")
	}

	m := &Manager{
		qtumClient: qtumClient,
		logger:     log.NewNopLogger(),
	}
	m.requestTransformers = map[string]RequestTransformerFunc{
		"eth_sendTransaction":       m.SendTransaction,
		"eth_call":                  m.Call,
		"eth_getTransactionByHash":  m.GetTransactionByHash,
		"eth_getTransactionReceipt": m.GetTransactionReceipt,
	}

	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Manager) TransformRequest(rpcReq *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	if trafo, ok := m.requestTransformers[rpcReq.Method]; ok {
		return trafo(rpcReq)
	}

	return NopRequestTransformer(nil)
}

func (m *Manager) SendTransaction(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	var txs []*eth.TransactionReq
	if err := json.Unmarshal(req.Params, &txs); err != nil {
		return nil, err
	}
	if len(txs) == 0 {
		return nil, errors.New("params must be set")
	}

	t := txs[0]
	if t.IsCreateContract() {
		return m.createcontract(req, t)
	} else if t.IsSendEther() {
		return m.sendtoaddress(req, t)
	} else if t.IsCallContract() {
		return m.sendtocontract(req, t)
	}

	return nil, &rpc.JSONRPCError{
		Code:    rpc.ErrUnknownOperation,
		Message: "unknown operation",
	}
}

func SetDebug(debug bool) func(*Manager) error {
	return func(m *Manager) error {
		m.debug = debug
		return nil
	}
}

func SetLogger(l log.Logger) func(*Manager) error {
	return func(m *Manager) error {
		m.logger = log.WithPrefix(l, "component", "transformer")
		return nil
	}
}
