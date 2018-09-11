package eth

import (
	"math/big"
	"sync"
	"sync/atomic"
)

type FilterType int

const (
	NewFilterTy FilterType = iota
	NewBlockFilterTy
	NewPendingTransactionFilterTy
)

type Filter struct {
	ID           uint64
	Type         FilterType
	Request      interface{}
	LastBlockNum *big.Int
	Data         sync.Map
}

type FilterSimulator struct {
	filters     sync.Map
	maxFilterID *uint64
}

func NewFilterSimulator() *FilterSimulator {
	id := uint64(0)
	return &FilterSimulator{
		maxFilterID: &id,
	}
}

func (f *FilterSimulator) New(ty FilterType, req ...interface{}) *Filter {
	id := atomic.AddUint64(f.maxFilterID, 1)
	filter := &Filter{ID: id, Type: ty}
	if ty == NewFilterTy {
		filter.Request = req[0]
	}

	f.filters.Store(id, filter)

	return filter
}

func (f *FilterSimulator) Uninstall(filterID uint64) {
	f.filters.Delete(filterID)
}

func (f *FilterSimulator) Filter(filterID uint64) (value interface{}, ok bool) {
	return f.filters.Load(filterID)
}
