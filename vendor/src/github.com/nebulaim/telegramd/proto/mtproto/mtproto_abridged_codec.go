/*
 *  Copyright (c) 2018, https://github.com/nebulaim
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

package mtproto

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/net2"
	"io"
)

// https://core.telegram.org/mtproto#tcp-transport
//
// There is an abridged version of the same protocol:
// if the client sends 0xef as the first byte (**important:** only prior to the very first data packet),
// then packet length is encoded by a single byte (0x01..0x7e = data length divided by 4;
// or 0x7f followed by 3 length bytes (little endian) divided by 4) followed
// by the data themselves (sequence number and CRC32 not added).
// In this case, server responses look the same (the server does not send 0xefas the first byte).
//
type MTProtoAbridgedCodec struct {
	conn *net2.BufferedConn
}

func NewMTProtoAbridgedCodec(conn *net2.BufferedConn) *MTProtoAbridgedCodec {
	return &MTProtoAbridgedCodec{
		conn: conn,
	}
}

func (c *MTProtoAbridgedCodec) Receive() (interface{}, error) {
	// minus padding
	//size := len(x.buf) / 4 - 1
	//
	//if size < 127 {
	//	x.buf[3] = byte(size)
	//	x.buf = x.buf[3:]
	//} else {
	//	binary.LittleEndian.PutUint32(x.buf, uint32(size << 8 | 127))
	//}
	//_, err := m.conn.Write(x.buf)
	//if err != nil {
	//	return err
	//}

	var size int
	var n int
	var err error

	b := make([]byte, 1)
	n, err = io.ReadFull(c.conn, b)
	if err != nil {
		return nil, err
	}

	// glog.Info("first_byte: ", hex.EncodeToString(b[:1]))
	needAck := bool(b[0]>>7 == 1)
	_ = needAck

	b[0] = b[0] & 0x7f
	// glog.Info("first_byte2: ", hex.EncodeToString(b[:1]))

	if b[0] < 0x7f {
		size = int(b[0]) << 2
		glog.Info("size1: ", size)
		if size == 0 {
			return nil, nil
		}
	} else {
		glog.Info("first_byte2: ", hex.EncodeToString(b[:1]))
		b2 := make([]byte, 3)
		n, err = io.ReadFull(c.conn, b2)
		if err != nil {
			return nil, err
		}
		size = (int(b2[0]) | int(b2[1])<<8 | int(b2[2])<<16) << 2
		glog.Info("size2: ", size)
	}

	left := size
	buf := make([]byte, size)
	for left > 0 {
		n, err = io.ReadFull(c.conn, buf[size-left:])
		if err != nil {
			glog.Error("ReadFull2 error: ", err)
			return nil, err
		}
		left -= n
	}
	if size > 10240 {
		glog.Info("ReadFull2: ", hex.EncodeToString(buf[:256]))
	}

	// TODO(@benqi): process report ack and quickack
	// 截断QuickAck消息，客户端有问题
	if size == 4 {
		glog.Errorf("Server response error: ", int32(binary.LittleEndian.Uint32(buf)))
		// return nil, fmt.Errorf("Recv QuickAckMessage, ignore!!!!") //  connId: ", c.stream, ", by client ", m.RemoteAddr())
		return nil, nil
	}

	authKeyId := int64(binary.LittleEndian.Uint64(buf))
	message := NewMTPRawMessage(authKeyId, 0, TRANSPORT_TCP)
	message.Decode(buf)
	return message, nil
}

func (c *MTProtoAbridgedCodec) Send(msg interface{}) error {
	message, ok := msg.(*MTPRawMessage)
	if !ok {
		err := fmt.Errorf("msg type error, only MTPRawMessage, msg: {%v}", msg)
		glog.Error(err)
		return err
	}

	b := message.Encode()

	sb := make([]byte, 4)
	// minus padding
	size := len(b) / 4

	if size < 127 {
		sb = []byte{byte(size)}
	} else {
		binary.LittleEndian.PutUint32(sb, uint32(size<<8|127))
	}

	b = append(sb, b...)
	_, err := c.conn.Write(b)

	if err != nil {
		glog.Errorf("Send msg error: %s", err)
	}

	return err
}

func (c *MTProtoAbridgedCodec) Close() error {
	return c.conn.Close()
}
