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
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/crypto"
	"github.com/nebulaim/telegramd/baselib/grpc_util/load_balancer"
	"github.com/nebulaim/telegramd/baselib/net2"
	"github.com/nebulaim/telegramd/baselib/net2/watcher2"
	"math/rand"
)

type ZProtoClientCallBack interface {
	OnNewClient(client *net2.TcpClient)
	OnClientMessageArrived(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, msg MessageBase) error
	OnClientClosed(client *net2.TcpClient)
	OnClientTimer(client *net2.TcpClient)
}

type ZProtoClientConfig struct {
	Clients []net2.ClientConfig
}

type Watcher struct {
	name    string
	watcher *watcher2.ClientWatcher
	ketama  *load_balancer.Ketama
}

type ZProtoClient struct {
	watchers []*Watcher
	clients  *net2.TcpClientGroupManager
	callback ZProtoClientCallBack
}

func NewZProtoClient(protoName string, conf *ZProtoClientConfig, cb ZProtoClientCallBack) *ZProtoClient {
	clients := map[string][]string{
		// "session": s.config.SessionClient.AddrList,
		// s.config.SessionClient.Name: s.config.SessionClient.AddrList,
	}

	c := &ZProtoClient{
		callback: cb,
	}

	c.clients = net2.NewTcpClientGroupManager(protoName, clients, c)

	// Check name
	for i := 0; i < len(conf.Clients); i++ {
		// service discovery
		etcdConfg := clientv3.Config{
			Endpoints: conf.Clients[i].EtcdAddrs,
		}
		watcher := &Watcher{
			name: conf.Clients[i].Name,
		}
		watcher.watcher, _ = watcher2.NewClientWatcher("/nebulaim", conf.Clients[i].Name, etcdConfg, c.clients)
		if conf.Clients[i].Balancer == "ketama" {
			watcher.ketama = load_balancer.NewKetama(10, nil)
		}
		c.watchers = append(c.watchers, watcher)
	}

	return c
}

///////////////////////////////////////////////////////////////////////////////////////////////
func (c *ZProtoClient) Serve() {
	for _, w := range c.watchers {
		if w.ketama != nil {
			go w.watcher.WatchClients(func(etype, addr string) {
				switch etype {
				case "add":
					w.ketama.Add(addr)
				case "delete":
					w.ketama.Remove(addr)
				}
			})
		} else {
			go w.watcher.WatchClients(nil)
		}
	}
}

func (c *ZProtoClient) Stop() {
	c.clients.Stop()
}

func (c *ZProtoClient) Pause() {
	// s.clients.Pause()
}

func (c *ZProtoClient) selectKetama(name string) *load_balancer.Ketama {
	for _, w := range c.watchers {
		if w.name == name && w.ketama != nil {
			return w.ketama
		}
	}

	return nil
}

func (c *ZProtoClient) SendKetamaMessage(name, key string, md *ZProtoMetadata, msg MessageBase, f func(addr string)) error {
	ketama := c.selectKetama(name)
	if ketama == nil {
		err := fmt.Errorf("not found ketama by name: %s", name)
		glog.Error(err)
		return err
	}

	if kaddr, ok := ketama.Get(key); ok {
		if f != nil {
			f(kaddr)
		}
		return c.SendMessageToAddress(name, kaddr, md, msg)
	} else {
		err := fmt.Errorf("not found kaddr by key: %s", key)
		glog.Error(err)
		return err
	}
}

func (c *ZProtoClient) SendMessage(name string, md *ZProtoMetadata, msg MessageBase) error {
	zmsg := &ZProtoMessage{
		Metadata: md,
		Message:  &ZProtoRawPayload{Payload: EncodeMessage(msg)},
	}
	return c.clients.SendData(name, zmsg)
}

func (c *ZProtoClient) SendMessageToAddress(name, addr string, md *ZProtoMetadata, msg MessageBase) error {
	zmsg := &ZProtoMessage{
		Metadata: md,
		Message:  &ZProtoRawPayload{Payload: EncodeMessage(msg)},
	}
	return c.clients.SendDataToAddress(name, addr, zmsg)
}

///////////////////////////////////////////////////////////////////////////////////////////////
func (c *ZProtoClient) OnNewClient(client *net2.TcpClient) {
	glog.Info("onNewClient - client: ", client, ", conn: ", client.GetConnection())
	////////////////////////////////////////////////////////////////
	// @benqi: hack
	codec := client.GetConnection().Codec()
	codec.(*ZProtoCodec).connID = client.GetConnection().GetConnID()

	// glog.Infof("onNewClient - peer(%s)", client.GetConnection())
	client.StartTimer()

	handshake := &ZProtoHandshakeReq{
		ProtoRevision: 1,
	}

	copy(handshake.RandomBytes[:], crypto.GenerateNonce(32))

	zmsg := &ZProtoMessage{
		// SessionId: client.GetConnection().GetConnID(),
		Message: handshake,
	}
	client.Send(zmsg)

	if c.callback != nil {
		c.callback.OnNewClient(client)
	}
}

