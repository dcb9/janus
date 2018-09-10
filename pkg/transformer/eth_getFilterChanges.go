package transformer

import (
	"math/big"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProxyETHGetFilterChanges implements ETHProxy
type ProxyETHGetFilterChanges struct {
	*qtum.Qtum
	blockFilter *BlockFilterSimulator
}

func (p *ProxyETHGetFilterChanges) Method() string {
	return "eth_getFilterChanges"
}

func (p *ProxyETHGetFilterChanges) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetFilterChangesRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHGetFilterChanges) request(ethreq *eth.GetFilterChangesRequest) (qtumresp eth.GetFilterChangesResponse, err error) {
	qtumresp = make(eth.GetFilterChangesResponse, 0)
	filterID, err := hexutil.DecodeUint64(string(*ethreq))
	if err != nil {
		return qtumresp, err
	}

	lastBlockNumber, err := p.blockFilter.GetBlockNumber(filterID)
	if err != nil {
		return qtumresp, err
	}

	blockCountBigInt, err := p.GetBlockCount()
	if err != nil {
		return qtumresp, err
	}
	blockCount := blockCountBigInt.Uint64()

	differ := blockCount - lastBlockNumber

	hashes := make(eth.GetFilterChangesResponse, differ)
	for i, _ := range hashes {
		blockNumber := new(big.Int).SetUint64(lastBlockNumber + uint64(i) + 1)

		resp, err := p.GetBlockHash(blockNumber)
		if err != nil {
			return qtumresp, err
		}

		hashes[i] = utils.AddHexPrefix(string(resp))
	}

	qtumresp = hashes
	return
}
