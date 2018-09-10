package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProxyETHGetBlockByNumber implements ETHProxy
type ProxyETHGetBlockByNumber struct {
	*qtum.Qtum
}

func (p *ProxyETHGetBlockByNumber) Method() string {
	return "eth_getBlockByNumber"
}

func (p *ProxyETHGetBlockByNumber) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetBlockByNumberRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}
func (p *ProxyETHGetBlockByNumber) request(req *eth.GetBlockByNumberRequest) (*eth.GetBlockByNumberResponse, error) {
	blockNum, err := getQtumBlockNumber(req.BlockNumber, 0)
	if err != nil {
		return nil, err
	}

	blockHash, err := p.GetBlockHash(blockNum)
	if err != nil {
		return nil, err
	}

	blockHeaderResp, err := p.GetBlockHeader(string(blockHash))
	if err != nil {
		return nil, err
	}

	blockResp, err := p.GetBlock(string(blockHash))
	if err != nil {
		return nil, err
	}

	txs := make([]string, 0, len(blockResp.Tx))
	for _, tx := range blockResp.Tx {
		txs = append(txs, utils.AddHexPrefix(tx))
	}

	return &eth.GetBlockByNumberResponse{
		Hash:             utils.AddHexPrefix(blockHeaderResp.Hash),
		Nonce:            hexutil.EncodeUint64(uint64(blockHeaderResp.Nonce)),
		Number:           hexutil.EncodeUint64(uint64(blockHeaderResp.Height)),
		ParentHash:       utils.AddHexPrefix(blockHeaderResp.Previousblockhash),
		Difficulty:       hexutil.EncodeUint64(uint64(blockHeaderResp.Difficulty)),
		Timestamp:        hexutil.EncodeUint64(blockHeaderResp.Time),
		StateRoot:        utils.AddHexPrefix(blockHeaderResp.HashStateRoot),
		Size:             hexutil.EncodeUint64(uint64(blockResp.Size)),
		Transactions:     txs,
		TransactionsRoot: utils.AddHexPrefix(blockResp.Merkleroot),

		ExtraData:       "0x0",
		Miner:           "0x0000000000000000000000000000000000000000",
		TotalDifficulty: "0x0",
		GasLimit:        "0x0",
		GasUsed:         "0x0",
		Uncles:          []string{},
	}, nil
}
