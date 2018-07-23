package transformer

import (
	"encoding/json"
	"math/big"

	"github.com/bitly/go-simplejson"
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
		"eth_blockNumber":           m.BlockNumber,
		"net_version":               m.NetVersion,
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

	return nil, nil
}

func (m *Manager) SendTransaction(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	var txs []*eth.TransactionReq
	if err := unmarshalRequest(req.Params, &txs); err != nil {
		return nil, err
	}
	if len(txs) == 0 {
		return nil, errors.New("params must be set")
	}

	var err error
	t := txs[0]
	if t.IsCreateContract() {
		err = m.createcontract(req, t)
	} else if t.IsSendEther() {
		err = m.sendtoaddress(req, t)
	} else if t.IsCallContract() {
		err = m.sendtocontract(req, t)
	} else {
		err = &rpc.JSONRPCError{
			Code:    rpc.ErrUnknownOperation,
			Message: "unknown operation",
		}
	}
	if err != nil {
		return nil, err
	}

	return m.sendTransactionResp, nil
}
func (m *Manager) sendTransactionResp(result json.RawMessage) (interface{}, error) {
	sj, err := simplejson.NewJson(result)
	if err != nil {
		return nil, err
	}
	txid, err := sj.Get("txid").String()
	if err != nil {
		return nil, err
	}

	return AddHexPrefix(txid), nil
}

func (m *Manager) BlockNumber(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	req.Method = qtum.MethodGetblockcount

	return func(result json.RawMessage) (newResult interface{}, err error) {
		var n *big.Int
		if err := json.Unmarshal(result, &n); err != nil {
			return nil, err
		}
		return hexutil.EncodeBig(n), nil
	}, nil
}

func (m *Manager) NetVersion(req *rpc.JSONRPCRequest) (ResponseTransformerFunc, error) {
	req.Method = qtum.MethodGetblockchaininfo

	return func(result json.RawMessage) (newResult interface{}, err error) {
		var i *qtum.BlockChainInfo
		if err := json.Unmarshal(result, &i); err != nil {
			return nil, err
		}

		return i.Chain, nil
	}, nil
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

func (m *Manager) getQtumWalletAddress(addr string) (string, error) {
	if IsEthHexAddress(addr) {
		return m.qtumClient.FromHexAddress(RemoveHexPrefix(addr))
	}
	return addr, nil
}
