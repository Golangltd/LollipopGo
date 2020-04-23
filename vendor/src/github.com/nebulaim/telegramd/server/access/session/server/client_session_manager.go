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
	"github.com/nebulaim/telegramd/baselib/queue2"
	"github.com/nebulaim/telegramd/baselib/sync2"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"sync"
	"time"
	// "github.com/nebulaim/telegramd/biz/core/user"
	"encoding/hex"
	"github.com/nebulaim/telegramd/proto/zproto"
)

const (
	kNetworkMessageStateNone             = 0 // created
	kNetworkMessageStateReceived         = 1 // received from client
	kNetworkMessageStateRunning          = 2 // invoke api
	kNetworkMessageStateWaitReplyTimeout = 3 // invoke timeout
	kNetworkMessageStateInvoked          = 4 // invoke ok, send to client
	kNetworkMessageStatePushSync         = 5 // invoke ok, send to client
	kNetworkMessageStateAck              = 6 // received client ack
	kNetworkMessageStateWaitAckTimeout   = 7 // wait ack timeout
	kNetworkMessageStateError            = 8 // invalid error
	kNetworkMessageStateEnd              = 9 // end state
)

//const (
//	kReqConnTypeTcp 			= 0
//	kReqConnTypeHttpRpc 		= 1
//	kReqConnTypeHttpWait 		= 2
//	kReqConnTypeHttpRpcAndWait 	= 3
//)

// client connID
//TRANSPORT_TCP  = 1 // TCP
//TRANSPORT_HTTP = 2 // HTTP
type ClientConnID struct {
	connType       int
	clientConnID   uint64 // client -> frontend netlib connID
	frontendConnID uint64 // frontend -> session netlib connID
	// receiveCount   int		// httpReq
	// sendCount      int		// httpRsp
}

func makeClientConnID(connType int, clientConnID, frontendConnID uint64) ClientConnID {
	connID := ClientConnID{
		connType:       connType,
		clientConnID:   clientConnID,
		frontendConnID: frontendConnID,
		// receiveCount:   0,
		// sendCount:      0,
	}
	return connID
}

func (c ClientConnID) String() string {
	return fmt.Sprintf("{conn_type: %d, client_conn_id: %d, frontend_conn_id: %d}", c.connType, c.clientConnID, c.frontendConnID)
}

func (c ClientConnID) Equal(id ClientConnID) bool {
	return c.connType == id.connType && c.clientConnID == id.clientConnID && c.frontendConnID == id.frontendConnID
}

type networkApiMessage struct {
	date       int64
	quickAckId int32 // 0: not use
	rpcRequest *mtproto.TLMessage2
	state      int // TODO(@benqi): sync.AtomicInt32
	rpcMsgId   int64
	rpcResult  mtproto.TLObject
}

type networkSyncMessage struct {
	date   int64
	update *mtproto.TLMessage2
	state  int
}

type rpcApiMessages struct {
	connID      ClientConnID
	md          *zproto.ZProtoMetadata
	sessionId   int64
	rpcMessages []*networkApiMessage
}

type sessionData struct {
	connID ClientConnID
	md     *zproto.ZProtoMetadata
	buf    []byte
}

type syncData struct {
	// sessionID int64
	md   *zproto.ZProtoMetadata
	data *messageData
}

type connData struct {
	isNew  bool
	connID ClientConnID
}

////////////////////////////////////////
const (
// inited --> work --> idle --> quit
)

type clientSessionManager struct {
	Layer           int32
	authKeyId       int64
	authKey         []byte
	AuthUserId      int32
	sessions        map[int64]*clientSessionHandler
	updatesSession  *clientUpdatesHandler
	bizRPCClient    *grpc_util.RPCClient
	nbfsRPCClient   *grpc_util.RPCClient
	syncRpcClient   mtproto.RPCSyncClient
	closeChan       chan struct{}
	sessionDataChan chan interface{} // receive from client
	rpcDataChan     chan interface{} // rpc reply
	rpcQueue        *queue2.SyncQueue
	finish          sync.WaitGroup
	running         sync2.AtomicInt32
	state           int
}

