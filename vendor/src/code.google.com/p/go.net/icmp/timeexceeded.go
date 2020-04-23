// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

// A TimeExceeded represents an ICMP time exceeded message body.
type TimeExceeded struct {
	Data []byte // data
}

// Len implements the Len method of MessageBody interface.
func (p *TimeExceeded) Len() int {
	if p == nil {
		return 0
	}
	return 4 + len(p.Data)
}

// Marshal implements the Marshal method of MessageBody interface.
func (p *TimeExceeded) Marshal() ([]byte, error) {
	b := make([]byte, 4+len(p.Data))
	copy(b[4:], p.Data)
	return b, nil
}

// parseTimeExceeded parses b as an ICMP time exceeded message body.
func parseTimeExceeded(b []byte) (MessageBody, error) {
	bodyLen := len(b)
	if bodyLen < 4 {
		return nil, errMessageTooShort
	}
	p := &TimeExceeded{}
	if bodyLen > 4 {
		p.Data = make([]byte, bodyLen-4)
		copy(p.Data, b[4:])
	}
	return p, nil
}
