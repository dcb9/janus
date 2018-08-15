package transformer

import (
	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
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

func (p *ProxyETHGetBlockByNumber) request(ethreq *eth.GetBlockByNumberRequest) (*eth.GetBlockByNumberResponse, error) {
	blockNumber, err := getQtumBlockNumber(ethreq.BlockNumber, -1)
	if err != nil {
		return nil, err
	}

	height := blockNumber.Int64()
	if height == -1 {
		status, err := p.Insight.GetStatus()
		if err != nil {
			return nil, err
		}

		height = status.Info.Blocks
	}

	blockIndex, err := p.Insight.GetBlockIndex(height)
	if err != nil {
		return nil, err
	}

	block, err := p.Insight.GetBlock(blockIndex.BlockHash)
	if err != nil {
		return nil, err
	}

	/*
		Number           string        `json:"number"`
		Hash             string        `json:"hash"`
		ParentHash       string        `json:"parentHash"`
		Nonce            string        `json:"nonce"`
		Sha3Uncles       string        `json:"sha3Uncles"`
		LogsBloom        string        `json:"logsBloom"`
		TransactionsRoot string        `json:"transactionsRoot"`
		StateRoot        string        `json:"stateRoot"`
		Miner            string        `json:"miner"`
		Difficulty       string        `json:"difficulty"`
		TotalDifficulty  string        `json:"totalDifficulty"`
		ExtraData        string        `json:"extraData"`
		Size             string        `json:"size"`
		GasLimit         string        `json:"gasLimit"`
		GasUsed          string        `json:"gasUsed"`
		Timestamp        string        `json:"timestamp"`
		Transactions     []interface{} `json:"transactions"`
		Uncles           []string      `json:"uncles"`
	*/
	miner, err := p.Base58AddressToHex(block.MinedBy)
	if err != nil {
		return nil, err
	}

	ethresp := new(eth.GetBlockByNumberResponse)
	ethresp.Hash = block.Hash
	ethresp.Number = hexutil.EncodeUint64(uint64(height))
	ethresp.ParentHash = block.Previousblockhash
	ethresp.Nonce = hexutil.EncodeUint64(uint64(block.Nonce))
	ethresp.Miner = miner

	return ethresp, nil
}
