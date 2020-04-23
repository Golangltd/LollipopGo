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

package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nebulaim/telegramd/baselib/app"
	"github.com/nebulaim/telegramd/baselib/net2"
	"github.com/nebulaim/telegramd/proto/zproto"
)

var (
	confPath string
	Conf     *zproto.ZProtoServerConfig
)

type TestServerInsance struct {
	server *zproto.ZProtoServer
}

func (s *TestServerInsance) Initialize() error {
	InitializeConfig()
	s.server = zproto.NewZProtoServer(Conf, s)
	return nil
}

func (s *TestServerInsance) RunLoop() {
	// go this.server.httpServer.Serve(this.server.httpListener)
	s.server.Serve()
	// this.client.Serve()
}

func (s *TestServerInsance) Destroy() {
	s.server.Stop()
}

func (s *TestServerInsance) OnServerNewConnection(conn *net2.TcpConnection) {

}

func (s *TestServerInsance) OnServerMessageDataArrived(c *net2.TcpConnection, md *zproto.ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, msg zproto.MessageBase) error {
	return nil
}

func (s *TestServerInsance) OnServerConnectionClosed(c *net2.TcpConnection) {

}

func init() {
	flag.StringVar(&confPath, "conf", "./test_server.toml", "config path")
}

func InitializeConfig() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	if err != nil {
		err = fmt.Errorf("decode file %s error: %v", confPath, err)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////
func main() {
	instance := &TestServerInsance{}
	app.DoMainAppInstance(instance)
}
