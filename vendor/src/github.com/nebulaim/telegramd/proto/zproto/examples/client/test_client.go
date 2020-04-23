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
	Conf     *zproto.ZProtoClientConfig
)

type TestClientInsance struct {
	client *zproto.ZProtoClient
}

func (inst *TestClientInsance) Initialize() error {
	InitializeConfig()
	inst.client = zproto.NewZProtoClient("zproto2", Conf, inst)
	return nil
}

func (inst *TestClientInsance) RunLoop() {
	// go this.server.httpServer.Serve(this.server.httpListener)
	inst.client.Serve()
	// this.client.Serve()
}

func (inst *TestClientInsance) Destroy() {
	inst.client.Stop()
}

func (inst *TestClientInsance) OnNewClient(client *net2.TcpClient) {

}

func (inst *TestClientInsance) OnClientMessageArrived(client *net2.TcpClient, md *zproto.ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, msg zproto.MessageBase) error {
	return nil
}

func (inst *TestClientInsance) OnClientClosed(client *net2.TcpClient) {

}

func (inst *TestClientInsance) OnClientTimer(client *net2.TcpClient) {

}

func init() {
	flag.StringVar(&confPath, "conf", "./test_client.toml", "config path")
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
	instance := &TestClientInsance{}
	app.DoMainAppInstance(instance)
}
