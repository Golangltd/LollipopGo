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

package service

import (
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/mysql_client"
	"github.com/nebulaim/telegramd/baselib/redis_client"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"google.golang.org/grpc"
	"github.com/nebulaim/telegramd/service/auth_session/service/rpc"
)

type authSessionServer struct {
	rpcServer *grpc_util.RPCServer
}

func NewAuthSessionServer() *authSessionServer {
	return &authSessionServer{}
}

// AppInstance interface
func (s *authSessionServer) Initialize() error {
	glog.Infof("authSessionServer - initialize...")

	err := InitializeConfig()
	if err != nil {
		glog.Fatal(err)
		return err
	}
	glog.Info("authSessionServer - load conf: ", Conf)

	// 初始化mysql_client、redis_client
	redis_client.InstallRedisClientManager(Conf.Redis)
	mysql_client.InstallMysqlClientManager(Conf.Mysql)

	s.rpcServer = grpc_util.NewRpcServer(Conf.RpcServer.Addr, &Conf.RpcServer.RpcDiscovery)
	return nil
}

func (s *authSessionServer) RunLoop() {
	glog.Infof("authSessionServer - runLoop...")

	// TODO(@benqi): check error
	s.rpcServer.Serve(func(s2 *grpc.Server) {
		impl := rpc.NewSessionServiceImpl("immaster", "cache")
		// &rpc.SessionServiceImpl{}
		mtproto.RegisterRPCSessionServer(s2, impl)
	})
}

func (s *authSessionServer) Destroy() {
	glog.Infof("authSessionServer - destroy...")
	s.rpcServer.Stop()
	//time.Sleep(1*time.Second)
}
