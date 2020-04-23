// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package perf_test

import (
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"golang.org/x/sys/unix"
	"golang.org/x/sys/unix/linux/perf"
)

func TestCount(t *testing.T) {
	t.Run("Hardware", testHardwareCounters)
	t.Run("Software", testSoftwareCounters)
	t.Run("HardwareCache", testHardwareCacheCounters)
	t.Run("Tracepoint", testSingleTracepoint)
	t.Run("IoctlAndCountIDsMatch", testIoctlAndCountIDsMatch)
}

func testHardwareCounters(t *testing.T) {
	requires(t, paranoid(1), hardwarePMU)

	t.Run("IPC", testIPC)
}

func testIPC(t *testing.T) {
	g := perf.Group{
		CountFormat: perf.CountFormat{
			ID: true,
		},
	}
	g.Add(perf.Instructions, perf.CPUCycles)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	hw, err := g.Open(perf.CallingThread, perf.AnyCPU)
	if err != nil {
		t.Fatal(err)
	}
	defer hw.Close()

	var sum int64
	gc, err := hw.MeasureGroup(func() {
		for i := int64(0); i < 1000000; i++ {
			sum += i
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range gc.Values {
		if c.Value == 0 {
			t.Fatalf("didn't count %q", c.Label)
		}
	}
	insns := gc.Values[0].Value
	cycles := gc.Values[1].Value
	ipc := float64(insns) / float64(cycles)
	t.Logf("got %d instructions, %d cycles: %f IPC", insns, cycles, ipc)
}

func testSoftwareCounters(t *testing.T) {
	requires(t, paranoid(1), softwarePMU)

	t.Run("PageFaults", testPageFaults)
}

var fault []byte

func testPageFaults(t *testing.T) {
	// TODO(acln): this test starts failing when run with -count > 4-5,
	// even though we're calling debug.FreeOSMemory. Why?
	pfa := &perf.Attr{
		CountFormat: perf.CountFormat{
			Running: true,
			Enabled: true,
		},
	}
	perf.PageFaults.Configure(pfa)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	faults, err := perf.Open(pfa, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer faults.Close()

	debug.FreeOSMemory()

	c, err := faults.Measure(func() {
		fault = make([]byte, 64*1024*1024)
		fault[0] = 1
		fault[63*1024*1024] = 1
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Value == 0 {
		t.Fatal("didn't see a page fault")
	}
	t.Logf("saw %v: enabled: %v, running: %v", c, c.Enabled, c.Running)
}

func testHardwareCacheCounters(t *testing.T) {
	// TODO(acln): add PMU requirement? but how?
	//
	// $ ls /sys/bus/event_source/devices/*/type | xargs cat
	//
	// does not contain a 3, which is the ABI-specified value of
	// perf.HardwareCacheEvent. Maybe it's under the "cpu" PMU
	// somewhere. Investigate.
	requires(t, paranoid(1))

	t.Run("L1DataMissesBadLocality", testL1DataMissesBadLocality)
	t.Run("L1DataMissesGoodLocality", testL1DataMissesGoodLocality)
	t.Run("L1Group", testL1Group)
}

func testL1DataMissesBadLocality(t *testing.T) {
	hwca := new(perf.Attr)
	hwcc := perf.HardwareCacheCounter{
		Cache:  perf.L1D,
		Op:     perf.Read,
		Result: perf.Miss,
	}
	hwcc.Configure(hwca)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	l1dmisses, err := perf.Open(hwca, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer l1dmisses.Close()

	rng := rand.New(rand.NewSource(time.Now().Unix()))

	max := 1000

	var bad []interface{}
	for i := 0; i < 10000; i++ {
		bad = append(bad, rng.Intn(max))
	}

	sink := 0
	c, err := l1dmisses.Measure(func() {
		for _, v := range bad {
			if v.(int) < max/2 {
				sink++
			}
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Value == 0 {
		t.Fatalf("recorded no L1 data cache misses")
	}

	t.Logf("bad locality: got %d L1 data cache misses", c.Value)
}

func testL1DataMissesGoodLocality(t *testing.T) {
	hwca := new(perf.Attr)
	hwcc := perf.HardwareCacheCounter{
		Cache:  perf.L1D,
		Op:     perf.Read,
		Result: perf.Miss,
	}
	hwcc.Configure(hwca)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	l1dmisses, err := perf.Open(hwca, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer l1dmisses.Close()

	rng := rand.New(rand.NewSource(time.Now().Unix()))

	max := 1000

	var contiguous []int
	for i := 0; i < 10000; i++ {
		contiguous = append(contiguous, rng.Intn(max))
	}

	sink := 0
	c, err := l1dmisses.Measure(func() {
		for _, v := range contiguous {
			if v < max/2 {
				sink++
			}
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Value == 0 {
		t.Fatalf("recorded no L1 data cache misses")
	}

	t.Logf("good locality: got %d L1 data cache misses", c.Value)
}

type l1testIdentity int

func (v l1testIdentity) value() int { return int(v) }

type l1testSquare int

func (v l1testSquare) value() int { return int(v * v) }

type l1testCube int

func (v l1testCube) value() int { return int(v * v * v) }

type valuer interface {
	value() int
}

func newValuer(n int) valuer {
	switch n % 3 {
	case 0:
		return l1testIdentity(n)
	case 1:
		return l1testSquare(n)
	default:
		return l1testCube(n)
	}
}

func testL1Group(t *testing.T) {
	caches := []perf.Cache{perf.L1D, perf.L1I}
	ops := []perf.CacheOp{perf.Read}
	results := []perf.CacheOpResult{perf.Miss}

	var g perf.Group
	g.Add(perf.HardwareCacheCounters(caches, ops, results)...)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	l1, err := g.Open(perf.CallingThread, perf.AnyCPU)
	if err != nil {
		t.Fatal(err)
	}
	defer l1.Close()

	const n = 100000

	valuers := make([]valuer, 0, n)
	for i := 0; i < n; i++ {
		valuers = append(valuers, newValuer(i))
	}

	sum := 0
	gc, err := l1.MeasureGroup(func() {
		for i := 0; i < n; i++ {
			sum += valuers[i].value()
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("got %d L1 data cache misses, %d L1 instruction cache misses",
		gc.Values[0].Value, gc.Values[1].Value)
}

func testSingleTracepoint(t *testing.T) {
	requires(t, paranoid(1), tracepointPMU, tracefs)

	tests := []singleTracepointTest{
		{
			category: "syscalls",
			event:    "sys_enter_getpid",
			trigger:  getpidTrigger,
		},
		{
			category: "syscalls",
			event:    "sys_enter_read",
			trigger:  readTrigger,
		},
		{
			category: "syscalls",
			event:    "sys_enter_write",
			trigger:  writeTrigger,
		},
	}
	for _, tt := range tests {
		t.Run(tt.String(), tt.run)
	}
}

type singleTracepointTest struct {
	category string
	event    string
	trigger  func()
}

func (tt singleTracepointTest) run(t *testing.T) {
	tp := perf.Tracepoint(tt.category, tt.event)
	attr := new(perf.Attr)
	if err := tp.Configure(attr); err != nil {
		t.Fatal(err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ev, err := perf.Open(attr, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ev.Close()

	c, err := ev.Measure(func() {
		tt.trigger()
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Value != 1 {
		t.Fatalf("got %d hits for %q, want 1 hit", c.Value, c.Label)
	}
}

func (tt singleTracepointTest) String() string {
	return tt.category + ":" + tt.event
}

func testIoctlAndCountIDsMatch(t *testing.T) {
	requires(t, paranoid(1), softwarePMU)

	pfa := new(perf.Attr)
	perf.PageFaults.Configure(pfa)
	pfa.CountFormat.ID = true

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	faults, err := perf.Open(pfa, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}

	runtime.GC()

	c, err := faults.Measure(func() {
		fault = make([]byte, 64*1024*1024)
		fault[0] = 1
		fault[63*1024*1024] = 1
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Value == 0 {
		t.Fatal("didn't see a page fault")
	}
	id, err := faults.ID()
	if err != nil {
		t.Fatal(err)
	}
	if id != c.ID {
		t.Fatalf("got ID %d from ioctl, but %d from count read", id, c.ID)
	}
}

func getpidTrigger() {
	unix.Getpid()
}

func readTrigger() {
	zero, err := os.Open("/dev/zero")
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 8)
	if _, err := zero.Read(buf); err != nil {
		panic(err)
	}
}

func writeTrigger() {
	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0200)
	if err != nil {
		panic(err)
	}
	if _, err := null.Write([]byte("big data")); err != nil {
		panic(err)
	}
}
