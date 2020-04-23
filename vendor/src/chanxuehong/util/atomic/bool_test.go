package atomic

import (
	"sync/atomic"
	"testing"
)

func TestBool_Store_Load(t *testing.T) {
	var b Bool // zero value is false
	v := b.Load()
	if v {
		t.Errorf("want false, not true")
		return
	}

	b.Store(true)
	v = b.Load()
	if !v {
		t.Errorf("want true, not false")
		return
	}

	b.Store(false)
	v = b.Load()
	if v {
		t.Errorf("want false, not true")
		return
	}

	b.Store(true)
	v = b.Load()
	if !v {
		t.Errorf("want true, not false")
		return
	}
}

func TestBool_Swap(t *testing.T) {
	var val Bool

	val.Store(true)
	old := val.Swap(true)
	if !old {
		t.Errorf("want true, not false")
		return
	}
	if !val.Load() {
		t.Errorf("want true, not false")
		return
	}

	val.Store(true)
	old = val.Swap(false)
	if !old {
		t.Errorf("want true, not false")
		return
	}
	if val.Load() {
		t.Errorf("want false, not true")
		return
	}

	val.Store(false)
	old = val.Swap(true)
	if old {
		t.Errorf("want false, not true")
		return
	}
	if !val.Load() {
		t.Errorf("want true, not false")
		return
	}

	val.Store(false)
	old = val.Swap(false)
	if old {
		t.Errorf("want false, not true")
		return
	}
	if val.Load() {
		t.Errorf("want false, not true")
		return
	}
}

func TestBool_CompareAndSwap(t *testing.T) {
	var val Bool

	val.Store(true)
	swapped := val.CompareAndSwap(false, false)
	if swapped {
		t.Errorf("want false, not true")
		return
	}
	if !val.Load() {
		t.Errorf("want true, not false")
		return
	}

	val.Store(true)
	swapped = val.CompareAndSwap(true, false)
	if !swapped {
		t.Errorf("want true, not false")
		return
	}
	if val.Load() {
		t.Errorf("want false, not true")
		return
	}

	val.Store(true)
	swapped = val.CompareAndSwap(false, true)
	if swapped {
		t.Errorf("want false, not true")
		return
	}
	if !val.Load() {
		t.Errorf("want true, not false")
		return
	}

	val.Store(true)
	swapped = val.CompareAndSwap(true, true)
	if !swapped {
		t.Errorf("want true, not false")
		return
	}
	if !val.Load() {
		t.Errorf("want true, not false")
		return
	}

	val.Store(false)
	swapped = val.CompareAndSwap(false, false)
	if !swapped {
		t.Errorf("want true, not false")
		return
	}
	if val.Load() {
		t.Errorf("want false, not true")
		return
	}

	val.Store(false)
	swapped = val.CompareAndSwap(true, false)
	if swapped {
		t.Errorf("want false, not true")
		return
	}
	if val.Load() {
		t.Errorf("want false, not true")
		return
	}

	val.Store(false)
	swapped = val.CompareAndSwap(false, true)
	if !swapped {
		t.Errorf("want true, not false")
		return
	}
	if !val.Load() {
		t.Errorf("want true, not false")
		return
	}

	val.Store(false)
	swapped = val.CompareAndSwap(true, true)
	if swapped {
		t.Errorf("want false, not true")
		return
	}
	if val.Load() {
		t.Errorf("want false, not true")
		return
	}
}

func BenchmarkLoadBool(b *testing.B) {
	var val Bool
	var result bool

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = val.Load()
	}
	_ = result
}

func BenchmarkLoadAtomicValue(b *testing.B) {
	var val atomic.Value
	val.Store(false)
	var result bool

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = val.Load().(bool)
	}
	_ = result
}

func BenchmarkStoreBool(b *testing.B) {
	var val Bool

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.Store(true)
	}
}

func BenchmarkStoreAtomicValue(b *testing.B) {
	var val atomic.Value
	val.Store(false)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.Store(true)
	}
}