func newClientSessionManager(authKeyId int64, authKey []byte, userId int32) *clientSessionManager {
	bizRPCClient, _ := getBizRPCClient()
	nbfsRPCClient, _ := getNbfsRPCClient()
	syncRpcClient, _ := getSyncRPCClient()

	return &clientSessionManager{
		authKeyId:       authKeyId,
		authKey:         authKey,
		AuthUserId:      userId,
		sessions:        make(map[int64]*clientSessionHandler),
		updatesSession:  newClientUpdatesHandler(),
		bizRPCClient:    bizRPCClient,
		nbfsRPCClient:   nbfsRPCClient,
		syncRpcClient:   syncRpcClient,
		closeChan:       make(chan struct{}),
		sessionDataChan: make(chan interface{}, 1024),
		rpcDataChan:     make(chan interface{}, 1024),
		rpcQueue:        queue2.NewSyncQueue(),
		finish:          sync.WaitGroup{},
	}
}

func (s *clientSessionManager) String() string {
	return fmt.Sprintf("{auth_key_id: %d, user_id: %d}", s.authKeyId, s.AuthUserId)
}

func (s *clientSessionManager) Start() {
	s.running.Set(1)
	s.finish.Add(1)
	go s.rpcRunLoop()
	go s.runLoop()
}

func (s *clientSessionManager) Stop() {
	s.running.Set(0)
	s.rpcQueue.Close()
	// close(s.closeChan)
}

func (s *clientSessionManager) runLoop() {
	defer func() {
		s.finish.Done()
		close(s.closeChan)
		s.finish.Wait()
	}()

	for s.running.Get() == 1 {
		select {
		case <-s.closeChan:
			// glog.Info("runLoop -> To Close ", this.String())
			return

		case sessionMsg, _ := <-s.sessionDataChan:
			switch sessionMsg.(type) {
			case *sessionData:
				s.onSessionData(sessionMsg.(*sessionData))
			case *syncData:
				s.onSyncData(sessionMsg.(*syncData))
			case *connData:

			default:
				panic("receive invalid type msg")
			}
		case rpcMessages, _ := <-s.rpcDataChan:
			results, _ := rpcMessages.(*rpcApiMessages)
			s.onRpcResult(results)
		case <-time.After(time.Second):
			s.onTimer()
		}
	}

	glog.Info("quit runLoop...")
}

func (s *clientSessionManager) rpcRunLoop() {
	for {
		apiRequests := s.rpcQueue.Pop()
		if apiRequests == nil {
			glog.Info("quit rpcRunLoop...")
			return
		} else {
			requests, _ := apiRequests.(*rpcApiMessages)
			s.onRpcRequest(requests)
		}
	}
}

func (s *clientSessionManager) onSessionClientNew(connID ClientConnID) error {
	select {
	case s.sessionDataChan <- &connData{true, connID}:
		return nil
	}
	return nil
}

func (s *clientSessionManager) OnSessionDataArrived(connID ClientConnID, md *zproto.ZProtoMetadata, buf []byte) error {
	select {
	case s.sessionDataChan <- &sessionData{connID, md, buf}:
		return nil
	}
	return nil
}

func (s *clientSessionManager) onSessionClientClosed(connID ClientConnID) error {
	select {
	case s.sessionDataChan <- &connData{false, connID}:
		return nil
	}
	return nil
}

func (s *clientSessionManager) OnSyncDataArrived(md *zproto.ZProtoMetadata, data *messageData) error {
	select {
	case s.sessionDataChan <- &syncData{md, data}:
		return nil
	}
	return nil
}

type messageListWrapper struct {
	messages []*mtproto.TLMessage2
}

