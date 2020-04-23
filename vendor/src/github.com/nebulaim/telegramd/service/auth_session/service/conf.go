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
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/mysql_client"
	"github.com/nebulaim/telegramd/baselib/redis_client"
)

var (
	confPath string
	Conf     *authSessionConfig
)

type authSessionConfig struct {
	Redis     []redis_client.RedisConfig
	Mysql     []mysql_client.MySQLConfig
	RpcServer *grpc_util.RPCServerConfig
}

func (c *authSessionConfig) String() string {
	return fmt.Sprintf("{redis: %v. mysql: %v, server: %v}",
		c.Redis,
		c.Mysql,
		c.RpcServer)
}

func init() {
	flag.StringVar(&confPath, "conf", "./auth_session.toml", "config path")
}

func InitializeConfig() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	if err != nil {
		err = fmt.Errorf("decode file %s error: %v", confPath, err)
	}
	return
}
