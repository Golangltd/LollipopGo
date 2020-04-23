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

package server

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/base"
	"github.com/nebulaim/telegramd/baselib/net2"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/proto/zproto"
	"github.com/nebulaim/telegramd/service/idgen/client"
	"sync"
	"time"
)

type handshakeState struct {
	state    int    // 状态
	resState int    // 后端握手返回的结果
	ctx      []byte // 握手上下文数据，透传给后端
}

type connContext struct {
	// TODO(@benqi): lock
	sync.Mutex
	state          int // 是否握手阶段
	md             *zproto.ZProtoMetadata
	handshakeState *zproto.HandshakeState
	seqNum         uint64

	sessionAddr string
	authKeyId   int64
}

func (ctx *connContext) getState() int {
	ctx.Lock()
	defer ctx.Unlock()
	return ctx.state
}

func (ctx *connContext) setState(state int) {
	ctx.Lock()
	defer ctx.Unlock()
	if ctx.state != state {
		ctx.state = state
	}
}

func (ctx *connContext) encryptedMessageAble() bool {
	ctx.Lock()
	defer ctx.Unlock()
	return ctx.state == zproto.STATE_CONNECTED2 ||
		ctx.state == zproto.STATE_AUTH_KEY ||
		(ctx.state == zproto.STATE_HANDSHAKE &&
			(ctx.handshakeState.State == zproto.STATE_pq_res ||
				(ctx.handshakeState.State == zproto.STATE_dh_gen_res &&
					ctx.handshakeState.ResState == zproto.RES_STATE_OK)))

}

type FrontendServer struct {
	idgen      idgen.UUIDGen
	server80   *mtproto.MTProtoServer
	server443  *mtproto.MTProtoServer
	server5222 *mtproto.MTProtoServer
	client     *zproto.ZProtoClient
}

