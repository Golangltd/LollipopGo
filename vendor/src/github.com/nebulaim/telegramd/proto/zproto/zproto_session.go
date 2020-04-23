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
	"encoding/hex"
	"fmt"
	// "github.com/nebulaim/telegramd/mtproto"
	"github.com/nebulaim/telegramd/baselib/bytes2"
)

// import "github.com/golang/glog"

const (
	STATE_ERROR = 0x0000

	STATE_CONNECTED2 = 0x0100
	STATE_HANDSHAKE  = 0x0200

	STATE_pq     = 0x0201
	STATE_pq_res = 0x0202
	STATE_pq_ack = 0x0203

	STATE_DH_params     = 0x0204
	STATE_DH_params_res = 0x0205
	STATE_DH_params_ack = 0x0206

	STATE_dh_gen     = 0x0207
	STATE_dh_gen_res = 0x0208
	STATE_dh_gen_ack = 0x0209

	STATE_AUTH_KEY = 0x0300
)

const (
	RES_STATE_NONE  = 0x00
	RES_STATE_OK    = 0x01
	RES_STATE_ERROR = 0x02
)

const (
	SESSION_HANDSHAKE             = 0xFF01
	SESSION_SESSION_DATA          = 0xFF02
	SYNC_DATA                     = 0xFF03
	SESSION_SESSION_CLIENT_NEW    = 0xFF04
	SESSION_SESSION_CLIENT_CLOSED = 0xFF05
)

//func isHandshake(state int) bool {
//	return state >= STATE_CONNECTED2 && state <= STATE_dh_gen_ack
//}

///////////////////////////////////////////////////////////////////////////////////////////
type HandshakeState struct {
	State    int    // 状态
	ResState int    // 后端握手返回的结果
	Ctx      []byte // 握手上下文数据，透传给后端
}

func (m *HandshakeState) String() string {
	return fmt.Sprintf("{state: %d, res_state: %d, ctx: %s}", m.State, m.ResState, hex.EncodeToString(m.Ctx))
}

func (m *HandshakeState) Encode(x *bytes2.BufferOutput) {
	x.Int32(int32(m.State))
	x.Int32(int32(m.ResState))
	bytes2.WriteBytes(x, m.Ctx)
	// return x.Buf()
}

func (m *HandshakeState) Decode(dbuf *bytes2.BufferInput) error {
	m.State = int(dbuf.Int32())
	m.ResState = int(dbuf.Int32())
	m.Ctx, _ = bytes2.ReadBytes(dbuf)
	// m.Payload = b
	return dbuf.Error()
}

//////////////////////////////////////////////////////////////////////////////
//func NewMTPRawMessage(authKeyId int64, quickAckId int32) *MTPRawMessage {
//	return &MTPRawMessage{
//		AuthKeyId:  authKeyId,
//		QuickAckId: quickAckId,
//	}
//}
//
//// 代理使用
//type MTPRawMessage struct {
//	AuthKeyId  int64 // 由原始数据解压获得
//	QuickAckId int32 // EncryptedMessage，则可能存在
//	// 原始数据
//	Payload []byte
//}
//
//func (m *MTPRawMessage) Encode(x *bytes2.BufferOutput) {
//	bytes2.WriteBytes(x, m.Payload)
//	// return x.Buf()
//}
//
//func (m *MTPRawMessage) Decode(dbuf *bytes2.BufferInput) error {
//	m.Payload, _ = bytes2.ReadBytes(dbuf)
//	return dbuf.Error()
//}

///////////////////////////////////////////////////////////////////////////////////////////
type ZProtoHandshakeMessage struct {
	SessionId  uint64
	State      *HandshakeState
	MTPRawData []byte
}

func (m *ZProtoHandshakeMessage) String() string {
	return fmt.Sprintf("{session_id: %d, state: %s, mtp_raw_data_len: %d, mtp_raw_data: %s}",
		m.SessionId,
		m.State,
		len(m.MTPRawData),
		bytes2.HexDump(m.MTPRawData))
}

func (m *ZProtoHandshakeMessage) Encode(x *bytes2.BufferOutput) {
	x.UInt32(SESSION_HANDSHAKE)
	x.UInt64(m.SessionId)
	m.State.Encode(x)
	bytes2.WriteBytes(x, m.MTPRawData)
}

func (m *ZProtoHandshakeMessage) Decode(dbuf *bytes2.BufferInput) error {
	m.SessionId = dbuf.UInt64()
	state := &HandshakeState{}
	if state.Decode(dbuf); dbuf.Error() != nil {
		return dbuf.Error()
	}
	m.State = state
	m.MTPRawData, _ = bytes2.ReadBytes(dbuf)
	return dbuf.Error()
}

///////////////////////////////////////////////////////////////////////////////////////////
type ZProtoSessionData struct {
	ConnType   int
	SessionId  uint64
	MtpRawData []byte
}

func (m *ZProtoSessionData) String() string {
	return fmt.Sprintf("{conn_type: %d, session_id: %d, mtp_raw_data_len: %d, mtp_raw_data: %s}",
		m.ConnType,
		m.SessionId,
		len(m.MtpRawData),
		bytes2.HexDump(m.MtpRawData))
}

func (m *ZProtoSessionData) Encode(x *bytes2.BufferOutput) {
	x.UInt32(SESSION_SESSION_DATA)
	x.Int32(int32(m.ConnType))
	x.UInt64(m.SessionId)
	bytes2.WriteBytes(x, m.MtpRawData)
}

