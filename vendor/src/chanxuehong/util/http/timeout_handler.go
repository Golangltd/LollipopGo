// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

type Skipper func(*http.Request) bool

// A BufferPool is an interface for getting and returning temporary bytes.Buffer.
type BufferPool interface {
	Get() *bytes.Buffer
	Put(*bytes.Buffer)
}

// TimeoutHandler returns a http.Handler that runs h with the given time limit.
//
// The new http.Handler calls h.ServeHTTP to handle each request, but if a
// call runs for longer than its time limit, the handler responds with
// the given http status code and the given message in its body.
// (If code is 0, http.StatusServiceUnavailable will be sent; If msg is empty, a suitable default message will be sent.)
// pool can be nil, and if it is nil, we use new(bytes.Buffer) to get bytes.Buffer every times.
// After such a timeout, writes by h to its http.ResponseWriter will return
// http.ErrHandlerTimeout.
//
// TimeoutHandler buffers all http.Handler writes to memory and does not
// support the http.Hijacker or http.Flusher interfaces.
func TimeoutHandler(h http.Handler, dt time.Duration, code int, msg string, pool BufferPool, skipper Skipper) http.Handler {
	return &timeoutHandler{
		handler: h,
		code:    code,
		body:    msg,
		dt:      dt,
		pool:    pool,
		skipper: skipper,
	}
}

type timeoutHandler struct {
	handler http.Handler
	code    int
	body    string
	dt      time.Duration
	pool    BufferPool
	skipper Skipper

	// When set, no context will be created and this context will
	// be used instead.
	testContext context.Context
}

func (h *timeoutHandler) errorCode() int {
	if h.code > 0 {
		return h.code
	}
	return http.StatusServiceUnavailable
}

func (h *timeoutHandler) errorBody() string {
	if h.body != "" {
		return h.body
	}
	return "<html><head><title>Timeout</title></head><body><h1>Timeout</h1></body></html>"
}

func (h *timeoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.skipper != nil && h.skipper(r) {
		h.handler.ServeHTTP(w, r)
		return
	}
	ctx := h.testContext
	if ctx == nil {
		var cancelCtx context.CancelFunc
		ctx, cancelCtx = context.WithTimeout(r.Context(), h.dt)
		defer cancelCtx()
	}
	r = r.WithContext(ctx)
	done := make(chan struct{})
	var wbuf *bytes.Buffer
	if h.pool != nil {
		wbuf = h.pool.Get()
		wbuf.Reset()
		defer h.pool.Put(wbuf)
	} else {
		wbuf = new(bytes.Buffer)
	}
	tw := &timeoutWriter{
		w:    w,
		h:    make(http.Header),
		wbuf: wbuf,
	}
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		h.handler.ServeHTTP(tw, r)
		close(done)
	}()
	select {
	case p := <-panicChan:
		panic(p)
	case <-done:
		tw.mu.Lock()
		defer tw.mu.Unlock()
		dst := w.Header()
		for k, vv := range tw.h {
			dst[k] = vv
		}
		if !tw.wroteHeader {
			tw.code = http.StatusOK
		}
		w.WriteHeader(tw.code)
		w.Write(tw.wbuf.Bytes())
		tw.handlerDone = true
		if tw.closeNotifyCh != nil {
			tw.closeNotifyCh <- true
		}
		return
	case <-ctx.Done():
		tw.mu.Lock()
		defer tw.mu.Unlock()
		w.WriteHeader(h.errorCode())
		io.WriteString(w, h.errorBody())
		tw.timedOut = true
		if tw.closeNotifyCh != nil {
			tw.closeNotifyCh <- true
		}
		return
	}
}

type timeoutWriter struct {
	w    http.ResponseWriter
	h    http.Header
	wbuf *bytes.Buffer

	closeNotifyCh chan bool

	mu          sync.Mutex
	timedOut    bool
	handlerDone bool
	wroteHeader bool
	code        int
}

func (tw *timeoutWriter) Header() http.Header { return tw.h }

func (tw *timeoutWriter) Write(p []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, http.ErrHandlerTimeout
	}
	if tw.handlerDone {
		return 0, http.ErrContentLength
	}
	if !tw.wroteHeader {
		tw.writeHeader(http.StatusOK)
	}
	return tw.wbuf.Write(p)
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader || tw.handlerDone {
		return
	}
	tw.writeHeader(code)
}

func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

func (tw *timeoutWriter) WriteString(s string) (n int, err error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, http.ErrHandlerTimeout
	}
	if tw.handlerDone {
		return 0, http.ErrContentLength
	}
	if !tw.wroteHeader {
		tw.writeHeader(http.StatusOK)
	}
	return tw.wbuf.WriteString(s)
}

func (tw *timeoutWriter) CloseNotify() <-chan bool {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if v, ok := tw.w.(http.CloseNotifier); ok {
		return v.CloseNotify()
	}
	if tw.closeNotifyCh == nil {
		tw.closeNotifyCh = make(chan bool, 1)
	}
	return tw.closeNotifyCh
}
