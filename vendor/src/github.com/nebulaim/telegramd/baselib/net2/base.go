/*
 *  Copyright (c) 2017, https://github.com/nebulaim
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package net2

import (
	"fmt"
	"io"
)

type Protocol interface {
	NewCodec(rw io.ReadWriter) (Codec, error)
}

type ProtocolFunc func(rw io.ReadWriter) (Codec, error)

func (pf ProtocolFunc) NewCodec(rw io.ReadWriter) (Codec, error) {
	return pf(rw)
}

type Codec interface {
	Receive() (interface{}, error)
	Send(interface{}) error
	Close() error
}

type MessageBase interface {
	Encode() []byte
	Decode(b []byte) error
}

type ConnectionFactory interface {
	NewConnection(serverName string) TcpConnection
}

type ClearSendChan interface {
	ClearSendChan(<-chan interface{})
}

var (
	protocolRegisters = make(map[string]Protocol)
)

func RegisterProtocol(name string, protocol Protocol) {
	protocolRegisters[name] = protocol
}

func NewCodecByName(name string, rw io.ReadWriter) (Codec, error) {
	protocol, ok := protocolRegisters[name]
	if !ok {
		return nil, fmt.Errorf("not found protocol name: %s", name)
	}
	return protocol.NewCodec(rw)
}
