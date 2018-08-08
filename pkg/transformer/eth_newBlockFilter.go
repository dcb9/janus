package transformer

import (
	"sync"

	"github.com/dcb9/janus/pkg/eth"
	"github.com/dcb9/janus/pkg/qtum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// ProxyETHNewBlockFilter implements ETHProxy
type ProxyETHNewBlockFilter struct {
	*qtum.Qtum
	blockFilter *BlockFilterSimulator
}

func (p *ProxyETHNewBlockFilter) Method() string {
	return "eth_newBlockFilter"
}

func (p *ProxyETHNewBlockFilter) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHNewBlockFilter) request() (eth.NewBlockFilterResponse, error) {
	blockCount, err := p.GetBlockCount()
	if err != nil {
		return "", err
	}

	id := p.blockFilter.New(blockCount.Int.Uint64())

	return eth.NewBlockFilterResponse(hexutil.EncodeUint64(id)), nil
}

// BlockFilterSimulatorr a map filter id => last block number
type BlockFilterSimulator struct {
	filters       sync.Map
	filterIDMutex sync.Mutex
	maxFilterID   uint64
}

func (f *BlockFilterSimulator) New(blockNumber uint64) uint64 {
	f.filterIDMutex.Lock()
	f.maxFilterID++
	id := f.maxFilterID
	f.filterIDMutex.Unlock()

	f.filters.Store(id, blockNumber)

	return id
}

func (f *BlockFilterSimulator) SetBlockNumber(filterID uint64, blockNumber uint64) error {
	if _, ok := f.filters.Load(filterID); !ok {
		return errors.Errorf("filter id %d does not exist", filterID)
	}

	f.filters.Store(filterID, blockNumber)

	return nil
}

func (f *BlockFilterSimulator) GetBlockNumber(filterID uint64) (uint64, error) {
	val, ok := f.filters.Load(filterID)
	if !ok {
		return 0, errors.Errorf("filter id %d does not exist", filterID)
	}
	blockNumber := val.(uint64)

	return blockNumber, nil
}

func (f *BlockFilterSimulator) Uninstall(filterID uint64) {
	f.filters.Delete(filterID)
}