func (s *clientSessionManager) onSessionData(sessionMsg *sessionData) {
	glog.Infof("onSessionData - receive data: {sess: %s, conn_id: %s, md: %s}", s, sessionMsg.connID, sessionMsg.md)
	message := mtproto.NewEncryptedMessage2(s.authKeyId)
	err := message.Decode(s.authKeyId, s.authKey, sessionMsg.buf[8:])
	if err != nil {
		// TODO(@benqi): close frontend conn??
		glog.Error(err)
		return
	}

	glog.Infof("sessionDataChan: ", message)

	if message.MessageId&0xffffffff == 0 {
		err = fmt.Errorf("the lower 32 bits of msg_id passed by the client must not be empty: %d", message.MessageId)
		glog.Error(err)

		// TODO(@benqi): replay-attack, close client conn.
		return
	}

	sess, ok := s.sessions[message.SessionId]
	if !ok {
		sess = newClientSessionHandler(message.SessionId, message.Salt, message.MessageId, s)
	}

	if !sess.CheckBadServerSalt(sessionMsg.connID, sessionMsg.md, message.MessageId, message.SeqNo, message.Salt) {
		glog.Infof("salt invalid - {sess: %s, conn_id: %s, md: %s}", s, sessionMsg.connID, sessionMsg.md)
		// glog.Error("salt invalid..")
		return
	}

	_, isContainer := message.Object.(*mtproto.TLMsgContainer)
	if !sess.CheckBadMsgNotification(sessionMsg.connID, sessionMsg.md, message.MessageId, message.SeqNo, isContainer) {
		glog.Infof("bad msg invalid - {sess: %s, conn_id: %s, md: %s}", s, sessionMsg.connID, sessionMsg.md)
		// glog.Error("bad msg invalid..")
		return
	}

	/*
		//=============================================================================================
		// Check Message Sequence Number (msg_seqno)
		//
		// https://core.telegram.org/mtproto/description#message-sequence-number-msg-seqno
		// Message Sequence Number (msg_seqno)
		//
		// A 32-bit number equal to twice the number of “content-related” messages
		// (those requiring acknowledgment, and in particular those that are not containers)
		// created by the sender prior to this message and subsequently incremented
		// by one if the current message is a content-related message.
		// A container is always generated after its entire contents; therefore,
		// its sequence number is greater than or equal to the sequence numbers of the messages contained in it.
		//

		if message.SeqNo < sess.lastSeqNo {
			err = fmt.Errorf("sequence number is greater than or equal to the sequence numbers of the messages contained in it: %d", message.SeqNo)
			glog.Error(err)

			// TODO(@benqi): ignore this message or close client conn??
			return
		}
		sess.lastSeqNo = message.SeqNo

		sess.onMessageData(sessionMsg.md, message.MessageId, message.SeqNo, message.Object)
	*/

	var messages = &messageListWrapper{[]*mtproto.TLMessage2{}}
	extractClientMessage(message.MessageId, message.SeqNo, message.Object, messages, func(layer int32) {
		s.Layer = layer
		// TODO(@benqi): clear session_manager
	})

	if !ok {
		s.sessions[message.SessionId] = sess
		glog.Info("newClientSession: ", sess)
		sess.onNewSessionCreated(sessionMsg.connID, sessionMsg.md, message.MessageId)
		// sess.clientConnID = sessionMsg.connID
		sess.clientState = kStateOnline
	} else {
		// New Session Creation Notification
		//
		// The server notifies the client that a new session (from the server’s standpoint)
		// had to be created to handle a client message.
		// If, after this, the server receives a message with an even smaller msg_id within the same session,
		// a similar notification will be generated for this msg_id as well.
		// No such notifications are generated for high msg_id values.
		//
		if message.MessageId < sess.firstMsgId {
			glog.Info("message.MessageId < sess.firstMsgId: ", message.MessageId, ", ", sess.firstMsgId, ", sessionId: ", message.SessionId)
			sess.firstMsgId = message.MessageId
			sess.onNewSessionCreated(sessionMsg.connID, sessionMsg.md, message.MessageId)
		}
	}

	// sess.onClientMessage(message.MessageId, message.SeqNo, message.Object, messages)
	sess.onMessageData(sessionMsg.connID, sessionMsg.md, messages.messages)
}

func (s *clientSessionManager) onTimer() {
	var delList = []int64{}
	for k, v := range s.sessions {
		if !v.onTimer() {
			delList = append(delList, k)
		}
	}

	for _, id := range delList {
		delete(s.sessions, id)
	}

	if len(s.sessions) == 0 {
		deleteClientSessionManager(s.authKeyId)
	}
}

func (s *clientSessionManager) onSyncData(syncMsg *syncData) {
	glog.Infof("onSyncData - receive data: {sess: %s, md: %s, data: {%v}}",
		s, syncMsg.md, syncMsg.data)

	s.updatesSession.onSyncData(syncMsg.md, syncMsg.data.obj)
}

func (s *clientSessionManager) onConnData(connMsg *connData) {
	if connMsg.isNew {

	} else {
		s.updatesSession.UnSubscribeUpdates(connMsg.connID)
	}
}

func (s *clientSessionManager) onRpcResult(rpcResults *rpcApiMessages) {
	if sess, ok := s.sessions[rpcResults.sessionId]; ok {
		msgList := sess.pendingMessages
		sess.pendingMessages = []*pendingMessage{}
		for _, m := range rpcResults.rpcMessages {
			msgList = append(msgList, &pendingMessage{mtproto.GenerateMessageId(), true, m.rpcResult})
		}
		if len(msgList) > 0 {
			sess.sendPendingMessagesToClient(rpcResults.connID, rpcResults.md, msgList)
		}
	}
}

func (s *clientSessionManager) PushApiRequest(apiRequest *mtproto.TLMessage2) {
	s.rpcQueue.Push(apiRequest)
}

