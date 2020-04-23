// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package perf_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"unsafe"

	"golang.org/x/sys/unix"
	"golang.org/x/sys/unix/linux/perf"
)

func TestOpen(t *testing.T) {
	t.Run("BadGroup", testOpenBadGroup)
	t.Run("BadAttrType", testOpenBadAttrType)
	t.Run("PopulatesLabel", testOpenPopulatesLabel)
	t.Run("EventIDsDifferentByCPU", testEventIDsDifferentByCPU)
}

func testOpenBadGroup(t *testing.T) {
	requires(t, paranoid(1), hardwarePMU)

	ca := new(perf.Attr)
	perf.CPUCycles.Configure(ca)
	ca.CountFormat.Group = true

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cycles, err := perf.Open(ca, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	cycles.Close()

	_, err = perf.Open(ca, perf.CallingThread, perf.AnyCPU, cycles)
	if err == nil {
		t.Fatal("successful Open with closed group *Event")
	}

	cycles = new(perf.Event) // uninitialized
	_, err = perf.Open(ca, perf.CallingThread, perf.AnyCPU, cycles)
	if err == nil {
		t.Fatal("successful Open with closed group *Event")
	}
}

func testOpenBadAttrType(t *testing.T) {
	a := &perf.Attr{
		Type: 42,
	}

	_, err := perf.Open(a, perf.CallingThread, perf.AnyCPU, nil)
	if err == nil {
		t.Fatal("got a valid *Event for bad Attr.Type 42")
	}
}

func testOpenPopulatesLabel(t *testing.T) {
	// TODO(acln): extend when we implement general label lookup
	requires(t, paranoid(1), hardwarePMU)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ca := &perf.Attr{
		Type:   perf.HardwareEvent,
		Config: uint64(perf.CPUCycles),
	}

	cycles, err := perf.Open(ca, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cycles.Close()

	c, err := cycles.Measure(getpidTrigger)
	if err != nil {
		t.Fatal(err)
	}
	if c.Label == "" {
		t.Fatal("Open did not set label on *Attr")
	}
}

func testEventIDsDifferentByCPU(t *testing.T) {
	requires(t, paranoid(1), hardwarePMU)

	if runtime.NumCPU() == 1 {
		t.Skip("only one CPU")
	}

	ca := new(perf.Attr)
	perf.CPUCycles.Configure(ca)

	cycles0, err := perf.Open(ca, perf.CallingThread, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cycles0.Close()

	cycles1, err := perf.Open(ca, perf.CallingThread, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cycles1.Close()

	id0, err := cycles0.ID()
	if err != nil {
		t.Fatal(err)
	}

	id1, err := cycles1.ID()
	if err != nil {
		t.Fatal(err)
	}

	if id0 == id1 {
		t.Fatalf("event has the same ID on different CPUs")
	}
}

func TestMain(m *testing.M) {
	if !perf.Supported() {
		fmt.Fprintln(os.Stderr, "perf_event_open not supported")
		os.Exit(2)
	}
	os.Exit(m.Run())
}

// perfTestEnv holds and caches information about the testing environment
// for package perf.
type perfTestEnv struct {
	cap struct {
		sync.Once
		sysadmin bool
	}

	paranoid struct {
		sync.Once
		value int
	}

	tracefs struct {
		sync.Once
		mounted  bool
		readable bool
		readErr  error
	}

	pmu struct {
		sync.Mutex
		ok      map[string]struct{}
		missing map[string]error
	}
}

func (env *perfTestEnv) capSysAdmin() bool {
	env.cap.Once.Do(env.initCap)
	return env.cap.sysadmin
}

type capHeader struct {
	version uint32
	pid     int32
}

type capData struct {
	effective uint32
	_         uint32 // permitted
	_         uint32 // inheritable
}

// constants from uapi/linux/capability.h
const (
	capSysAdmin = 21
	capV3       = 0x20080522
)

func (env *perfTestEnv) initCap() {
	header := &capHeader{
		version: capV3,
		pid:     int32(unix.Getpid()),
	}
	data := make([]capData, 2)
	_, _, e := unix.Syscall(unix.SYS_CAPGET, uintptr(unsafe.Pointer(header)), uintptr(unsafe.Pointer(&data[0])), 0)
	if e != 0 {
		return
	}
	if data[0].effective&(1<<capSysAdmin) != 0 {
		env.cap.sysadmin = true
	}
}

func (env *perfTestEnv) paranoidLevel() int {
	env.paranoid.Once.Do(env.initParanoid)
	return env.paranoid.value
}

func (env *perfTestEnv) initParanoid() {
	content, err := ioutil.ReadFile("/proc/sys/kernel/perf_event_paranoid")
	if err != nil {
		env.paranoid.value = 3
		return
	}
	nr := strings.TrimSpace(string(content))
	paranoid, err := strconv.ParseInt(nr, 10, 32)
	if err != nil {
		env.paranoid.value = 3
		return
	}
	env.paranoid.value = int(paranoid)
}

func (env *perfTestEnv) initTracefs() {
	_, err := os.Stat("/sys/kernel/debug/tracing")
	if err != nil {
		return
	}
	env.tracefs.mounted = true
	_, err = ioutil.ReadDir("/sys/kernel/debug/tracing")
	if err != nil {
		env.tracefs.readErr = err
		return
	}
	env.tracefs.readable = true
}

func (env *perfTestEnv) tracefsMounted() bool {
	env.tracefs.Once.Do(env.initTracefs)
	return env.tracefs.mounted
}

func (env *perfTestEnv) tracefsReadable() (bool, error) {
	env.tracefs.Once.Do(env.initTracefs)
	return env.tracefs.readable, env.tracefs.readErr
}

func (env *perfTestEnv) havePMU(u string) (bool, error) {
	env.pmu.Lock()
	defer env.pmu.Unlock()

	if env.pmu.ok == nil {
		env.pmu.ok = map[string]struct{}{}
	}
	if env.pmu.missing == nil {
		env.pmu.missing = map[string]error{}
	}

	if _, ok := env.pmu.ok[u]; ok {
		return true, nil
	}
	if err, ok := env.pmu.missing[u]; ok {
		return false, err
	}

	_, err := perf.LookupEventType(u)
	if err != nil {
		env.pmu.missing[u] = err
		return false, err
	}

	env.pmu.ok[u] = struct{}{}
	return true, nil
}

var testenv perfTestEnv

// paranoid specifies a perf_event_paranoid level requirement for a test.
//
// For example, a value of 1 for paranoid means that the test requires a
// perf_event_paranoid level of 1 or less.
type paranoid int

func (p paranoid) Evaluate() error {
	if testenv.capSysAdmin() {
		return nil
	}
	want, have := int(p), testenv.paranoidLevel()
	if have > want {
		return fmt.Errorf("want perf_event_paranoid <= %d, have %d", want, have)
	}
	return nil
}

// tracefsreq specifies a tracefs requirement for a test: tracefs must be
// mounted at /sys/kernel/debug/tracing, and it must be readable.
type tracefsreq struct{}

func (tracefsreq) Evaluate() error {
	if !testenv.tracefsMounted() {
		return errors.New("tracefs is not mounted at /sys/kernel/debug/tracing")
	}
	if ok, err := testenv.tracefsReadable(); !ok {
		return fmt.Errorf("tracefs is not readable: %v", err)
	}
	return nil
}

var tracefs = tracefsreq{}

// pmu specifies a PMU requirement for a test.
type pmu string

var (
	hardwarePMU   = pmu("hardware")
	softwarePMU   = pmu("software")
	tracepointPMU = pmu("tracepoint")
)

func (u pmu) Evaluate() error {
	device := string(u)
	if device == "hardware" {
		device = "cpu" // TODO(acln): investigate
	}
	if ok, err := testenv.havePMU(device); !ok {
		return fmt.Errorf("%s PMU not supported: %v", device, err)
	}
	return nil
}

type testRequirement interface {
	Evaluate() error
}

func requires(t *testing.T, reqs ...testRequirement) {
	t.Helper()

	sb := new(strings.Builder)
	unmet := 0

	for _, req := range reqs {
		if err := req.Evaluate(); err != nil {
			if unmet > 0 {
				sb.WriteString("; ")
			}
			fmt.Fprint(sb, err)
			unmet++
		}
	}

	switch unmet {
	case 0:
		return
	case 1:
		t.Skipf("unmet requirement: %s", sb.String())
	default:
		t.Skipf("unmet requirements: %s", sb.String())
	}
}
