// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package perf_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
	"golang.org/x/sys/unix/linux/perf"
)

func TestPoll(t *testing.T) {
	t.Run("Timeout", testPollTimeout)
	t.Run("Cancel", testPollCancel)
	t.Run("Expired", testPollExpired)
	t.Run("DisabledExplicitly", testPollDisabledExplicitly)
	t.Run("DisabledByRefresh", testPollDisabledByRefresh)
	t.Run("DisabledByExit", testPollDisabledByExit)
}

func TestReadRecord(t *testing.T) {
	t.Run("Comm", testComm)
	t.Run("Exit", testExit)
	t.Run("CPUWideSwitch", testCPUWideSwitch)
	t.Run("SampleGetpid", testSampleGetpid)
	t.Run("SampleGetpidConcurrent", testSampleGetpidConcurrent)
	t.Run("SampleTracepointStack", testSampleTracepointStack)
	t.Run("RedirectedOutput", testRedirectedOutput)

	// TODO(acln): a test for the case when a record straddles the head
	// of the ring is missing. See readRawRecordNonblock.
}

func testPollTimeout(t *testing.T) {
	requires(t, paranoid(1), tracepointPMU, tracefs)

	ga := new(perf.Attr)
	ga.SetSamplePeriod(1)
	ga.SetWakeupEvents(1)
	gtp := perf.Tracepoint("syscalls", "sys_enter_getpid")
	if err := gtp.Configure(ga); err != nil {
		t.Fatal(err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	getpid, err := perf.Open(ga, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getpid.Close()
	if err := getpid.MapRing(); err != nil {
		t.Fatal(err)
	}

	errch := make(chan error)
	timeout := 20 * time.Millisecond

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		for i := 0; i < 2; i++ {
			_, err := getpid.ReadRecord(ctx)
			errch <- err
		}
	}()

	c, err := getpid.Measure(getpidTrigger)
	if err != nil {
		t.Fatal(err)
	}
	if c.Value != 1 {
		t.Fatalf("got %d hits for %q, want 1", c.Value, c.Label)
	}

	// For the first event, we should get a valid sample immediately.
	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatalf("didn't get the first sample: timeout")
	case err := <-errch:
		if err != nil {
			t.Fatalf("got %v, want valid first sample", err)
		}
	}

	// Now, we should get a timeout.
	select {
	case <-time.After(2 * timeout):
		t.Logf("didn't time out, waiting")
		err := <-errch
		t.Fatalf("got %v", err)
	case err := <-errch:
		if err != context.DeadlineExceeded {
			t.Fatalf("got %v, want context.DeadlineExceeded", err)
		}
	}
}

func testPollCancel(t *testing.T) {
	requires(t, paranoid(1), tracepointPMU, tracefs)

	ga := new(perf.Attr)
	ga.SetSamplePeriod(1)
	ga.SetWakeupEvents(1)
	gtp := perf.Tracepoint("syscalls", "sys_enter_getpid")
	if err := gtp.Configure(ga); err != nil {
		t.Fatal(err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	getpid, err := perf.Open(ga, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getpid.Close()
	if err := getpid.MapRing(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errch := make(chan error)

	go func() {
		for i := 0; i < 2; i++ {
			_, err := getpid.ReadRecord(ctx)
			errch <- err
		}
	}()

	c, err := getpid.Measure(getpidTrigger)
	if err != nil {
		t.Fatal(err)
	}
	if c.Value != 1 {
		t.Fatalf("got %d hits for %q, want 1", c.Value, c.Label)
	}

	// For the first event, we should get a valid sample.
	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatalf("didn't get the first sample: timeout")
	case err := <-errch:
		if err != nil {
			t.Fatalf("got %v, want valid first sample", err)
		}
	}

	// The goroutine reading the records is now blocked in ReadRecord.
	// Cancel the context and observe the results. We should see
	// context.Canceled quite quickly.
	cancel()

	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatalf("context cancel didn't unblock ReadRecord")
	case err := <-errch:
		if err != context.Canceled {
			t.Fatalf("got %v, want %v", err, context.Canceled)
		}
	}
}

func testPollExpired(t *testing.T) {
	requires(t, paranoid(1), softwarePMU)

	da := new(perf.Attr)
	perf.Dummy.Configure(da)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	dummy, err := perf.Open(da, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer dummy.Close()
	if err := dummy.MapRing(); err != nil {
		t.Fatal(err)
	}

	timeout := 1 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait until the deadline is in the past.
	time.Sleep(2 * timeout)

	rec, err := dummy.ReadRecord(ctx)
	if err == nil {
		t.Fatalf("got nil error and record %#v", rec)
	}
	if err != context.DeadlineExceeded {
		t.Fatalf("got %v, want context.DeadlineExceeded", err)
	}
}

const errDisabledTestEnv = "PERF_TEST_ERR_DISABLED"

func init() {
	// In child process of testErrDisabledProcessExist.
	if os.Getenv(errDisabledTestEnv) != "1" {
		return
	}

	readyevfd := 3
	startevfd := 4

	// Signal to the parent that we can start.
	evsig(readyevfd)

	// Wait for the parent to tell us that they have set up performance
	// monitoring, and are ready to observe the event.
	evwait(startevfd)

	// Call getpid, then exit. Parent will see POLLIN for getpid, then
	// POLLHUP because we exited.
	unix.Getpid()
	os.Exit(0)
}

func testPollDisabledByExit(t *testing.T) {
	requires(t, paranoid(1), tracepointPMU, tracefs)

	// Re-exec ourselves with PERF_TEST_ERR_DISABLED=1.
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	readyevfd, err := unix.Eventfd(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(readyevfd)

	startevfd, err := unix.Eventfd(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(startevfd)

	cmd := exec.Command(self)
	cmd.Env = append(os.Environ(), errDisabledTestEnv+"=1")
	cmd.ExtraFiles = []*os.File{
		os.NewFile(uintptr(readyevfd), "readyevfd"),
		os.NewFile(uintptr(startevfd), "startevfd"),
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Set up performance monitoring for the child process.
	ga := &perf.Attr{
		Options: perf.Options{
			Disabled: true,
		},
		SampleFormat: perf.SampleFormat{
			Tid: true,
		},
	}
	ga.SetSamplePeriod(1)
	ga.SetWakeupEvents(1)
	gtp := perf.Tracepoint("syscalls", "sys_enter_getpid")
	if err := gtp.Configure(ga); err != nil {
		t.Fatal(err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	getpid, err := perf.Open(ga, cmd.Process.Pid, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getpid.Close()
	if err := getpid.MapRing(); err != nil {
		t.Fatal(err)
	}

	// Wait for the child process to be ready.
	evwait(readyevfd)

	// Now that it is, enable the event.
	if err := getpid.Enable(); err != nil {
		t.Fatal(err)
	}

	// Signal to the child that it should call getpid now.
	// It will call getpid, then exit.
	evsig(startevfd)
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}

	// Read two records. The first one should be valid,
	// the second one should not, and the second error
	// should be ErrDisabled.
	timeout := 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rec1, err1 := getpid.ReadRecord(ctx)
	rec2, err2 := getpid.ReadRecord(ctx)

	if err1 != nil {
		t.Errorf("first error was %v, want nil", err1)
	}
	sr, ok := rec1.(*perf.SampleRecord)
	if !ok {
		t.Errorf("first record: got %T, want *perf.SampleRecord", rec1)
	}
	if int(sr.Pid) != cmd.Process.Pid {
		t.Errorf("first record: got pid %d in the sample, want %d",
			sr.Pid, cmd.Process.Pid)
	}

	if err2 != perf.ErrDisabled {
		t.Errorf("second record: error was %v, want ErrDisabled", err2)
	}
	if rec2 != nil {
		t.Errorf("second record: got %#v, want nil", rec2)
	}
}

func testPollDisabledExplicitly(t *testing.T) {
	requires(t, paranoid(1), tracepointPMU, tracefs)

	ga := &perf.Attr{
		SampleFormat: perf.SampleFormat{
			Tid: true,
		},
		Options: perf.Options{
			Disabled: true,
		},
	}
	ga.SetSamplePeriod(1)
	ga.SetWakeupEvents(1)
	gtp := perf.Tracepoint("syscalls", "sys_enter_getpid")
	if err := gtp.Configure(ga); err != nil {
		t.Fatal(err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	getpid, err := perf.Open(ga, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getpid.Close()
	if err := getpid.MapRing(); err != nil {
		t.Fatal(err)
	}

	const n = 3

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	seen := 0

	go func() {
		for i := 0; i < 2*n; i++ {
			_, err := getpid.ReadRecord(ctx)
			if err == nil {
				seen++
			}
		}
		close(done)
	}()

	if err := getpid.Enable(); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < n; i++ {
		getpidTrigger()
	}

	if err := getpid.Disable(); err != nil {
		getpidTrigger()
	}

	for i := 0; i < n; i++ {
		getpidTrigger()
	}

	cancel()
	<-done

	if seen != n {
		t.Fatalf("saw %d events, want %d", seen, n)
	}
}

func testPollDisabledByRefresh(t *testing.T) {
	// TODO(acln): investigate the following: the man page says that
	// POLLHUP should be indicated on the file descriptor when the counter
	// associated with a call to Refresh reaches zero.  I have not been
	// able to observe this. When the counter reaches zero, the event
	// is disabled (which is what this test shows), but POLLHUP doesn't
	// seem to be indicated on the file descriptor.
	//
	// If we ever figure out how to observe a HUP there, we should
	// make ReadRawRecord return ErrDisabled. In the meantime, leave
	// things as-is.
	requires(t, paranoid(1), tracepointPMU, tracefs)

	ga := &perf.Attr{
		SampleFormat: perf.SampleFormat{
			Tid: true,
		},
		Options: perf.Options{
			Disabled: true,
		},
	}
	ga.SetSamplePeriod(1)
	ga.SetWakeupEvents(1)
	gtp := perf.Tracepoint("syscalls", "sys_enter_getpid")
	if err := gtp.Configure(ga); err != nil {
		t.Fatal(err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	getpid, err := perf.Open(ga, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getpid.Close()
	if err := getpid.MapRing(); err != nil {
		t.Fatal(err)
	}

	const n = 3

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	seen := 0

	go func() {
		for i := 0; i < 2*n; i++ {
			_, err := getpid.ReadRecord(ctx)
			if err == nil {
				seen++
			}
		}
		close(done)
	}()

	if err := getpid.Refresh(n); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < n; i++ {
		getpidTrigger()
	}

	for i := 0; i < n; i++ {
		getpidTrigger()
	}

	cancel()
	<-done

	if seen != n {
		t.Fatalf("saw %d events, want %d", seen, n)
	}
}

const (
	commTestEnv  = "PERF_TEST_COMM"
	commTestName = "commtest"
)

func init() {
	// In child process of testComm.
	if os.Getenv(commTestEnv) != "1" {
		return
	}

	readyevfd := 3
	startevfd := 4
	sawcommevfd := 5

	// Signal to the parent that we can start.
	evsig(readyevfd)

	// Wait for the parent to tell us that they have set up performance
	// monitoring, and are ready to observe the event.
	evwait(startevfd)

	// Change our name.
	b := make([]byte, len(commTestName)+1)
	copy(b, commTestName)
	err := unix.Prctl(unix.PR_SET_NAME, uintptr(unsafe.Pointer(&b[0])), 0, 0, 0)
	runtime.KeepAlive(&b[0])
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(2)
	}

	// TODO(acln): investigate the legitimacy of the following crutch.
	//
	// Wait for the parent to see that we changed our name, then exit.
	//
	// If we do not wait here, there is a terrible race condition waiting
	// to happen: If we PR_SET_NAME in the child, then immediately exit,
	// the other side may not see POLLIN on the comm record: it may see
	// POLLHUP directly, even though a comm record was actually written
	// to the ring in the meantime. Why we get POLLHUP directly, and not
	// POLLIN before it, is unclear. The machinery to deal with this
	// eventuality in the poller does not exist yet, and at the time
	// when this comment was written, I have found no good solutions to
	// this conundrum.
	//
	// So we live with it, but still try to make our test pass.
	evwait(sawcommevfd)
	os.Exit(0)
}

func testComm(t *testing.T) {
	t.Skip("flaky. TODO(acln): investigate")

	requires(t, paranoid(1), softwarePMU)

	// Re-exec ourselves with PERF_TEST_COMM=1.
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	readyevfd, err := unix.Eventfd(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(readyevfd)

	startevfd, err := unix.Eventfd(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(startevfd)

	sawcommevfd, err := unix.Eventfd(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(sawcommevfd)

	cmd := exec.Command(self)
	cmd.Env = append(os.Environ(), commTestEnv+"=1")
	cmd.ExtraFiles = []*os.File{
		os.NewFile(uintptr(readyevfd), "readyevfd"),
		os.NewFile(uintptr(startevfd), "startevfd"),
		os.NewFile(uintptr(sawcommevfd), "sawcommevfd"),
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Set up performance monitoring for the child process.
	ca := &perf.Attr{
		Options: perf.Options{
			Disabled: true,
			Comm:     true,
		},
		SampleFormat: perf.SampleFormat{
			Tid: true,
		},
	}
	ca.SetSamplePeriod(1)
	ca.SetWakeupEvents(1)
	perf.Dummy.Configure(ca)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	comm, err := perf.Open(ca, cmd.Process.Pid, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer comm.Close()
	if err := comm.MapRing(); err != nil {
		t.Fatal(err)
	}

	// Wait for the child process to be ready.
	evwait(readyevfd)

	// Now that it is, enable the event.
	if err := comm.Enable(); err != nil {
		t.Fatal(err)
	}

	// Signal to the child that it should change its name.
	evsig(startevfd)

	// Read the CommRecord.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	rec, rerr := comm.ReadRecord(ctx)

	// Signal to the child that it should exit, and wait for it to do so.
	evsig(sawcommevfd)
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}

	// Observe the CommRecord.
	if rerr != nil {
		t.Fatalf("got %v, want valid record", rerr)
	}
	cr, ok := rec.(*perf.CommRecord)
	if !ok {
		t.Fatalf("got %T, want *perf.CommRecord", rec)
	}
	if int(cr.Pid) != cmd.Process.Pid {
		t.Errorf("got pid %d, want %d", cr.Pid, cmd.Process.Pid)
	}
	if cr.NewName != commTestName {
		t.Errorf("new name = %q, want %q", cr.NewName, commTestName)
	}
	if cr.WasExec() {
		t.Error("got WasExec() == true, want false")
	}
}

const (
	exitTestEnv  = "PERF_TEST_EXIT"
	exitTestCode = 42
)

func init() {
	// In the child process of testExit.
	if os.Getenv("PERF_TEST_EXIT") != "1" {
		return
	}

	readyevfd := 3
	startevfd := 4

	// Signal to the parent that we can start.
	evsig(readyevfd)

	// Wait for the parent to tell us that they have set up performance
	// monitoring, and are ready to observe the event.
	evwait(startevfd)

	os.Exit(exitTestCode)
}

func testExit(t *testing.T) {
	requires(t, paranoid(1), softwarePMU)

	// Re-exec ourselves with PERF_TEST_EXIT=1.
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	readyevfd, err := unix.Eventfd(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(readyevfd)

	startevfd, err := unix.Eventfd(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(startevfd)

	cmd := exec.Command(self)
	cmd.Env = append(os.Environ(), exitTestEnv+"=1")
	cmd.ExtraFiles = []*os.File{
		os.NewFile(uintptr(readyevfd), "readyevfd"),
		os.NewFile(uintptr(startevfd), "startevfd"),
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	pid := cmd.Process.Pid

	// Set up performance monitoring for the child process.
	ca := &perf.Attr{
		Options: perf.Options{
			Disabled: true,
			Task:     true,
		},
		SampleFormat: perf.SampleFormat{
			Tid: true,
		},
	}
	ca.SetSamplePeriod(1)
	ca.SetWakeupEvents(1)
	perf.Dummy.Configure(ca)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	comm, err := perf.Open(ca, pid, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer comm.Close()
	if err := comm.MapRing(); err != nil {
		t.Fatal(err)
	}

	// Wait for the child process to be ready.
	evwait(readyevfd)

	// Now that it is, enable the event.
	if err := comm.Enable(); err != nil {
		t.Fatal(err)
	}

	// Signal to the child that it should exit now.
	evsig(startevfd)

	// Observe the exit code from os/exec first.
	err = cmd.Wait()
	if err == nil {
		t.Fatal("child exited with code 0")
	}
	ee, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("got %T, want *exec.ExitError", err)
	}
	if got := ee.ExitCode(); got != exitTestCode {
		t.Fatalf("got exit code %d, want %d", got, exitTestCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	rec, err := comm.ReadRecord(ctx)
	if err != nil {
		t.Fatalf("got %v, want valid record", err)
	}
	er, ok := rec.(*perf.ExitRecord)
	if !ok {
		t.Fatalf("got %T, want *perf.ExitRecord", rec)
	}
	if int(er.Pid) != pid {
		t.Errorf("got pid %d, want %d", er.Pid, pid)
	}
	// Unfortunately, no er.Ppid and er.Ptid test. The Go runtime
	// interferes with us.
}

func testCPUWideSwitch(t *testing.T) {
	requires(t, paranoid(0), softwarePMU)

	var wg sync.WaitGroup
	ready := make(chan error)
	start := make(chan struct{})
	pingpong := make(chan struct{})
	var recvtid, sendtid int

	const numpingpongs = 4
	const cpu = 0

	fn := func(recv bool) {
		defer wg.Done()

		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		var cpuset unix.CPUSet
		cpuset.Set(cpu)
		if err := unix.SchedSetaffinity(0, &cpuset); err != nil {
			ready <- err
			return
		}

		if !recv {
			sendtid = unix.Gettid()
			ready <- nil
			<-start
			for i := 0; i < numpingpongs; i++ {
				pingpong <- struct{}{}
				<-pingpong
			}
		} else {
			recvtid = unix.Gettid()
			ready <- nil
			<-start
			for i := 0; i < numpingpongs; i++ {
				<-pingpong
				pingpong <- struct{}{}
			}
		}
	}

	wg.Add(2)

	go fn(true)
	go fn(false)

	if err := <-ready; err != nil {
		t.Fatal(err)
	}
	if err := <-ready; err != nil {
		t.Fatal(err)
	}

	sa := &perf.Attr{
		Options: perf.Options{
			ExcludeKernel: true,
			Disabled:      true,
			ContextSwitch: true,
		},
	}
	sa.SetSamplePeriod(1)
	sa.SetWakeupEvents(1)
	perf.ContextSwitches.Configure(sa)

	switches, err := perf.Open(sa, perf.AllThreads, cpu, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer switches.Close()
	if err := switches.MapRing(); err != nil {
		t.Fatal(err)
	}

	if err := switches.Enable(); err != nil {
		t.Fatal(err)
	}

	// Run the ping-pong game.
	close(start)
	wg.Wait()

	intorecv, outofrecv := 0, 0
	intosend, outofsend := 0, 0
	intosched, outofsched := 0, 0

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	var rerr error

	for {
		sawinto := intorecv >= numpingpongs && intosend >= numpingpongs
		sawoutof := outofrecv >= numpingpongs && outofsend >= numpingpongs
		if sawinto && sawoutof {
			break
		}
		rec, err := switches.ReadRecord(ctx)
		if err != nil {
			rerr = err
			break
		}
		sr, ok := rec.(*perf.SwitchCPUWideRecord)
		if !ok {
			t.Errorf("got %T, want *perf.SwitchCPUWideRecord", rec)
		}
		switch int(sr.Tid) {
		case 0:
			if sr.Out() {
				outofsched++
			} else {
				intosched++
			}
		case recvtid:
			if sr.Out() {
				outofrecv++
			} else {
				intorecv++
			}
		case sendtid:
			if sr.Out() {
				outofsend++
			} else {
				intosend++
			}
		}
	}

	if rerr != nil {
		t.Fatal(err)
	}

	t.Logf("%d ping-pongs", numpingpongs)
	t.Logf("recv switches: %d in, %d out", intorecv, outofrecv)
	t.Logf("send switches: %d in, %d out", intosend, outofsend)
	t.Logf("scheduler switches: %d in, %d out", intosched, outofsched)
}

func testSampleGetpid(t *testing.T) {
	requires(t, paranoid(1), tracepointPMU, tracefs)

	ga := &perf.Attr{
		SampleFormat: perf.SampleFormat{
			Tid: true,
		},
	}
	ga.SetSamplePeriod(1)
	ga.SetWakeupEvents(1)
	gtp := perf.Tracepoint("syscalls", "sys_enter_getpid")
	if err := gtp.Configure(ga); err != nil {
		t.Fatal(err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	getpid, err := perf.Open(ga, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getpid.Close()
	if err := getpid.MapRing(); err != nil {
		t.Fatal(err)
	}

	c, err := getpid.Measure(getpidTrigger)
	if err != nil {
		t.Fatal(err)
	}
	if c.Value != 1 {
		t.Fatalf("got %d hits for %q, want 1 hit", c.Value, c.Label)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	rec, err := getpid.ReadRecord(ctx)
	if err != nil {
		t.Fatalf("got %v, want a valid sample record", err)
	}
	sr, ok := rec.(*perf.SampleRecord)
	if !ok {
		t.Fatalf("got a %T, want a SampleRecord", rec)
	}
	pid, tid := unix.Getpid(), unix.Gettid()
	if int(sr.Pid) != pid || int(sr.Tid) != tid {
		t.Fatalf("got pid=%d tid=%d, want pid=%d tid=%d", sr.Pid, sr.Tid, pid, tid)
	}
}

func testSampleGetpidConcurrent(t *testing.T) {
	requires(t, paranoid(1), tracepointPMU, tracefs)

	ga := &perf.Attr{
		SampleFormat: perf.SampleFormat{
			Tid: true,
		},
	}
	ga.SetSamplePeriod(1)
	ga.SetWakeupEvents(1)
	gtp := perf.Tracepoint("syscalls", "sys_enter_getpid")
	if err := gtp.Configure(ga); err != nil {
		t.Fatal(err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	getpid, err := perf.Open(ga, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getpid.Close()
	if err := getpid.MapRing(); err != nil {
		t.Fatal(err)
	}

	const n = 6
	sawSample := make(chan bool)

	go func() {
		for i := 0; i < n; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			rec, err := getpid.ReadRecord(ctx)
			_, isSample := rec.(*perf.SampleRecord)
			if err == nil && isSample {
				sawSample <- true
			} else {
				sawSample <- false
			}
		}
	}()

	seen := 0

	c, err := getpid.Measure(func() {
		for i := 0; i < n; i++ {
			getpidTrigger()
			if ok := <-sawSample; ok {
				seen++
			}
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Value != n {
		t.Fatalf("got %d hits for %q, want %d", c.Value, c.Label, n)
	}
	if seen != n {
		t.Fatalf("saw %d samples, want %d", seen, n)
	}
}

func testSampleTracepointStack(t *testing.T) {
	requires(t, paranoid(1), tracepointPMU, tracefs)

	ga := &perf.Attr{
		Options: perf.Options{
			Disabled: true,
		},
		SampleFormat: perf.SampleFormat{
			Tid:       true,
			Time:      true,
			CPU:       true,
			IP:        true,
			Callchain: true,
		},
	}
	ga.SetSamplePeriod(1)
	ga.SetWakeupEvents(1)
	gtp := perf.Tracepoint("syscalls", "sys_enter_getpid")
	if err := gtp.Configure(ga); err != nil {
		t.Fatal(err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	getpid, err := perf.Open(ga, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getpid.Close()
	if err := getpid.MapRing(); err != nil {
		t.Fatal(err)
	}

	pcs := make([]uintptr, 10)
	var n int

	c, err := getpid.Measure(func() {
		n = runtime.Callers(2, pcs)
		getpidTrigger()
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Value != 1 {
		t.Fatalf("want 1 hit for %q, got %d", c.Label, c.Value)
	}

	pcs = pcs[:n]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	rec, err := getpid.ReadRecord(ctx)
	if err != nil {
		t.Fatal(err)
	}
	getpidsample, ok := rec.(*perf.SampleRecord)
	if !ok {
		t.Fatalf("got a %T, want a *SampleRecord", rec)
	}

	i := len(pcs) - 1
	j := len(getpidsample.Callchain) - 1

	for i >= 0 && j >= 0 {
		gopc := pcs[i]
		kpc := getpidsample.Callchain[j]
		if gopc != uintptr(kpc) {
			t.Fatalf("Go (%#x) and kernel (%#x) PC differ", gopc, kpc)
		}
		i--
		j--
	}

	logFrame := func(pc uintptr) {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			t.Logf("%#x <nil>", pc)
		} else {
			file, line := fn.FileLine(pc)
			t.Logf("%#x %s:%d %s", pc, file, line, fn.Name())
		}
	}

	t.Log("kernel callchain:")
	for _, kpc := range getpidsample.Callchain {
		logFrame(uintptr(kpc))
	}

	t.Log()

	t.Logf("Go stack:")
	for _, gopc := range pcs {
		logFrame(gopc)
	}
}

func testRedirectedOutput(t *testing.T) {
	requires(t, paranoid(1), tracepointPMU, tracefs)

	ga := &perf.Attr{
		SampleFormat: perf.SampleFormat{
			Tid:      true,
			Time:     true,
			CPU:      true,
			Addr:     true,
			StreamID: true,
		},
		CountFormat: perf.CountFormat{
			Group: true,
		},
		Options: perf.Options{
			Disabled: true,
		},
	}
	ga.SetSamplePeriod(1)
	ga.SetWakeupEvents(1)
	gtp := perf.Tracepoint("syscalls", "sys_enter_getpid")
	if err := gtp.Configure(ga); err != nil {
		t.Fatalf("Configure: %v", err)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	leader, err := perf.Open(ga, perf.CallingThread, perf.AnyCPU, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer leader.Close()
	if err := leader.MapRing(); err != nil {
		t.Fatal(err)
	}

	wa := &perf.Attr{
		SampleFormat: perf.SampleFormat{
			Tid:      true,
			Time:     true,
			CPU:      true,
			Addr:     true,
			StreamID: true,
		},
	}
	wa.SetSamplePeriod(1)
	wa.SetWakeupEvents(1)
	wtp := perf.Tracepoint("syscalls", "sys_enter_write")
	if err := wtp.Configure(wa); err != nil {
		t.Fatal(err)
	}

	follower, err := perf.Open(wa, perf.CallingThread, perf.AnyCPU, leader)
	if err != nil {
		t.Fatal(err)
	}
	defer follower.Close()
	if err := follower.SetOutput(leader); err != nil {
		t.Fatal(err)
	}

	errch := make(chan error)
	go func() {
		for i := 0; i < 2; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			_, err := leader.ReadRecord(ctx)
			errch <- err
		}
	}()

	gc, err := leader.MeasureGroup(func() {
		getpidTrigger()
		writeTrigger()
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := gc.Values[0]; got.Value != 1 {
		t.Fatalf("got %d hits for %q, want 1 hit", got.Value, got.Label)
	}
	if got := gc.Values[1]; got.Value != 1 {
		t.Fatalf("got %d hits for %q, want 1 hit", got.Value, got.Label)
	}

	for i := 0; i < 2; i++ {
		select {
		case <-time.After(10 * time.Millisecond):
			t.Errorf("did not get sample record: timeout")
		case err := <-errch:
			if err != nil {
				t.Fatalf("did not get sample record: %v", err)
			}
		}
	}
}

func evsig(fd int) {
	val := uint64(1)
	buf := (*[8]byte)(unsafe.Pointer(&val))[:]
	unix.Write(fd, buf)
}

func evwait(fd int) {
	var val uint64
	buf := (*[8]byte)(unsafe.Pointer(&val))[:]
	unix.Read(fd, buf)
}
