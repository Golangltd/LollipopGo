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
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/bytes2"
)

const (
	PROTO              = 0xFF00
	PING               = 0x0100
	PONG               = 0x0200
	DROP               = 0x0300
	REDIRECT           = 0x0400
	ACK                = 0x0500
	HANDSHAKE_REQ      = 0x0600
	HANDSHAKE_RSP      = 0x0700
	MARS_SIGNAL        = 0x0800
	MESSAGE_ACK        = 0x0001
	RPC_REQUEST        = 0x000F
	RPC_OK             = 0x0010
	RPC_ERROR          = 0x0011
	RPC_FLOOD_WAIT     = 0x0012
	RPC_INTERNAL_ERROR = 0x0013
	PUSH               = 0x0014
)

const (
	kFrameHeaderLen = 12
	kMagicNumber    = 0x5A4D5450 // "ZMTP"
	kVersion        = 1
)

type MessageBase interface {
	Encode(x *bytes2.BufferOutput)
	Decode(dbuf *bytes2.BufferInput) error
}

type newZProtoMessage func() MessageBase

var zprotoFactories = make(map[uint32]newZProtoMessage)

//func CheckPackageType(packageType uint32) (r bool) {
//	_, r = zprotoFactories[packageType]
//	return
//}
//

func NewZProtoMessage(msgType uint32) MessageBase {
	m, ok := zprotoFactories[msgType]
	if !ok {
		glog.Errorf("invalid msgType: %d", msgType)
		return nil
	}
	return m()
}

/*
type ZProtoPackageData struct {
	magicNumber   uint32 // "ZMTP"
	packageLength uint32 // 整个数据包长度 packageLen + bodyLen: metaDataLength + len(metaData) + len(payload) + len(crc32)
	packageIndex  uint32 // Index of package starting from zero. If packageIndex is broken connection need to be dropped.
	// version       uint16 // version
	// reserved      uint16 // reserved
	sessionId     uint64
	seqNum        uint64
	metadata      []byte
	packageType   uint32 // frameType
	body          []byte
	crc32         uint32 // CRC32 of body
}
*/

type ZProtoMessage struct {
	sessionId uint64
	messageId uint64
	seqNo     uint32
	Metadata  *ZProtoMetadata
	Message   MessageBase
}

func (m *ZProtoMessage) Encode(x *bytes2.BufferOutput) {
	// return nil
}

func (m *ZProtoMessage) Decode(dbuf *bytes2.BufferInput) error {
	return nil
}

/*
type ZProtoMessageData struct {
	SessionId uint64
	SeqNum    uint64
	Metadata  *ZProtoMetadata
	Message   net2.MessageBase
}

func (m *ZProtoMessageData) Encode() ([]byte) {
	return nil
}

func (m *ZProtoMessageData) Decode(b []byte) error {
	return nil
}
*/

/////////////////////////////////////////////////////////////////
type ZProtoMetadata struct {
	ServerId     int
	ClientConnId uint64
	ClientAddr   string
	TraceId      int64
	SpanId       int64
	ReceiveTime  int64
	From         string
	Options      map[string]string
	extend       []byte
}

func (m *ZProtoMetadata) Encode(x *bytes2.BufferOutput) {
	x.UInt32(uint32(m.ServerId))
	x.UInt64(uint64(m.ClientConnId))
	bytes2.WriteString(x, m.ClientAddr)
	x.UInt64(uint64(m.TraceId))
	x.UInt64(uint64(m.SpanId))
	x.UInt64(uint64(m.ReceiveTime))
	bytes2.WriteString(x, m.From)
	x.UInt32(uint32(len(m.Options)))
	for k, v := range m.Options {
		bytes2.WriteString(x, k)
		bytes2.WriteString(x, v)
	}
	bytes2.WriteBytes(x, m.extend)
}

func (m *ZProtoMetadata) Decode(dbuf *bytes2.BufferInput) (err error) {
	m.ServerId = int(dbuf.UInt32())
	m.ClientConnId = dbuf.UInt64()
	m.ClientAddr, _ = bytes2.ReadString(dbuf)
	m.TraceId = dbuf.Int64()
	m.SpanId = dbuf.Int64()
	m.ReceiveTime = dbuf.Int64()
	m.From, _ = bytes2.ReadString(dbuf)
	len := int(dbuf.Int32())
	for i := 0; i < len; i++ {
		k, _ := bytes2.ReadString(dbuf)
		v, _ := bytes2.ReadString(dbuf)
		m.Options[k] = v
	}
	m.extend, _ = bytes2.ReadBytes(dbuf)
	return dbuf.Error()
}