func (s *clientSessionManager) onRpcRequest(requests *rpcApiMessages) {
	glog.Infof("onRpcRequest - receive data: {sess: %s, session_id: %d, conn_id: %d, md: %s, data: {%v}}",
		s, requests.sessionId, requests.connID, requests.md, requests.rpcMessages)

	rpcMessageList := make([]*networkApiMessage, 0, len(requests.rpcMessages))

	for i := 0; i < len(requests.rpcMessages); i++ {
		var (
			err         error
			rpcResult   mtproto.TLObject
		)

		// 初始化metadata
		rpcMetadata := &grpc_util.RpcMetadata{
			ServerId:        getServerID(),
			NetlibSessionId: int64(requests.connID.clientConnID),
			AuthId:          s.authKeyId,
			SessionId:       requests.sessionId,
			TraceId:         requests.md.TraceId,
			SpanId:          getUUID(),
			ReceiveTime:     time.Now().Unix(),
			UserId:          s.AuthUserId,
			ClientMsgId:     requests.rpcMessages[i].rpcRequest.MsgId,
		}

		if s.Layer == 0 {
			s.Layer = getCacheApiLayer(s.authKeyId)
		}
		rpcMetadata.Layer = s.Layer

		// TODO(@benqi): change state.
		requests.rpcMessages[i].state = kNetworkMessageStateRunning

		// TODO(@benqi): rpc proxy
		if checkRpcUploadRequest(requests.rpcMessages[i].rpcRequest.Object) {
			rpcResult, err = s.nbfsRPCClient.Invoke(rpcMetadata, requests.rpcMessages[i].rpcRequest.Object)
		} else if checkRpcDownloadRequest(requests.rpcMessages[i].rpcRequest.Object) {
			rpcResult, err = s.nbfsRPCClient.Invoke(rpcMetadata, requests.rpcMessages[i].rpcRequest.Object)
		} else {
			rpcResult, err = s.bizRPCClient.Invoke(rpcMetadata, requests.rpcMessages[i].rpcRequest.Object)
		}

		reply := &mtproto.TLRpcResult{
			ReqMsgId: requests.rpcMessages[i].rpcRequest.MsgId,
		}

		if err != nil {
			glog.Error(err)
			rpcErr, _ := err.(*mtproto.TLRpcError)
			if rpcErr.GetErrorCode() == int32(mtproto.TLRpcErrorCodes_NOTRETURN_CLIENT) {
				continue
			}
			reply.Result = rpcErr
		} else {
			// glog.Infof("OnMessage - rpc_result: {%v}\n", rpcResult)
			reply.Result = rpcResult
		}

		requests.rpcMessages[i].state = kNetworkMessageStateInvoked
		requests.rpcMessages[i].rpcResult = reply

		rpcMessageList = append(rpcMessageList, requests.rpcMessages[i])
	}

	// TODO(@benqi): rseult metadata
	requests.rpcMessages = rpcMessageList
	s.rpcDataChan <- requests
}

// TODO(@benqi): status_client
func (s *clientSessionManager) setUserOnline(sessionId int64, connID ClientConnID) {
	defer func() {
		if r := recover(); r != nil {
			glog.Error(r)
		}
	}()

	setOnline(s.AuthUserId, s.authKeyId, getServerID(), s.Layer)
}

//==================================================================================================
type InitConnectionHandler func(layer int32)

