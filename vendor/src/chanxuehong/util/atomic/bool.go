package atomic

import (
	"sync/atomic"
	"unsafe"
)

type Bool uint32 // zero value represents false

func (b *Bool) Load() (val bool) {
	n := atomic.LoadUint32((*uint32)(unsafe.Pointer(b)))
	return n != 0
}

func (b *Bool) Store(val bool) {
	if val {
		atomic.StoreUint32((*uint32)(unsafe.Pointer(b)), 1)
	} else {
		atomic.StoreUint32((*uint32)(unsafe.Pointer(b)), 0)
	}
}

func (b *Bool) Swap(new bool) (old bool) {
	var _new uint32
	if new {
		_new = 1
	}
	_old := atomic.SwapUint32((*uint32)(unsafe.Pointer(b)), _new)
	return _old != 0
}

func (b *Bool) CompareAndSwap(old, new bool) (swapped bool) {
	var _old, _new uint32
	if old {
		_old = 1
	}
	if new {
		_new = 1
	}
	return atomic.CompareAndSwapUint32((*uint32)(unsafe.Pointer(b)), _old, _new)
}