func (m *ZProtoMetadata) String() string {
	return fmt.Sprintf("{server_id: %d, conn_id: %d, client_addr: %d, trace_id: %d, span_id: %d, recveive_time: %d, from: %s}",
		m.ServerId,
		m.ClientConnId,
		m.ClientAddr,
		m.TraceId,
		m.SpanId,
		m.ReceiveTime,
		m.From)
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoRawPayload struct {
	Payload []byte
}

func (m *ZProtoRawPayload) String() string {
	return fmt.Sprintf("{payload_len: %d, payload: %s}", len(m.Payload), bytes2.HexDump(m.Payload))
}

func (m *ZProtoRawPayload) Encode(x *bytes2.BufferOutput) {
	x.UInt32(PROTO)
	bytes2.WriteBytes(x, m.Payload)
}

func (m *ZProtoRawPayload) Decode(dbuf *bytes2.BufferInput) error {
	m.Payload, _ = bytes2.ReadBytes(dbuf)
	return dbuf.Error()
}

type ZProtoPing struct {
	PingId int64
}

////////////////////////////////////////////////////////////////////////////////
func (m *ZProtoPing) Encode(x *bytes2.BufferOutput) {
	x.UInt32(PING)
	x.Int64(m.PingId)
	// return x.Buf()
}

func (m *ZProtoPing) Decode(dbuf *bytes2.BufferInput) error {
	m.PingId = dbuf.Int64()
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoPong struct {
	MessageId uint64
	PingId    int64
}

func (m *ZProtoPong) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(20)
	x.UInt32(PONG)
	x.UInt64(m.MessageId)
	x.Int64(m.PingId)
	// return x.Buf()
}

func (m *ZProtoPong) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.MessageId = dbuf.UInt64()
	m.PingId = dbuf.Int64()
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoDrop struct {
	MessageId    int64
	ErrorCode    int32
	ErrorMessage string
}

func (m *ZProtoDrop) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(64)
	x.UInt32(DROP)
	x.Int64(m.MessageId)
	x.Int32(m.ErrorCode)
	bytes2.WriteString(x, m.ErrorMessage)
	// x.String(m.ErrorMessage)
	// return x.Buf()
}

func (m *ZProtoDrop) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.MessageId = dbuf.Int64()
	m.ErrorCode = dbuf.Int32()
	m.ErrorMessage, _ = bytes2.ReadString(dbuf) // dbuf.String()
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
// RPC很少使用这条消息
type ZProtoRedirect struct {
	Host    string
	Port    int
	Timeout int
}

func (m *ZProtoRedirect) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(64)
	x.UInt32(REDIRECT)
	// x.String(m.Host)
	bytes2.WriteString(x, m.Host)
	x.Int32(int32(m.Port))
	x.Int32(int32(m.Timeout))
	// return x.Buf()
}

func (m *ZProtoRedirect) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.Host, _ = bytes2.ReadString(dbuf) // dbuf.String()
	m.Port = int(dbuf.Int32())
	m.Timeout = int(dbuf.Int32())
	return nil
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoAck struct {
	ReceivedPackageIndex int
}

func (m *ZProtoAck) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(8)
	x.UInt32(ACK)
	x.Int32(int32(m.ReceivedPackageIndex))
	// return x.Buf()
}

func (m *ZProtoAck) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.ReceivedPackageIndex = int(dbuf.Int32())
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoHandshakeReq struct {
	ProtoRevision int
	RandomBytes   [32]byte
}

func (m *ZProtoHandshakeReq) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(64)
	x.UInt32(HANDSHAKE_REQ)
	x.Int32(int32(m.ProtoRevision))
	x.Bytes(m.RandomBytes[:])
	// return x.Buf()
}

func (m *ZProtoHandshakeReq) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.ProtoRevision = int(dbuf.Int32())
	randomBytes := dbuf.Bytes(32)
	copy(m.RandomBytes[:], randomBytes)
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoHandshakeRes struct {
	ProtoRevision int
	Sha1          [32]byte
}

