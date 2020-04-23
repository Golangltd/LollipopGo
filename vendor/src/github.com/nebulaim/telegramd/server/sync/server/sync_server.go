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
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/base"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/mysql_client"
	"github.com/nebulaim/telegramd/baselib/net2"
	"github.com/nebulaim/telegramd/baselib/redis_client"
	"github.com/nebulaim/telegramd/biz/dal/dao"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/proto/zproto"
	"github.com/nebulaim/telegramd/server/sync/biz/core/update"
	"github.com/nebulaim/telegramd/service/idgen/client"
	"github.com/nebulaim/telegramd/service/status/client"
	"google.golang.org/grpc"
	"sync"
	"time"
	"github.com/nebulaim/telegramd/server/sync/server/rpc"
)

func init() {
	proto.RegisterType((*mtproto.TLSyncConnectToSessionServer)(nil), "mtproto.TLSyncConnectToSessionServer")
	proto.RegisterType((*mtproto.TLSyncSessionServerConnected)(nil), "mtproto.TLSyncSessionServerConnected")
	proto.RegisterType((*mtproto.TLSyncPushUpdatesData)(nil), "mtproto.TLSyncPushUpdatesData")
	proto.RegisterType((*mtproto.TLSyncPushRpcResultData)(nil), "mtproto.TLSyncPushRpcResultData")
	proto.RegisterType((*mtproto.PushData)(nil), "mtproto.PushData")
	proto.RegisterType((*mtproto.Bool)(nil), "mtproto.Bool")
}

type connContext struct {
	serverId  int32
	sessionId uint64
}

type syncServer struct {
	idgen idgen.UUIDGen
	// update     *update.UpdateModel
	status     status_client.StatusClient
	client     *zproto.ZProtoClient
	server     *grpc_util.RPCServer
	impl       *rpc.SyncServiceImpl
	sessionMap sync.Map
}

func NewSyncServer() *syncServer {
	return &syncServer{}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// AppInstance interface
func (s *syncServer) Initialize() error {
	var err error

	err = InitializeConfig()
	if err != nil {
		glog.Fatal(err)
		return err
	}
	glog.Info("config loaded: ", Conf)

	// idgen
	s.idgen, _ = idgen.NewUUIDGen("snowflake", base.Int32ToString(Conf.ServerId))

	// 初始化mysql_client、redis_client
	mysql_client.InstallMysqlClientManager(Conf.Mysql)
	redis_client.InstallRedisClientManager(Conf.Redis)

	// 初始化redis_dao、mysql_dao
	dao.InstallMysqlDAOManager(mysql_client.GetMysqlClientManager())
	dao.InstallRedisDAOManager(redis_client.GetRedisClientManager())

	s.status, _ = status_client.NewStatusClient("redis", "cache")

	s.server = grpc_util.NewRpcServer(Conf.Server.Addr, &Conf.Server.RpcDiscovery)
	s.client = zproto.NewZProtoClient("zproto", Conf.SessionClient, s)

	return nil
}

func (s *syncServer) RunLoop() {
	go s.server.Serve(func(s2 *grpc.Server) {
		updateModel := update.NewUpdateModel(Conf.ServerId, "immaster", "cache")
		s.impl = rpc.NewSyncService(s, s.status, updateModel)
		mtproto.RegisterRPCSyncServer(s2, s.impl)
	})
	s.client.Serve()
}

func (s *syncServer) Destroy() {
	if s.impl != nil {
		s.impl.Destroy()
	}

	s.server.Stop()
	s.client.Stop()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *syncServer) newMetadata() *zproto.ZProtoMetadata {
	md := &zproto.ZProtoMetadata{
		From:        "sync",
		ReceiveTime: time.Now().Unix(),
	}
	md.SpanId, _ = s.idgen.GetUUID()
	md.TraceId, _ = s.idgen.GetUUID()
	return md
}

func protoToSyncData(m proto.Message) (*zproto.ZProtoSyncData, error) {
	x := mtproto.NewEncodeBuf(128)
	n := proto.MessageName(m)
	glog.Info("message: name - ", n, ", ", m)
	x.Int(int32(len(n)))
	x.Bytes([]byte(n))
	b, err := proto.Marshal(m)
	x.Bytes(b)
	return &zproto.ZProtoSyncData{SyncRawData: x.GetBuf()}, err
}

///////////////////////////////////////////////////////////////////////////////////////
// Impl ZProtoClientCallBack
func (s *syncServer) OnNewClient(client *net2.TcpClient) {
	glog.Infof("OnNewConnection")
	connectToSession := &mtproto.TLSyncConnectToSessionServer{Data2: &mtproto.ConnectToServer_Data{
		SyncServerId: 1,
	}}

	req, _ := protoToSyncData(connectToSession)
	zproto.SendMessageByClient(client, s.newMetadata(), req)
}

func (s *syncServer) OnClientMessageArrived(client *net2.TcpClient, md *zproto.ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, msg zproto.MessageBase) error {
	switch msg.(type) {
	case *zproto.ZProtoSyncData:
		buf := msg.(*zproto.ZProtoSyncData).SyncRawData
		dbuf := mtproto.NewDecodeBuf(buf)
		len2 := int(dbuf.Int())
		messageName := string(dbuf.Bytes(len2))
		message, err := grpc_util.NewMessageByName(messageName)
		if err != nil {
			glog.Error(err)
			return err
		}

		err = proto.Unmarshal(buf[4+len2:], message)
		if err != nil {
			glog.Error(err)
			return err
		}

		switch message.(type) {
		case *mtproto.TLSyncSessionServerConnected:
			glog.Infof("onSyncData - request(SessionServerConnectedRsp): {%v}", message)
			// TODO(@benqi): bind server_id, server_name
			res, _ := message.(*mtproto.TLSyncSessionServerConnected)
			// res.GetServerId()
			ctx := &connContext{serverId: res.GetData2().SessionServerId, sessionId: client.GetConnection().GetConnID()}
			client.GetConnection().Context = ctx
			glog.Info("store serverId: ", ctx)
			s.sessionMap.Store(ctx.serverId, client)
			// glog.Info("store serverId: ", ctx)
		case *mtproto.Bool:
			glog.Infof("onSyncData - request(PushUpdatesData): {%v}", message)
		default:
			glog.Errorf("invalid register proto type: {%v}", message)
		}
	default:
		return fmt.Errorf("invalid payload type: %v", msg)
	}
	return nil
}

//func (s *syncServer) onSessionServerAdd(sid int32, sessionId int64) {
//	//s.mu.Lock()
//	//defer s.mu.Unlock()
//}

func (s *syncServer) OnClientClosed(client *net2.TcpClient) {
	glog.Infof("OnConnectionClosed")

	ctx := client.GetConnection().Context
	if ctx != nil {
		if connCtx, ok := ctx.(*connContext); ok {
			s.sessionMap.Delete(connCtx.serverId)
		}
	}
}

func (s *syncServer) OnClientTimer(client *net2.TcpClient) {
	// glog.Infof("OnTimer")
}

func (s *syncServer) SendToSessionServer(serverId int, m proto.Message) {
	if c, ok := s.sessionMap.Load(int32(serverId)); ok {
		client := c.(*net2.TcpClient)
		if client != nil {
			message, _ := protoToSyncData(m)
			zproto.SendMessageByClient(client, s.newMetadata(), message)
		} else {
			glog.Error("client type invalid")
		}
	} else {
		glog.Error("not found server id: ", serverId)
	}
}
