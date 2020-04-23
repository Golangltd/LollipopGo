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
	"github.com/nebulaim/telegramd/baselib/etcd_util"
	"github.com/nebulaim/telegramd/baselib/grpc_util/service_discovery"
	"github.com/nebulaim/telegramd/baselib/grpc_util/service_discovery/etcd3"
	"github.com/nebulaim/telegramd/baselib/net2"
	"net"
)

type ZProtoServerCallback interface {
	OnServerNewConnection(conn *net2.TcpConnection)
	OnServerMessageDataArrived(conn *net2.TcpConnection, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, msg MessageBase) error
	OnServerConnectionClosed(conn *net2.TcpConnection)
}

type ZProtoServerConfig struct {
	Server    net2.ServerConfig
	Discovery service_discovery.ServiceDiscoveryServerConfig
}

type ZProtoServer struct {
	server   *net2.TcpServer
	registry *etcd3.EtcdReigistry
	callback ZProtoServerCallback
}

func NewZProtoServer(conf *ZProtoServerConfig, cb ZProtoServerCallback) *ZProtoServer {
	lsn, err := net.Listen("tcp", conf.Server.Addr)
	if err != nil {
		glog.Fatal("listen error: %v", err)
	}

	server := &ZProtoServer{
		callback: cb,
	}
	server.server = net2.NewTcpServer(net2.TcpServerArgs{
		Listener:           lsn,
		ServerName:         conf.Server.Name,
		ProtoName:          conf.Server.ProtoName,
		SendChanSize:       1024,
		ConnectionCallback: server,
	}) // todo (yumcoder): set max connection

	server.registry, err = etcd_util.NewEtcdRegistry(conf.Discovery)
	if err != nil {
		glog.Fatal(err)
	}

	return server
}

///////////////////////////////////////////////////////////////////////////////////////////////
func (s *ZProtoServer) Serve() {
	go s.server.Serve()
	go s.registry.Register()
}

func (s *ZProtoServer) Stop() {
	s.registry.Deregister()
	s.server.Stop()
}

func (s *ZProtoServer) Pause() {
	s.server.Pause()
}

///////////////////////////////////////////////////////////////////////////////////////////////
func (s *ZProtoServer) SendMessageByConnID(connID uint64, md *ZProtoMetadata, msg MessageBase) error {
	conn := s.server.GetConnection(connID)
	if conn != nil {
		return SendMessageByConn(conn, md, msg)
	} else {
		err := fmt.Errorf("send data error, conn offline, connID: %d", connID)
		glog.Error(err)
		return err
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
func (s *ZProtoServer) OnNewConnection(conn *net2.TcpConnection) {
	glog.Info("onNewConnection - ", conn)

	////////////////////////////////////////////////////////////////
	// @benqi: hack
	codec := conn.Codec()
	codec.(*ZProtoCodec).connID = conn.GetConnID()

	if s.callback != nil {
		s.callback.OnServerNewConnection(conn)
	}
}

func (s *ZProtoServer) OnConnectionDataArrived(conn *net2.TcpConnection, msg interface{}) error {
	zmsg, ok := msg.(*ZProtoMessage)
	if !ok {
		return fmt.Errorf("recv invalid zmsg - {%v}", zmsg)
	}

	// TODO(@benqi): check sessionId and seqNo
	switch zmsg.Message.(type) {
	case *ZProtoRawPayload:
		s.onRawPayload(conn, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoRawPayload))
	case *ZProtoPing:
		s.onPing(conn, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoPing))
	case *ZProtoAck:
		s.onAck(conn, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoAck))
	case *ZProtoHandshakeReq:
		s.onHandshakeReq(conn, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoHandshakeReq))
	case *ZProtoMarsSignal:
		s.onMarsSignal(conn, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoMarsSignal))
	case *ZProtoRpcRequest:
		s.onRpcRequest(conn, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoRpcRequest))
	default:
		err := fmt.Errorf("invalid message - {conn: %s, msg: {%v}}", conn, zmsg.Message)
		glog.Error(err)
		return err
	}

	return nil
}

func (s *ZProtoServer) OnConnectionClosed(conn *net2.TcpConnection) {
	glog.Info("onConnectionClosed - ", conn)

	if s.callback != nil {
		s.callback.OnServerConnectionClosed(conn)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
func (s *ZProtoServer) onRawPayload(conn *net2.TcpConnection, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, payload *ZProtoRawPayload) error {
	glog.Info("onRawPayload - conn: %s, md: %s, payload_len: %d", conn, md, len(payload.Payload))

	var (
		err error
		m2  MessageBase
	)

	if s.callback != nil {
		m2, err = DecodeMessage(payload.Payload)
		if err != nil {
			return err
		}
		err = s.callback.OnServerMessageDataArrived(conn, md, sessionId, messageId, seqNo, m2)
	}

	return err
}

func (s *ZProtoServer) onPing(conn *net2.TcpConnection, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, ping *ZProtoPing) error {
	glog.Info("onPing: ", ping)

	zmsg := &ZProtoMessage{
		// SessionId: client.GetConnection().GetConnID(),
		Message: &ZProtoPong{
			MessageId: messageId,
			PingId:    ping.PingId,
		},
	}

	return conn.Send(zmsg)
}

func (s *ZProtoServer) onAck(conn *net2.TcpConnection, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, ack *ZProtoAck) error {
	glog.Info("onAck: ", ack)

	return nil
}

func (s *ZProtoServer) onHandshakeReq(conn *net2.TcpConnection, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, handshake *ZProtoHandshakeReq) error {
	glog.Info("onHandshakeReq: ", handshake)

	// TODO(@benqi): check ProtoRevision
	zmsg := &ZProtoMessage{
		// SessionId: client.GetConnection().GetConnID(),
		Message: &ZProtoHandshakeRes{
			ProtoRevision: handshake.ProtoRevision,
			Sha1:          handshake.RandomBytes,
		},
	}

	return conn.Send(zmsg)
}

func (s *ZProtoServer) onMarsSignal(conn *net2.TcpConnection, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, marsSignal *ZProtoMarsSignal) error {
	glog.Info("onMarsSignal: ", marsSignal)

	// wechat open source mars - marsSignal support
	return nil
}

func (s *ZProtoServer) onRpcRequest(conn *net2.TcpConnection, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, rpc *ZProtoRpcRequest) error {
	glog.Info("onRpcRequest: ", rpc)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func SendMessageByConn(conn *net2.TcpConnection, md *ZProtoMetadata, msg MessageBase) error {
	zmsg := &ZProtoMessage{
		Metadata: md,
		Message:  &ZProtoRawPayload{Payload: EncodeMessage(msg)},
	}
	return conn.Send(zmsg)
}
