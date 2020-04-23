/*
 * Copyright (c) 2018-present, Yumcoder, LLC.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */
package app

import (
	"github.com/golang/glog"
	"sync"
	"testing"
)

type NullInstance struct {
	state int
	m     func()
}

func (e *NullInstance) Initialize() error {
	glog.Info("null instance initialize...")
	e.state = 1
	return nil
}

func (e *NullInstance) RunLoop() {
	glog.Info("null run_loop...")
	e.state = 2
	e.m()
}

func (e *NullInstance) Destroy() {
	glog.Info("null destroy...")
	e.state = 3
}

func TestRun1(t *testing.T) {
	instance := &NullInstance{}
	instance.m = func() {
		QuitAppInstance()
	}

	DoMainAppInstance(instance)

	result := instance.state
	expect := 3

	if result != expect {
		t.Error(`expect:`, expect, `result:`, result)
	}
}

func TestRun2(t *testing.T) {
	instance := &NullInstance{}
	wg := sync.WaitGroup{}

	wg.Add(1)
	instance.m = func() {
		glog.Info("done...")
		wg.Done()
	}

	go DoMainAppInstance(instance)
	wg.Wait()
	result := instance.state
	expect := 2

	if result != expect {
		t.Error(`expect:`, expect, `result:`, result)
	}
}