func extractClientMessage(msgId int64, seqNo int32, object mtproto.TLObject, messages *messageListWrapper, init InitConnectionHandler) {
	switch object.(type) {
	case *mtproto.TLMsgContainer:
		msgContainer, _ := object.(*mtproto.TLMsgContainer)
		// Simple Container
		//
		// A simple container carries several messages as follows:
		//
		//  msg_container#73f1f8dc messages:vector message = MessageContainer;
		//
		// Here message refers to any message together with its length and msg_id:
		//
		//  message msg_id:long seqno:int bytes:int body:Object = Message;
		//
		// bytes is the number of bytes in the body serialization.
		//
		// All messages in a container must have msg_id lower than that of the container itself.
		// A container does not require an acknowledgment and may not carry other simple containers.
		// When messages are re-sent, they may be combined into a container in a different manner or sent individually.
		//
		// Empty containers are also allowed.
		// They are used by the server, for example,
		// to respond to an HTTP request when the timeout specified in http_wait expires,
		// and there are no messages to transmit.
		//

		// A container does not require an acknowledgment
		if seqNo%2 != 0 {
			// invalid
			// TODO(@benqi): close client and add to banned??
			glog.Error("A container does not require an acknowledgment.")
			return
		}

		// TODO(@benqi): 19: container msg_id is the same as msg_id of a previously received message (this must never happen)
		//

		for _, m := range msgContainer.Messages {
			// glog.Info("processMsgContainer - request data: ", m)
			if m.Object == nil {
				continue
			}

			// Check msgId
			//
			// A container is always generated after its entire contents; therefore,
			// its sequence number is greater than or equal to the sequence numbers of the messages contained in it.
			//
			if m.Seqno > seqNo {
				glog.Errorf("sequence number is greater than or equal to the sequence numbers of the messages contained in it: %d", seqNo)
				// TODO(@benqi): close client and add to banned??
				continue
			}

			// may not carry other simple containers
			if _, ok := m.Object.(*mtproto.TLMsgContainer); ok {
				glog.Error("may not carry other simple containers")
				// TODO(@benqi): close client and add to banned??
				continue
			}

			extractClientMessage(m.MsgId, m.Seqno, m.Object, messages, init)
		}

	case *mtproto.TLGzipPacked:
		gzipPacked, _ := object.(*mtproto.TLGzipPacked)
		glog.Info("processGzipPacked - request data: ", gzipPacked)

		dbuf := mtproto.NewDecodeBuf(gzipPacked.PackedData)
		o := dbuf.Object()
		if o == nil {
			glog.Errorf("Decode query error: %s", hex.EncodeToString(gzipPacked.PackedData))
			return
		}
		// return s.onGzipPacked(sessionId, msgId, seqNo, request)
		extractClientMessage(msgId, seqNo, o, messages, init)

	case *mtproto.TLMsgCopy:
		// not use in client
		glog.Error("android client not use msg_copy: ", object)

	case *mtproto.TLInvokeAfterMsg:
		invokeAfterMsg := object.(*mtproto.TLInvokeAfterMsg)
		invokeAfterMsgExt := NewInvokeAfterMsgExt(invokeAfterMsg)
		messages.messages = append(messages.messages, &mtproto.TLMessage2{MsgId: msgId, Seqno: seqNo, Object: invokeAfterMsgExt})

	case *mtproto.TLInvokeAfterMsgs:
		invokeAfterMsgs := object.(*mtproto.TLInvokeAfterMsgs)
		invokeAfterMsgsExt := NewInvokeAfterMsgsExt(invokeAfterMsgs)
		messages.messages = append(messages.messages, &mtproto.TLMessage2{MsgId: msgId, Seqno: seqNo, Object: invokeAfterMsgsExt})

	case *mtproto.TLInvokeWithLayer:
		invokeWithLayer := object.(*mtproto.TLInvokeWithLayer)
		//if invokeWithLayer.Layer != c.manager.Layer {
		//	c.manager.Layer = invokeWithLayer.Layer
		//	// TODO(@benqi):
		//}

		if invokeWithLayer.GetQuery() == nil {
			glog.Errorf("invokeWithLayer Query is nil, query: {%v}", invokeWithLayer)
			return
		} else {
			dbuf := mtproto.NewDecodeBuf(invokeWithLayer.Query)
			classID := dbuf.Int()
			if classID != int32(mtproto.TLConstructor_CRC32_initConnection) {
				glog.Errorf("Not initConnection classID: %d", classID)
				return
			}

			initConnection := &mtproto.TLInitConnection{}
			err := initConnection.Decode(dbuf)
			if err != nil {
				glog.Error("Decode initConnection error: ", err)
				return
			}

			// InitConnectionHandler
			init(invokeWithLayer.Layer)

			initConnectionExt := NewInitConnectionExt(initConnection)
			messages.messages = append(messages.messages, &mtproto.TLMessage2{MsgId: msgId, Seqno: seqNo, Object: initConnectionExt})
		}

	case *mtproto.TLInvokeWithoutUpdates:
		// TODO(@benqi): macOS client used.
		// glog.Error("android client not use invokeWithoutUpdates: ", object)
		invokeWithoutUpdates := object.(*mtproto.TLInvokeWithoutUpdates)
		invokeWithoutUpdatesExt := NewInvokeWithoutUpdatesExt(invokeWithoutUpdates)
		messages.messages = append(messages.messages, &mtproto.TLMessage2{MsgId: msgId, Seqno: seqNo, Object: invokeWithoutUpdatesExt})

	default:
		// glog.Info("processOthers - request data: ", object)
		messages.messages = append(messages.messages, &mtproto.TLMessage2{MsgId: msgId, Seqno: seqNo, Object: object})
	}
}