func (m *ZProtoHandshakeRes) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(64)
	x.UInt32(HANDSHAKE_RSP)
	x.Int32(int32(m.ProtoRevision))
	x.Bytes(m.Sha1[:])
	// return x.Buf()
}

func (m *ZProtoHandshakeRes) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.ProtoRevision = int(dbuf.Int32())
	sha1 := dbuf.Bytes(32)
	copy(m.Sha1[:], sha1)
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoMarsSignal struct {
}

func (m *ZProtoMarsSignal) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(4)
	x.UInt32(MARS_SIGNAL)
	// return x.Buf()
}

func (m *ZProtoMarsSignal) Decode(dbuf *bytes2.BufferInput) error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoMessageAck struct {
	MessageIds []uint64
}

func (m *ZProtoMessageAck) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(512)
	x.UInt32(MESSAGE_ACK)
	x.Int32(int32(len(m.MessageIds)))
	for _, id := range m.MessageIds {
		x.UInt64(id)
	}
	// return x.Buf()
}

func (m *ZProtoMessageAck) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	len := int(dbuf.Int32())
	for i := 0; i < len; i++ {
		m.MessageIds = append(m.MessageIds, dbuf.UInt64())
	}
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoRpcRequest struct {
	MethodId string
	Body     []byte
}

func (m *ZProtoRpcRequest) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(512)
	x.UInt32(RPC_REQUEST)
	// x.String(m.MethodId)
	bytes2.WriteString(x, m.MethodId)
	// x.Bytes(m.Body)
	bytes2.WriteBytes(x, m.Body)
	// return x.Buf()
}

func (m *ZProtoRpcRequest) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	// m.MethodId = dbuf.String()
	m.MethodId, _ = bytes2.ReadString(dbuf)
	// m.Body = dbuf.StringBytes()
	m.Body, _ = bytes2.ReadBytes(dbuf)
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoRpcOk struct {
	RequestMessageId int64
	MethodResponseId string
	Body             []byte
}

func (m *ZProtoRpcOk) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(512)
	x.UInt32(RPC_OK)
	x.Int64(m.RequestMessageId)
	// x.String(m.MethodResponseId)
	bytes2.WriteString(x, m.MethodResponseId)
	// x.StringBytes(m.Body)
	bytes2.WriteBytes(x, m.Body)
	// return x.Buf()
}

func (m *ZProtoRpcOk) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.RequestMessageId = dbuf.Int64()
	// m.MethodResponseId = dbuf.String()
	m.MethodResponseId, _ = bytes2.ReadString(dbuf)
	// m.Body = dbuf.StringBytes()
	m.Body, _ = bytes2.ReadBytes(dbuf)
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoRpcError struct {
	RequestMessageId int64
	ErrorCode        int
	ErrorTag         string
	UserMessage      string
	CanTryAgain      bool
	ErrorData        []byte
}

func (m *ZProtoRpcError) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(512)
	x.UInt32(RPC_ERROR)
	x.Int64(m.RequestMessageId)
	x.Int32(int32(m.ErrorCode))
	// x.String(m.ErrorTag)
	bytes2.WriteString(x, m.ErrorTag)
	// x.String(m.UserMessage)
	bytes2.WriteString(x, m.UserMessage)
	if m.CanTryAgain {
		x.Int32(1)
	} else {
		x.Int32(0)
	}
	// x.StringBytes(m.ErrorData)
	bytes2.WriteBytes(x, m.ErrorData)

	// return x.Buf()
}

func (m *ZProtoRpcError) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.RequestMessageId = dbuf.Int64()
	m.ErrorCode = int(dbuf.Int32())
	// m.ErrorTag = dbuf.String()
	m.ErrorTag, _ = bytes2.ReadString(dbuf)
	// m.UserMessage = dbuf.String()
	m.UserMessage, _ = bytes2.ReadString(dbuf)
	canTryAgain := dbuf.Int32()
	m.CanTryAgain = canTryAgain == 1
	// m.ErrorData = dbuf.StringBytes()
	m.ErrorData, _ = bytes2.ReadBytes(dbuf)
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoRpcFloodWait struct {
	RequestMessageId int64
	Delay            int
}