func (c *ZProtoClient) OnClientDataArrived(client *net2.TcpClient, msg interface{}) error {
	zmsg, ok := msg.(*ZProtoMessage)
	if !ok {
		return fmt.Errorf("recv invalid zmsg - {%v}", zmsg)
	}

	// TODO(@benqi): check sessionId and seqNo
	switch zmsg.Message.(type) {
	case *ZProtoRawPayload:
		c.onRawPayload(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoRawPayload))
	case *ZProtoPong:
		c.onPong(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoPong))
	case *ZProtoAck:
		c.onAck(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoAck))
	case *ZProtoHandshakeRes:
		c.onHandshakeRes(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoHandshakeRes))
	case *ZProtoMarsSignal:
		c.onMarsSignal(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoMarsSignal))
	case *ZProtoDrop:
		c.onDrop(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoDrop))
	case *ZProtoRedirect:
		c.onRedirect(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoRedirect))
	case *ZProtoRpcOk:
		c.onRpcOk(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoRpcOk))
	case *ZProtoRpcError:
		c.onRpcError(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoRpcError))
	case *ZProtoRpcInternalError:
		c.onRpcInternalError(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoRpcInternalError))
	case *ZProtoRpcFloodWait:
		c.onRpcFloodWait(client, zmsg.Metadata, zmsg.sessionId, zmsg.messageId, zmsg.seqNo, zmsg.Message.(*ZProtoRpcFloodWait))
	default:
		err := fmt.Errorf("invalid message - {%v}", zmsg.Message)
		glog.Error(err)
		return err
	}

	return nil
}

func (c *ZProtoClient) OnClientClosed(client *net2.TcpClient) {
	glog.Infof("onClientClosed - peer(%s)", client.GetConnection())

	if c.callback != nil {
		c.callback.OnClientClosed(client)
	}

	if client.AutoReconnect() {
		client.Reconnect()
	}
}

func (c *ZProtoClient) OnClientTimer(client *net2.TcpClient) {
	zmsg := &ZProtoMessage{
		// SessionId: client.GetConnection().GetConnID(),
		Message: &ZProtoPing{
			PingId: rand.Int63(),
		},
	}

	glog.Info("onClientTimer - sendPing: ", zmsg)
	client.Send(zmsg)

	if c.callback != nil {
		c.callback.OnClientTimer(client)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
func (c *ZProtoClient) onRawPayload(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, payload *ZProtoRawPayload) error {
	// glog.Info("onRawPayload \n", bytes2.DumpSize(256, payload.Payload))

	var (
		err error
		m2  MessageBase
	)

	if c.callback != nil {
		m2, err = DecodeMessage(payload.Payload)
		if err != nil {
			return err
		}
		err = c.callback.OnClientMessageArrived(client, md, sessionId, messageId, seqNo, m2)
	}

	return err
}

func (c *ZProtoClient) onPong(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, pong *ZProtoPong) error {
	glog.Info("onPong: ", pong)

	return nil
}

func (c *ZProtoClient) onAck(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, ack *ZProtoAck) error {
	glog.Info("onAck: ", ack)

	return nil
}

func (c *ZProtoClient) onHandshakeRes(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, handshake *ZProtoHandshakeRes) error {
	glog.Info("onHandshakeRes: ", handshake)

	// TODO(@benqi): check handshake.ProtoRevision
	return nil
}

func (c *ZProtoClient) onDrop(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, drop *ZProtoDrop) error {
	glog.Info("onDrop: ", drop)

	// TODO(@benqi): close client
	return nil
}

func (c *ZProtoClient) onRedirect(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, redirect *ZProtoRedirect) error {
	glog.Info("onRedirect: ", redirect)

	// TODO(@benqi): close client, redirect
	return nil
}

func (c *ZProtoClient) onMarsSignal(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, marsSignal *ZProtoMarsSignal) error {
	glog.Info("onMarsSignal: ", marsSignal)

	// wechat open source mars - marsSignal support
	return nil
}

func (c *ZProtoClient) onRpcOk(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, rpcOk *ZProtoRpcOk) error {
	glog.Info("onRpcOk: ", rpcOk)

	return nil
}

func (c *ZProtoClient) onRpcError(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, rpcError *ZProtoRpcError) error {
	glog.Info("onRpcError: ", rpcError)

	return nil
}

func (c *ZProtoClient) onRpcInternalError(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, rpcError *ZProtoRpcInternalError) error {
	glog.Info("onRpcInternalError: ", rpcError)

	return nil
}

func (c *ZProtoClient) onRpcFloodWait(client *net2.TcpClient, md *ZProtoMetadata, sessionId, messageId uint64, seqNo uint32, rpcError *ZProtoRpcFloodWait) error {
	glog.Info("onRpcFloodWait: ", rpcError)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func SendMessageByClient(client *net2.TcpClient, md *ZProtoMetadata, msg MessageBase) error {
	zmsg := &ZProtoMessage{
		Metadata: md,
		Message:  &ZProtoRawPayload{Payload: EncodeMessage(msg)},
	}
	return client.Send(zmsg)
}
