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
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/mysql_client"
	"github.com/nebulaim/telegramd/baselib/net2"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/proto/zproto"
	"github.com/nebulaim/telegramd/server/access/auth_key/dal/dao"
	"time"
)

type AuthKeyServer struct {
	handshake            *handshake
	server               *zproto.ZProtoServer
	authSessionRpcClient mtproto.RPCSessionClient
}

func NewAuthKeyServer() *AuthKeyServer {
	return &AuthKeyServer{}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// AppInstance interface
func (s *AuthKeyServer) Initialize() error {
	err := InitializeConfig()
	if err != nil {
		glog.Fatal(err)
		return err
	}
	glog.Info("load conf: ", Conf)

	// 初始化mysql_client、redis_client
	mysql_client.InstallMysqlClientManager(Conf.Mysql)
	// 初始化redis_dao、mysql_dao
	dao.InstallMysqlDAOManager(mysql_client.GetMysqlClientManager())

	s.server = zproto.NewZProtoServer(Conf.Server, s)
	// s.rpcServer = grpc_util.NewRpcServer(Conf.RpcServer.Addr, &Conf.RpcServer.RpcDiscovery)

	return nil
}

func (s *AuthKeyServer) RunLoop() {
	c, _ := grpc_util.NewRPCClient(&Conf.AuthSessionRpcClient)
	s.authSessionRpcClient = mtproto.NewRPCSessionClient(c.GetClientConn())
	s.handshake = newHandshake(s.authSessionRpcClient)

	go s.server.Serve()
}

func (s *AuthKeyServer) Destroy() {
	glog.Infof("sessionServer - destroy...")
	s.server.Stop()
	// s.rpcServer.Stop()
	time.Sleep(1 * time.Second)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *AuthKeyServer) OnServerNewConnection(conn *net2.TcpConnection) {
	glog.Infof("onNewConnection %v", conn.RemoteAddr())
}

func (s *AuthKeyServer) OnServerMessageDataArrived(conn *net2.TcpConnection, md *zproto.ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, msg zproto.MessageBase) error {
	glog.Infof("onServerMessageDataArrived - msg: %v", msg)

	hmsg, ok := msg.(*zproto.ZProtoHandshakeMessage)
	if !ok {
		err := fmt.Errorf("invalid handshakeMessage: {%v}", msg)
		glog.Error(err)
		return err
	}

	hrsp, err := s.handshake.onHandshake(conn, hmsg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	// Fix onMsgAck return nil bug.
	if hrsp == nil {
		return nil
	}

	return zproto.SendMessageByConn(conn, md, hrsp)
}

func (s *AuthKeyServer) OnServerConnectionClosed(conn *net2.TcpConnection) {
	glog.Infof("onConnectionClosed - %v", conn.RemoteAddr())
}