func NewFrontendServer() *FrontendServer {
	return &FrontendServer{}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// AppInstance interface
func (s *FrontendServer) Initialize() error {
	err := InitializeConfig()
	if err != nil {
		glog.Fatal(err)
		return err
	}
	glog.Info("load conf: ", Conf)

	// idgen
	s.idgen, _ = idgen.NewUUIDGen("snowflake", base.Int32ToString(Conf.ServerId))

	// mtproto_server
	s.server80 = mtproto.NewMTProtoServer(Conf.Server80, s)
	s.server443 = mtproto.NewMTProtoServer(Conf.Server443, s)
	s.server5222 = mtproto.NewMTProtoServer(Conf.Server5222, s)

	// client
	s.client = zproto.NewZProtoClient("zproto", Conf.Clients, s)
	return nil
}

func (s *FrontendServer) RunLoop() {
	s.server80.Serve()
	s.server443.Serve()
	s.server5222.Serve()
	s.client.Serve()
}

func (s *FrontendServer) Destroy() {
	s.server80.Stop()
	s.server443.Stop()
	s.server5222.Stop()
	s.client.Stop()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *FrontendServer) newMetadata(conn *net2.TcpConnection) *zproto.ZProtoMetadata {
	md := &zproto.ZProtoMetadata{
		ServerId:     int(Conf.ServerId),
		ClientConnId: conn.GetConnID(),
		ClientAddr:   conn.RemoteAddr().String(),
		From:         "frontend",
		ReceiveTime:  time.Now().Unix(),
	}
	md.SpanId, _ = s.idgen.GetUUID()
	md.TraceId, _ = s.idgen.GetUUID()
	return md
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MTProtoServerCallback
func (s *FrontendServer) OnServerNewConnection(conn *net2.TcpConnection) {
	conn.Context = &connContext{
		state: zproto.STATE_CONNECTED2,
		md: &zproto.ZProtoMetadata{
			ServerId:     int(Conf.ServerId),
			ClientConnId: conn.GetConnID(),
			ClientAddr:   conn.RemoteAddr().String(),
			From:         "frontend",
		},
		handshakeState: &zproto.HandshakeState{
			State:    zproto.STATE_CONNECTED2,
			ResState: zproto.RES_STATE_NONE,
		},
		seqNum: 1,
	}
	glog.Infof("onServerNewConnection - {peer: %s, ctx: {%v}}", conn, conn.Context)
}

func (s *FrontendServer) OnServerMessageDataArrived(conn *net2.TcpConnection, msg *mtproto.MTPRawMessage) error {
	md := s.newMetadata(conn)
	glog.Infof("onServerMessageDataArrived - receive data: {peer: %s, md: %s, msg: %s}", conn, md, msg)

	ctx, _ := conn.Context.(*connContext)

	var err error
	if msg.AuthKeyId() == 0 {
		if ctx.getState() == zproto.STATE_AUTH_KEY {
			err = fmt.Errorf("invalid state STATE_AUTH_KEY: %d", ctx.getState())
			glog.Errorf("process msg error: {%v} - {peer: %s, md: %s, msg: %s}", err, conn, md, msg)
			conn.Close()
		} else {
			err = s.onServerUnencryptedRawMessage(ctx, conn, md, msg)
		}
	} else {
		if !ctx.encryptedMessageAble() {
			err = fmt.Errorf("invalid state: {state: %d, handshakeState: {%v}}", ctx.state, ctx.handshakeState)
			glog.Errorf("process msg error: {%v} - {peer: %s, md: %s, msg: %s}", err, conn, md, msg)
			conn.Close()
		} else {
			err = s.onServerEncryptedRawMessage(ctx, conn, md, msg)
		}
	}

	return err
}

func (s *FrontendServer) OnServerConnectionClosed(conn *net2.TcpConnection) {
	glog.Infof("onServerConnectionClosed - {peer: %s}", conn)
	// s.sendClientClosed(conn)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// ZProtoClientCallBack
func (s *FrontendServer) OnNewClient(client *net2.TcpClient) {
	glog.Infof("onNewClient - peer(%s)", client.GetConnection())
}

func (s *FrontendServer) OnClientMessageArrived(client *net2.TcpClient, md *zproto.ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, msg zproto.MessageBase) error {
	var err error

	switch msg.(type) {
	case *zproto.ZProtoHandshakeMessage:
		err = s.onClientHandshakeMessage(client, md, msg.(*zproto.ZProtoHandshakeMessage))
	case *zproto.ZProtoSessionData:
		err = s.onClientSessionData(client, md, msg.(*zproto.ZProtoSessionData))
	default:
		err = fmt.Errorf("invalid msg: %v", msg)
		glog.Errorf("onClientMessageArrived - invalid msg: peer(%s), zmsg: {%v}",
			client.GetConnection(),
			msg)
	}

	return err
}

func (s *FrontendServer) OnClientClosed(client *net2.TcpClient) {
	glog.Infof("onClientClosed - peer(%s)", client.GetConnection())
}

func (s *FrontendServer) OnClientTimer(client *net2.TcpClient) {
	// Impl timer logic
	glog.Info("onClientTimer")
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *FrontendServer) onClientHandshakeMessage(client *net2.TcpClient, md *zproto.ZProtoMetadata, handshake *zproto.ZProtoHandshakeMessage) error {
	glog.Infof("onClientHandshakeMessage - handshake: peer(%s), state: {%v}",
		client.GetConnection(),
		handshake.State)

	///////////////////////////////////////////////////////////////////
	conn := s.getConnBySessionID(handshake.SessionId)
	// s.server443.GetConnection(zmsg.SessionId)
	if conn == nil {
		glog.Warning("conn closed, connID = ", handshake.SessionId)
		return nil
	}

	if handshake.State.ResState == zproto.RES_STATE_ERROR {
		// TODO(@benqi): Close.
		glog.Warning(" handshake.State.ResState error, connID = ", handshake.SessionId)
		// conn.Close()
		return nil
	} else {
		ctx := conn.Context.(*connContext)
		ctx.Lock()
		ctx.handshakeState = handshake.State
		ctx.Unlock()

		glog.Infof("onClientHandshakeMessage - sendToClient to: {peer: %s, md: %s, handshake: %s}",
			conn,
			md,
			handshake)

		return conn.Send(&mtproto.MTPRawMessage{Payload: handshake.MTPRawData})
	}
}

func (s *FrontendServer) onClientSessionData(client *net2.TcpClient, md *zproto.ZProtoMetadata, sessData *zproto.ZProtoSessionData) error {
	///////////////////////////////////////////////////////////////////
	conn := s.getConnBySessionID(sessData.SessionId)
	// s.server443.GetConnection(zmsg.SessionId)
	if conn == nil {
		glog.Warning("conn closed, connID = ", sessData.SessionId)
		return nil
	}

	glog.Infof("onClientSessionData - sendToClient to: {peer: %s, md: %s, sessData: %s}",
		conn,
		md,
		sessData)
	return conn.Send(&mtproto.MTPRawMessage{Payload: sessData.MtpRawData})
}

func (s *FrontendServer) genSessionId(conn *net2.TcpConnection) uint64 {
	var sid = conn.GetConnID()
	if conn.Name() == "frontend443" {
		// sid = sid | 0 << 60
	} else if conn.Name() == "frontend80" {
		sid = sid | 1<<60
	} else if conn.Name() == "frontend5222" {
		sid = sid | 2<<60
	}

	return sid
}

func (s *FrontendServer) getConnBySessionID(id uint64) *net2.TcpConnection {
	//
	var server *mtproto.MTProtoServer
	sid := id >> 60
	if sid == 0 {
		server = s.server443
	} else if sid == 1 {
		server = s.server80
	} else if sid == 2 {
		server = s.server5222
	} else {
		return nil
	}

	id = id & 0xffffffffffffff
	return server.GetConnection(id)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *FrontendServer) onServerUnencryptedRawMessage(ctx *connContext, conn *net2.TcpConnection, md *zproto.ZProtoMetadata, mmsg *mtproto.MTPRawMessage) error {
	glog.Infof("onServerUnencryptedRawMessage - receive data: {peer: %s, md: %s, ctx: %s, msg: %s}", conn, ctx, md, mmsg)

	ctx.Lock()
	if ctx.state == zproto.STATE_CONNECTED2 {
		ctx.state = zproto.STATE_HANDSHAKE
	}
	if ctx.handshakeState.State == zproto.STATE_CONNECTED2 {
		ctx.handshakeState.State = zproto.STATE_pq
	}
	ctx.Unlock()

	// sentToClient
	hmsg := &zproto.ZProtoHandshakeMessage{
		SessionId:  s.genSessionId(conn),
		State:      ctx.handshakeState,
		MTPRawData: mmsg.Payload,
	}

	// glog.Infof("SendMessage - handshake: {peer: %s, md: %s, ctx: %s, msg: %s}", conn, ctx, md, mmsg)
	return s.client.SendMessage("handshake", s.newMetadata(conn), hmsg)
}

func (s *FrontendServer) onServerEncryptedRawMessage(ctx *connContext, conn *net2.TcpConnection, md *zproto.ZProtoMetadata, mmsg *mtproto.MTPRawMessage) error {
	glog.Infof("onServerEncryptedRawMessage - receive data: {peer: %s, md: %s, ctx: %s, msg: %s}", conn, ctx, md, mmsg)

	sessData := &zproto.ZProtoSessionData{
		ConnType:   mmsg.ConnType(),
		SessionId:  s.genSessionId(conn),
		MtpRawData: mmsg.Payload,
	}

	return s.client.SendKetamaMessage("session", base.Int64ToString(mmsg.AuthKeyId()), md, sessData, func(addr string) {
		// s.checkAndSendClientNew(ctx, conn, addr, mmsg.AuthKeyId(), md)
	})
}

func (s *FrontendServer) checkAndSendClientNew(ctx *connContext, conn *net2.TcpConnection, kaddr string, authKeyId int64, md *zproto.ZProtoMetadata) error {
	var err error
	if ctx.sessionAddr == "" {
		clientNew := &zproto.ZProtoSessionClientNew{
			// MTPMessage: mmsg,
		}
		err = s.client.SendMessageToAddress("session", kaddr, s.newMetadata(conn), clientNew)
		if err == nil {
			ctx.sessionAddr = kaddr
			ctx.authKeyId = authKeyId
		} else {
			glog.Error(err)
		}
	} else {
		// TODO(@benqi): check ctx.sessionAddr == kaddr
	}

	return err
}

func (s *FrontendServer) sendClientClosed(conn *net2.TcpConnection) {
	if conn.Context == nil {
		return
	}

	ctx, _ := conn.Context.(*connContext)
	if ctx.sessionAddr == "" || ctx.authKeyId == 0 {
		return
	}

	s.client.SendKetamaMessage("session", base.Int64ToString(ctx.authKeyId), s.newMetadata(conn), &zproto.ZProtoSessionClientClosed{}, nil)
}
