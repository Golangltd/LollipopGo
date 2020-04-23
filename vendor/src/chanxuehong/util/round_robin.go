package util

import (
	"fmt"
	"math"
	"sync/atomic"
)

func NewRoundRobinIndex(bound int) *RoundRobinIndex {
	if bound <= 0 {
		panic(fmt.Sprintf("invalid bound: %d", bound))
	}
	return &RoundRobinIndex{
		bound: uint64(bound),
		index: math.MaxUint64,
	}
}

type RoundRobinIndex struct {
	bound uint64
	index uint64
}

func (rr *RoundRobinIndex) Next() int {
	index := atomic.AddUint64(&rr.index, 1)
	return int(index % rr.bound)
}
