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
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nebulaim/telegramd/baselib/grpc_util/service_discovery"
	"github.com/nebulaim/telegramd/baselib/redis_client"
	"github.com/nebulaim/telegramd/proto/zproto"
)

var (
	confPath string
	Conf     *sessionConfig
)

type sessionConfig struct {
	ServerId             int32 // 服务器ID
	Redis                []redis_client.RedisConfig
	SaltCache            redis_client.RedisConfig
	// AuthKeyRpcClient     service_discovery.ServiceDiscoveryClientConfig
	BizRpcClient         service_discovery.ServiceDiscoveryClientConfig
	NbfsRpcClient        service_discovery.ServiceDiscoveryClientConfig
	SyncRpcClient        service_discovery.ServiceDiscoveryClientConfig
	AuthSessionRpcClient service_discovery.ServiceDiscoveryClientConfig
	Server               *zproto.ZProtoServerConfig
}

func init() {
	flag.StringVar(&confPath, "conf", "./session.toml", "config path")
}

func InitializeConfig() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	if err != nil {
		err = fmt.Errorf("decode file %s error: %v", confPath, err)
	}
	return
}