func (m *ZProtoRpcFloodWait) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(512)
	x.UInt32(RPC_FLOOD_WAIT)
	x.Int64(m.RequestMessageId)
	x.Int32(int32(m.Delay))
	// return x.Buf()
}

func (m *ZProtoRpcFloodWait) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.RequestMessageId = dbuf.Int64()
	m.Delay = int(dbuf.Int32())
	return dbuf.Error()
}

////////////////////////////////////////////////////////////////////////////////
type ZProtoRpcInternalError struct {
	RequestMessageId int64
	CanTryAgain      bool
	TryAgainDelay    int
}

func (m *ZProtoRpcInternalError) Encode(x *bytes2.BufferOutput) {
	// x := NewEncodeBuf(512)
	x.UInt32(RPC_INTERNAL_ERROR)
	x.Int64(m.RequestMessageId)
	if m.CanTryAgain {
		x.Int32(1)
	} else {
		x.Int32(0)
	}
	x.Int32(int32(m.TryAgainDelay))
	// return x.Buf()
}

func (m *ZProtoRpcInternalError) Decode(dbuf *bytes2.BufferInput) error {
	// dbuf := NewDecodeBuf(b)
	m.RequestMessageId = dbuf.Int64()
	canTryAgain := dbuf.Int32()
	m.CanTryAgain = canTryAgain == 1
	m.TryAgainDelay = int(dbuf.Int32())
	return dbuf.Error()
}

func init() {
	zprotoFactories[PROTO] = func() MessageBase { return &ZProtoRawPayload{} }
	zprotoFactories[PING] = func() MessageBase { return &ZProtoPing{} }
	zprotoFactories[PONG] = func() MessageBase { return &ZProtoPong{} }
	zprotoFactories[DROP] = func() MessageBase { return &ZProtoDrop{} }
	zprotoFactories[REDIRECT] = func() MessageBase { return &ZProtoRedirect{} }
	zprotoFactories[ACK] = func() MessageBase { return &ZProtoAck{} }
	zprotoFactories[HANDSHAKE_REQ] = func() MessageBase { return &ZProtoHandshakeReq{} }
	zprotoFactories[HANDSHAKE_RSP] = func() MessageBase { return &ZProtoHandshakeRes{} }
	zprotoFactories[MARS_SIGNAL] = func() MessageBase { return &ZProtoMarsSignal{} }
	zprotoFactories[MESSAGE_ACK] = func() MessageBase { return &ZProtoMessageAck{} }
	zprotoFactories[RPC_REQUEST] = func() MessageBase { return &ZProtoRpcRequest{} }
	zprotoFactories[RPC_OK] = func() MessageBase { return &ZProtoRpcOk{} }
	zprotoFactories[RPC_ERROR] = func() MessageBase { return &ZProtoRpcError{} }
	zprotoFactories[RPC_FLOOD_WAIT] = func() MessageBase { return &ZProtoRpcFloodWait{} }
	zprotoFactories[RPC_INTERNAL_ERROR] = func() MessageBase { return &ZProtoRpcInternalError{} }
}

///////////////////////////////////////////////////////////////////////////////////////////////
func DecodeMessage(buf []byte) (MessageBase, error) {
	// glog.Info("decodeMessage: \n%s", bytes2.DumpSize(128, buf))

	var err error
	if len(buf) < 4 {
		err = fmt.Errorf("buf [len < 4]: len = %d", len(buf))
		glog.Error(err)
		return nil, err
	}

	return DecodeMessageByBuffer(bytes2.NewBufferInput(buf))
}

func DecodeMessageByBuffer(dbuf *bytes2.BufferInput) (MessageBase, error) {
	var err error
	msgType := dbuf.UInt32()
	m2 := NewZProtoMessage(msgType)
	if m2 == nil {
		err = fmt.Errorf("unregister msgType: %d, payload: \n%s", msgType, dbuf.DumpSize(128))
		glog.Error(err)
		return nil, err
	}

	err = m2.Decode(dbuf)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return m2, nil
}

func EncodeMessage(msg MessageBase) []byte {
	x := bytes2.NewBufferOutput(512)
	msg.Encode(x)
	return x.Buf()
}
