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

package zproto

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/bytes2"
	"github.com/nebulaim/telegramd/baselib/net2"
	"hash/crc32"
	"io"
	"net"
	"time"
)

func init() {
	net2.RegisterProtocol("zproto", &ZProto{})
}

func generateMessageId(isReq bool) uint64 {
	const nano = 1000 * 1000 * 1000
	unixnano := time.Now().UnixNano()

	messageId := ((unixnano / nano) << 32) | ((unixnano % nano) & -4)
	for {
		if isReq {
			// rpc_request
			if (messageId % 4) != 3 {
				messageId += 1
			} else {
				break
			}
		} else {
			//rpc_response
			if (messageId % 4) != 1 {
				messageId += 1
			} else {
				break
			}
		}
	}

	return uint64(messageId)
}

type ZProto struct {
}

func (this *ZProto) NewCodec(rw io.ReadWriter) (net2.Codec, error) {
	codec := &ZProtoCodec{
		conn:                 rw.(*net.TCPConn),
		headBuf:              make([]byte, kFrameHeaderLen),
		recvLastPackageIndex: 0,
		sendLastPackageIndex: 0,
	}
	return codec, nil
}

type ZProtoCodec struct {
	conn                 *net.TCPConn
	headBuf              []byte
	recvLastPackageIndex uint32
	sendLastPackageIndex uint32
	nextSeqNo            uint32
	connID               uint64
}

func (c *ZProtoCodec) encodeMessage(x *bytes2.BufferOutput, m *ZProtoMessage) {
	m.sessionId = c.connID
	m.messageId = generateMessageId(false)
	m.seqNo = c.generateMessageSeqNo(true)
	x.UInt64(m.sessionId)
	x.UInt64(m.messageId)
	x.UInt32(m.seqNo)
	len := x.Len()
	x.UInt32(0)
	if m.Metadata != nil {
		m.Metadata.Encode(x)
		binary.LittleEndian.PutUint32(x.Buf()[len:], uint32(x.Len()-len))
	}
	m.Message.Encode(x)
}

func (c *ZProtoCodec) decodeMessage(b []byte) (*ZProtoMessage, error) {
	var err error

	dbuf := bytes2.NewBufferInput(b)

	m := &ZProtoMessage{}
	m.sessionId = dbuf.UInt64() // binary.LittleEndian.Uint64(payload[:8])
	m.messageId = dbuf.UInt64()
	m.seqNo = dbuf.UInt32() // binary.LittleEndian.Uint64(payload[8:16])

	// TODO(@benqi): check mdLen
	mdLen := dbuf.UInt32() // binary.LittleEndian.Uint32(payload[16:20])
	if mdLen > uint32(len(b)-28) {
		err := fmt.Errorf("metadata len invalid - mdLen: %d, bLen: %d", mdLen, len(b))
		glog.Error(err)
		return nil, err
	}
	//
	if mdLen > 0 {
		md := &ZProtoMetadata{}
		// b[20 : mdLen+20]
		err := md.Decode(dbuf)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		m.Metadata = md
	} else {
		//
	}

	m.Message, err = DecodeMessageByBuffer(dbuf)
	if err != nil {
		glog.Errorf("decode error - msg: {%v}, buf: {%s}", m, bytes2.DumpSize(512, b))
	}
	return m, err
}

func (c *ZProtoCodec) generateMessageSeqNo(increment bool) uint32 {
	value := c.nextSeqNo
	if increment {
		c.nextSeqNo++
		return uint32(value*2 + 1)
	} else {
		return uint32(value * 2)
	}
}

func (c *ZProtoCodec) Send(msg interface{}) error {
	switch msg.(type) {
	case *ZProtoMessage:
		message, _ := msg.(*ZProtoMessage)
		x := bytes2.NewBufferOutput(512)
		x.UInt32(kMagicNumber)
		x.UInt32(0)
		//.Add(1)
		x.UInt32(uint32(c.sendLastPackageIndex))
		c.sendLastPackageIndex++
		//x.UInt16(0)
		//x.UInt16(kVersion)

		// encode ZProtoMessage
		c.encodeMessage(x, message)

		binary.LittleEndian.PutUint32(x.Buf()[4:8], uint32(x.Len()+4))

		// write crc32
		crc32 := crc32.ChecksumIEEE(x.Buf())
		x.UInt32(crc32)

		_, err := c.conn.Write(x.Buf())
		return err
	default:
		return fmt.Errorf("invalid zproto message: %v", msg)
	}
}

func (c *ZProtoCodec) Receive() (interface{}, error) {
	_, err := io.ReadFull(c.conn, c.headBuf)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	// 1. check packageLength
	dbuf := bytes2.NewBufferInput(c.headBuf)
	// 1. check magic
	magic := dbuf.UInt32() // binary.LittleEndian.Uint32(c.headBuf[4:8])
	if magic != kMagicNumber {
		err = fmt.Errorf("invalid magic number: %d", magic)
		glog.Error(err)
		return nil, err
	}

	// TODO(@benqi): check packageLength
	packageLength := dbuf.UInt32() // binary.LittleEndian.Uint32(c.headBuf[:4])
	// glog.Infof("packageLength: %d", packageLength)

	// 2. check packageIndex
	// TODO(@benqi): check packageIndex
	packageIndex := dbuf.UInt32() // binary.LittleEndian.Uint32(c.headBuf[8:12])
	if packageIndex != 0 && packageIndex != c.recvLastPackageIndex+1 {
		err = fmt.Errorf("invalid packageIndex - lastPackageIndex: %d, packageIndex: %d", c.recvLastPackageIndex, packageIndex)
		glog.Error(err)
		return nil, err
	}
	c.recvLastPackageIndex = packageIndex

	payload := make([]byte, packageLength-kFrameHeaderLen) // recv crc32
	_, err = io.ReadFull(c.conn, payload)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	crcHash := crc32.NewIEEE()
	crcHash.Write(c.headBuf)
	crcHash.Write(payload[:len(payload)-4])

	// 3. check crc32
	crc := binary.LittleEndian.Uint32(payload[len(payload)-4:])

	if crc != crcHash.Sum32() {
		err = fmt.Errorf("invalid crc32: %d", crc)
		glog.Error(err)
		return nil, err
	}

	// glog.Info("onRawPayload \n", bytes2.DumpSize(256, payload))
	message, err := c.decodeMessage(payload)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return message, nil
}

func (c *ZProtoCodec) Close() error {
	return c.conn.Close()
}
