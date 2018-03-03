/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月3日
*/

package controller

import (
	"sync"
)

// 控制器的结构体
type Tcontroller struct {
	info map[int]string
	mu   sync.RWMutex
}

// 申请对象
func NewController(paras ...interface{}) (m *Tcontroller) {

	return
}

//func newConcurrentMap3(initialCapacity int,
//	loadFactor float32, concurrencyLevel int) (m *ConcurrentMap) {
//	m = &ConcurrentMap{}

//	if !(loadFactor > 0) || initialCapacity < 0 || concurrencyLevel <= 0 {
//		panic(IllegalArgError)
//	}

//	if concurrencyLevel > MAX_SEGMENTS {
//		concurrencyLevel = MAX_SEGMENTS
//	}

//	// Find power-of-two sizes best matching arguments
//	sshift := 0
//	ssize := 1
//	for ssize < concurrencyLevel {
//		sshift++
//		ssize = ssize << 1
//	}

//	m.segmentShift = uint(32) - uint(sshift)
//	m.segmentMask = ssize - 1

//	m.segments = make([]*Segment, ssize)

//	if initialCapacity > MAXIMUM_CAPACITY {
//		initialCapacity = MAXIMUM_CAPACITY
//	}

//	c := initialCapacity / ssize
//	if c*ssize < initialCapacity {
//		c++
//	}
//	cap := 1
//	for cap < c {
//		cap <<= 1
//	}

//	for i := 0; i < len(m.segments); i++ {
//		m.segments[i] = m.newSegment(cap, loadFactor)
//	}
//	m.engChecker = new(Once)
//	return
//}

//func NewConcurrentMap(paras ...interface{}) (m *ConcurrentMap) {
//	ok := false
//	cap := DEFAULT_INITIAL_CAPACITY
//	factor := DEFAULT_LOAD_FACTOR
//	concurrent_lvl := DEFAULT_CONCURRENCY_LEVEL

//	if len(paras) >= 1 {
//		if cap, ok = paras[0].(int); !ok {
//			panic(IllegalArgError)
//		}
//	}

//	if len(paras) >= 2 {
//		if factor, ok = paras[1].(float32); !ok {
//			panic(IllegalArgError)
//		}
//	}

//	if len(paras) >= 3 {
//		if concurrent_lvl, ok = paras[2].(int); !ok {
//			panic(IllegalArgError)
//		}
//	}

//	m = newConcurrentMap3(cap, factor, concurrent_lvl)
//	return
//}