func (m *ZProtoSessionData) Decode(dbuf *bytes2.BufferInput) error {
	m.ConnType = int(dbuf.Int32())
	m.SessionId = dbuf.UInt64()
	m.MtpRawData, _ = bytes2.ReadBytes(dbuf)
	return dbuf.Error()
}

///////////////////////////////////////////////////////////////////////////////////////////
type ZProtoSyncData struct {
	SyncRawData []byte
}

func (m *ZProtoSyncData) String() string {
	return fmt.Sprintf("{sync_raw_data_len: %d, sync_raw_data: %s}",
		len(m.SyncRawData),
		bytes2.HexDump(m.SyncRawData))
}

func (m *ZProtoSyncData) Encode(x *bytes2.BufferOutput) {
	x.UInt32(SYNC_DATA)
	bytes2.WriteBytes(x, m.SyncRawData)
}

func (m *ZProtoSyncData) Decode(dbuf *bytes2.BufferInput) error {
	m.SyncRawData, _ = bytes2.ReadBytes(dbuf)
	return dbuf.Error()
}

///////////////////////////////////////////////////////////////////////////////////////////
type ZProtoSessionClientNew struct {
	// proto int32
	ConnType  int
	SessionId uint64
	AuthKeyId int64
}

func (m *ZProtoSessionClientNew) String() string {
	return fmt.Sprintf("{conn_type: %d, session_id: %d, auth_key_id: %d}", m.ConnType, m.SessionId, m.AuthKeyId)
}

func (m *ZProtoSessionClientNew) Encode(x *bytes2.BufferOutput) {
	x.UInt32(SESSION_SESSION_CLIENT_NEW)
	x.Int32(int32(m.ConnType))
	x.UInt64(m.SessionId)
	x.Int64(m.AuthKeyId)
}

func (m *ZProtoSessionClientNew) Decode(dbuf *bytes2.BufferInput) error {
	m.ConnType = int(dbuf.Int32())
	m.SessionId = dbuf.UInt64()
	m.AuthKeyId = dbuf.Int64()
	return dbuf.Error()
}

///////////////////////////////////////////////////////////////////////////////////////////
type ZProtoSessionClientClosed struct {
	ConnType  int
	SessionId uint64
	AuthKeyId int64
}

func (m *ZProtoSessionClientClosed) String() string {
	return fmt.Sprintf("{conn_type: %d, session_id: %d, auth_key_id: %d}", m.ConnType, m.SessionId, m.AuthKeyId)
}

func (m *ZProtoSessionClientClosed) Encode(x *bytes2.BufferOutput) {
	x.UInt32(SESSION_SESSION_CLIENT_CLOSED)
	x.Int32(int32(m.ConnType))
	x.UInt64(m.SessionId)
	x.Int64(m.AuthKeyId)
}

func (m *ZProtoSessionClientClosed) Decode(dbuf *bytes2.BufferInput) error {
	m.ConnType = int(dbuf.Int32())
	m.SessionId = dbuf.UInt64()
	m.AuthKeyId = dbuf.Int64()
	return dbuf.Error()
}

///////////////////////////////////////////////////////////////////////////////////////////
//type SessionHandshakeMessage struct {
//	State      *HandshakeState
//	MTPMessage *UnencryptedMessage
//}
//
//func (m *SessionHandshakeMessage) Encode(x *bytes2.BufferOutput) []byte {
//	x := NewEncodeBuf(512)
//	x.UInt(SESSION_HANDSHAKE)
//	x.Int(int32(m.State.State))
//	x.Int(int32(m.State.ResState))
//	x.StringBytes(m.State.Ctx)
//	x.Bytes(m.MTPMessage.Encode())
//	return x.GetBuf()
//}
//
//func (m *SessionHandshakeMessage) Decode(dbuf *bytes2.BufferInput) error {
//	// glog.Info(b)
//	dbuf := NewDecodeBuf(b)
//	m.State.State = int(dbuf.Int())
//	m.State.ResState = int(dbuf.Int())
//	m.State.Ctx = dbuf.StringBytes()
//	m.MTPMessage = &UnencryptedMessage{}
//	err := dbuf.GetError()
//	if err == nil {
//		dbuf.Long()
//		err = m.MTPMessage.Decode(b[dbuf.off:])
//	}
//	return err
//}
//
/////////////////////////////////////////////////////////////////////////////////////////////
//type SessionDataMessage struct {
//	MTPMessage *EncryptedMessage2
//}
//
//func (m *SessionDataMessage) Encode(x *bytes2.BufferOutput) []byte {
//	x := NewEncodeBuf(512)
//	x.UInt(SESSION_SESSION_DATA)
//	// x.Bytes(m.MTPMessage.Encode())
//	return x.GetBuf()
//}
//
//func (m *SessionDataMessage) Decode(dbuf *bytes2.BufferInput) error {
//	//m.MTPMessage = &EncryptedMessage2{}
//	//return m.MTPMessage.Decode(b)
//	return nil
//}

func init() {
	zprotoFactories[SESSION_HANDSHAKE] = func() MessageBase { return &ZProtoHandshakeMessage{} }
	zprotoFactories[SESSION_SESSION_DATA] = func() MessageBase { return &ZProtoSessionData{} }
	zprotoFactories[SYNC_DATA] = func() MessageBase { return &ZProtoSyncData{} }
	zprotoFactories[SESSION_SESSION_CLIENT_NEW] = func() MessageBase { return &ZProtoSessionClientNew{} }
	zprotoFactories[SESSION_SESSION_CLIENT_CLOSED] = func() MessageBase { return &ZProtoSessionClientClosed{} }
}
