package transformer

import (
	"encoding/json"
	"math/big"

	"github.com/pkg/errors"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProxyETHGetFilterChanges implements ETHProxy
type ProxyETHGetFilterChanges struct {
	*qtum.Qtum
	filter *eth.FilterSimulator
}

func (p *ProxyETHGetFilterChanges) Method() string {
	return "eth_getFilterChanges"
}

func (p *ProxyETHGetFilterChanges) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetFilterChangesRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	filterID, err := hexutil.DecodeUint64(string(req))
	if err != nil {
		return nil, err
	}

	_filter, ok := p.filter.Filter(filterID)
	if !ok {
		return nil, errors.New("Invalid filter id")
	}
	filter := _filter.(*eth.Filter)

	switch filter.Type {
	case eth.NewFilterTy:
		return p.requestFilter(filter)
	case eth.NewBlockFilterTy:
		return p.requestBlockFilter(filter)
	case eth.NewPendingTransactionFilterTy:
		fallthrough
	default:

		return nil, errors.New("Unknown filter type")
	}
}

func (p *ProxyETHGetFilterChanges) requestBlockFilter(filter *eth.Filter) (qtumresp eth.GetFilterChangesResponse, err error) {
	qtumresp = make(eth.GetFilterChangesResponse, 0)

	_lastBlockNumber, ok := filter.Data.Load("lastBlockNumber")
	if !ok {
		return qtumresp, errors.New("Could not get lastBlockNumber")
	}
	lastBlockNumber := _lastBlockNumber.(uint64)

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
	filter.Data.Store("lastBlockNumber", blockCount)
	return
}
func (p *ProxyETHGetFilterChanges) requestFilter(filter *eth.Filter) (qtumresp eth.GetFilterChangesResponse, err error) {
	qtumresp = make(eth.GetFilterChangesResponse, 0)

	_lastBlockNumber, ok := filter.Data.Load("lastBlockNumber")
	if !ok {
		return qtumresp, errors.New("Could not get lastBlockNumber")
	}
	lastBlockNumber := _lastBlockNumber.(uint64)

	blockCountBigInt, err := p.GetBlockCount()
	if err != nil {
		return qtumresp, err
	}
	blockCount := blockCountBigInt.Uint64()

	differ := blockCount - lastBlockNumber

	if differ == 0 {
		return eth.GetFilterChangesResponse{}, nil
	}

	searchLogsReq, err := p.toSearchLogsReq(filter, big.NewInt(int64(lastBlockNumber+1)), big.NewInt(int64(blockCount)))
	if err != nil {
		return nil, err
	}

	return p.doSearchLogs(searchLogsReq)
}

func (p *ProxyETHGetFilterChanges) doSearchLogs(req *qtum.SearchLogsRequest) (eth.GetFilterChangesResponse, error) {
	resp, err := p.SearchLogs(req)
	if err != nil {
		return nil, err
	}

	receiptToResult := func(receipt *qtum.TransactionReceiptStruct) []interface{} {
		logs := getEthLogs(receipt)
		res := make([]interface{}, len(logs))
		for i, _ := range res {
			res[i] = logs[i]
		}
		return res
	}
	results := make(eth.GetFilterChangesResponse, 0)
	for _, receipt := range resp {
		r := qtum.TransactionReceiptStruct(receipt)
		results = append(results, receiptToResult(&r)...)
	}

	return results, nil
}

func (p *ProxyETHGetFilterChanges) toSearchLogsReq(filter *eth.Filter, from, to *big.Int) (*qtum.SearchLogsRequest, error) {
	ethreq := filter.Request.(*eth.NewFilterRequest)
	var err error
	var addresses []string
	if ethreq.Address != nil {
		if isString(ethreq.Address) {
			var addr string
			if err = json.Unmarshal(ethreq.Address, &addr); err != nil {
				return nil, err
			}
			addresses = append(addresses, addr)
		} else {
			if err = json.Unmarshal(ethreq.Address, &addresses); err != nil {
				return nil, err
			}
		}
		for i, _ := range addresses {
			addresses[i] = utils.RemoveHexPrefix(addresses[i])
		}
	}

	qtumreq := &qtum.SearchLogsRequest{
		Addresses: addresses,
		FromBlock: from,
		ToBlock:   to,
	}

	topics, ok := filter.Data.Load("topics")
	if ok {
		qtumreq.Topics = topics.([]interface{})
	}

	return qtumreq, nil
}
