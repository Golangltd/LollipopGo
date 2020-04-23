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
	"container/list"
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/proto/zproto"
	"math/rand"
	"reflect"
	"time"
)

const (
	kDefaultPingTimeout = 30
	kPingAddTimeout     = 15
)

const (
	kStateCreated = iota
	kStateOnline
	kStateOffline
)

type messageData struct {
	confirmFlag  bool
	compressFlag bool
	obj          mtproto.TLObject
}

func (m *messageData) String() string {
	return fmt.Sprintf("{confirmFlag: %v, compressFlag: %v, obj: {%s}}", m.confirmFlag, m.compressFlag, m.obj)
}

type pendingMessage struct {
	messageId int64
	confirm   bool
	tl        mtproto.TLObject
}

func makePendingMessage(messageId int64, confirm bool, tl mtproto.TLObject) *pendingMessage {
	return &pendingMessage{messageId, confirm, tl}
}

type clientSessionHandler struct {
	closeDate        int64
	closeSessionDate int64
	salt             int64
	nextSeqNo        uint32
	sessionId        int64
	manager          *clientSessionManager
	apiMessages      *list.List
	firstMsgId       int64
	uniqueId         int64
	clientState      int
	pendingMessages  []*pendingMessage
	isUpdates        bool
	rpcMessages      []*networkApiMessage
}

func newClientSessionHandler(sessionId, salt, firstMsgId int64, m *clientSessionManager) *clientSessionHandler {
	return &clientSessionHandler{
		closeDate:        time.Now().Unix() + kDefaultPingTimeout + kPingAddTimeout,
		closeSessionDate: 0,
		salt:             salt,
		sessionId:        sessionId,
		manager:          m,
		apiMessages:      list.New(), // []*networkApiMessage{},
		firstMsgId:       firstMsgId,
		uniqueId:         rand.Int63(),
		clientState:      kStateCreated,
		pendingMessages:  []*pendingMessage{},
		isUpdates:        false,
	}
}

func (c *clientSessionHandler) String() string {
	return fmt.Sprintf("{sesses: %s, session_id: %d, client_state: %v, is_updates: %d}",
		c.manager,
		c.sessionId,
		c.clientState,
		c.isUpdates)
}

//============================================================================================
// return false, will delete this clientSession
func (c *clientSessionHandler) onTimer() bool {
	date := time.Now().Unix()

	for e := c.apiMessages.Front(); e != nil; e = e.Next() {
		if date-e.Value.(*networkApiMessage).date > 300 {
			c.apiMessages.Remove(e)
		}
	}

	//for e := c.syncMessages.Front(); e != nil; e = e.Next() {
	//	if date - e.Value.(*networkSyncMessage).date > 300 {
	//		c.apiMessages.Remove(e)
	//	}
	//}
	//
	//if date >= c.closeDate {
	//	// glog.Infof("onClose: {date: %d, c: {%v}}", date, c)
	//	c.onCloseSessionClient()
	//}
	//
	//if c.clientState == kStateOnline {
	//	for e := c.syncMessages.Front(); e != nil; e = e.Next() {
	//		v, _ := e.Value.(*networkSyncMessage)
	//		// resend
	//		if v.state != kNetworkMessageStateAck {
	//			c.sendToClient(c.clientConnID, &mtproto.ZProtoMetadata{}, v.update.MsgId, true, v.update.Object)
	//		}
	//	}
	//}
	//
	//if c.closeSessionDate != 0 && date >= c.closeSessionDate{
	//	return false
	//} else {
	//	return true
	//}
	return true
}

//============================================================================================
func (c *clientSessionHandler) encodeMessage(authKeyId int64, authKey []byte, messageId int64, confirm bool, tl mtproto.TLObject) ([]byte, error) {
	message := &mtproto.EncryptedMessage2{
		Salt:      c.salt,
		SeqNo:     c.generateMessageSeqNo(confirm),
		MessageId: messageId,
		SessionId: c.sessionId,
		Object:    tl,
	}
	return message.EncodeToLayer(authKeyId, authKey, int(c.manager.Layer))
	// return message.Encode(authKeyId, authKey)
}

func (c *clientSessionHandler) generateMessageSeqNo(increment bool) int32 {
	value := c.nextSeqNo
	if increment {
		c.nextSeqNo++
		return int32(value*2 + 1)
	} else {
		return int32(value * 2)
	}
}

func (c *clientSessionHandler) sendToClient(connID ClientConnID, md *zproto.ZProtoMetadata, messageId int64, confirm bool, obj mtproto.TLObject) error {
	// glog.Infof("sendToClient - manager: %v", c.manager)
	b, err := c.encodeMessage(c.manager.authKeyId, c.manager.authKey, messageId, confirm, obj)
	if err != nil {
		glog.Error(err)
		return err
	}

	sessData := &zproto.ZProtoSessionData{
		SessionId:  connID.frontendConnID,
		MtpRawData: b,
	}

	glog.Infof("sendSessionDataByConnID - {sess: %s, connID: %s, md: %s, sessData: %s}", c, connID, md, sessData)
	return sendSessionDataByConnID(connID.clientConnID, md, sessData)
}

func (c *clientSessionHandler) sendPendingMessagesToClient(connID ClientConnID, md *zproto.ZProtoMetadata, pendingMessages []*pendingMessage) error {
	if len(pendingMessages) == 0 {
		return nil
	}

	// glog.Infof("sendPendingMessagesToClient - connID: {%v}, pendingLen: {%v}", connID, len(pendingMessages))
	if len(pendingMessages) == 1 {
		// return c.sendToClient(connID, md, pendingMessages[0].messageId, pendingMessages[0].confirm, pendingMessages[0].tl)
		return c.sendToClient(connID, md, 0, pendingMessages[0].confirm, pendingMessages[0].tl)
	} else {
		msgContainer := &mtproto.TLMsgContainer{
			Messages: make([]mtproto.TLMessage2, 0, len(pendingMessages)),
		}
		// var seqno int32
		for _, m := range pendingMessages {
			//msgId := m.messageId
			//if msgId == 0 {
			//	msgId = mtproto.GenerateMessageId()
			//}
			message2 := mtproto.TLMessage2{
				//MsgId:  msgId,
				MsgId:  mtproto.GenerateMessageId(),
				Seqno:  c.generateMessageSeqNo(m.confirm),
				Bytes:  int32(len(m.tl.EncodeToLayer(int(c.manager.Layer)))),
				// Bytes:  int32(len(m.tl.Encode())),
				Object: m.tl,
			}
			msgContainer.Messages = append(msgContainer.Messages, message2)
		}

		return c.sendToClient(connID, md, 0, false, msgContainer)
	}
}

//// Check Server Salt
func (c *clientSessionHandler) CheckBadServerSalt(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, salt int64) bool {
	// Notice of Ignored Error Message
	//
	// Here, error_code can also take on the following values:
	//  48: incorrect server salt (in this case,
	//      the bad_server_salt response is received with the correct salt,
	//      and the message is to be re-sent with it)
	//
	if !CheckBySalt(c.manager.authKeyId, salt) {
		c.salt, _ = GetOrInsertSalt(c.manager.authKeyId)
		badServerSalt := &mtproto.TLBadServerSalt{Data2: &mtproto.BadMsgNotification_Data{
			BadMsgId:      msgId,
			ErrorCode:     48,
			BadMsgSeqno:   seqNo,
			NewServerSalt: c.salt,
		}}

		glog.Infof("invalid salt: %d, send badServerSalt: {%v}", salt, badServerSalt)
		c.sendToClient(connID, md, 0, false, badServerSalt.To_BadMsgNotification())
		return false
	}

	return true
}

func (c *clientSessionHandler) CheckBadMsgNotification(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, isContainer bool) bool {
	// Notice of Ignored Error Message
	//
	// In certain cases, a server may notify a client that its incoming message was ignored for whatever reason.
	// Note that such a notification cannot be generated unless a message is correctly decoded by the server.
	//
	// bad_msg_notification#a7eff811 bad_msg_id:long bad_msg_seqno:int error_code:int = BadMsgNotification;
	// bad_server_salt#edab447b bad_msg_id:long bad_msg_seqno:int error_code:int new_server_salt:long = BadMsgNotification;
	//
	// Here, error_code can also take on the following values:
	//
	//  16: msg_id too low (most likely, client time is wrong;
	//      it would be worthwhile to synchronize it using msg_id notifications
	//      and re-send the original message with the “correct” msg_id or wrap
	//      it in a container with a new msg_id
	//      if the original message had waited too long on the client to be transmitted)
	//  17: msg_id too high (similar to the previous case,
	//      the client time has to be synchronized, and the message re-sent with the correct msg_id)
	//  18: incorrect two lower order msg_id bits (the server expects client message msg_id to be divisible by 4)
	//  19: container msg_id is the same as msg_id of a previously received message (this must never happen)
	//  20: message too old, and it cannot be verified whether the server has received a message with this msg_id or not
	//  32: msg_seqno too low (the server has already received a message with a lower msg_id
	//      but with either a higher or an equal and odd seqno)
	//  33: msg_seqno too high (similarly, there is a message with a higher msg_id
	//      but with either a lower or an equal and odd seqno)
	//  34: an even msg_seqno expected (irrelevant message), but odd received
	//  35: odd msg_seqno expected (relevant message), but even received
	//  48: incorrect server salt (in this case,
	//      the bad_server_salt response is received with the correct salt,
	//      and the message is to be re-sent with it)
	//  64: invalid container.
	//
	// The intention is that error_code values are grouped (error_code >> 4):
	// for example, the codes 0x40 - 0x4f correspond to errors in container decomposition.
	//
	// Notifications of an ignored message do not require acknowledgment (i.e., are irrelevant).
	//
	// Important: if server_salt has changed on the server or if client time is incorrect,
	// any query will result in a notification in the above format.
	// The client must check that it has, in fact,
	// recently sent a message with the specified msg_id, and if that is the case,
	// update its time correction value (the difference between the client’s and the server’s clocks)
	// and the server salt based on msg_id and the server_salt notification,
	// so as to use these to (re)send future messages. In the meantime,
	// the original message (the one that caused the error message to be returned)
	// must also be re-sent with a better msg_id and/or server_salt.
	//
	// In addition, the client can update the server_salt value used to send messages to the server,
	// based on the values of RPC responses or containers carrying an RPC response,
	// provided that this RPC response is actually a match for the query sent recently.
	// (If there is doubt, it is best not to update since there is risk of a replay attack).
	//

	//=============================================================================================
	// TODO(@benqi): Time Synchronization, https://core.telegram.org/mtproto#time-synchronization
	//
	// Time Synchronization
	//
	// If client time diverges widely from server time,
	// a server may start ignoring client messages,
	// or vice versa, because of an invalid message identifier (which is closely related to creation time).
	// Under these circumstances,
	// the server will send the client a special message containing the correct time and
	// a certain 128-bit salt (either explicitly provided by the client in a special RPC synchronization request or
	// equal to the key of the latest message received from the client during the current session).
	// This message could be the first one in a container that includes other messages
	// (if the time discrepancy is significant but does not as yet result in the client’s messages being ignored).
	//
	// Having received such a message or a container holding it,
	// the client first performs a time synchronization (in effect,
	// simply storing the difference between the server’s time
	// and its own to be able to compute the “correct” time in the future)
	// and then verifies that the message identifiers for correctness.
	//
	// Where a correction has been neglected,
	// the client will have to generate a new session to assure the monotonicity of message identifiers.
	//

	var errorCode int32 = 0

	timeMessage := msgId / 4294967296.0
	date := time.Now().Unix()
	// glog.Info("date: ", date, ", timeMessage: ", timeMessage)

	if timeMessage+30 < date {
		errorCode = 16
	}
	if date > timeMessage+300 {
		errorCode = 17
	}

	//=================================================================================================
	// Check Message Identifier (msg_id)
	//
	// https://core.telegram.org/mtproto/description#message-identifier-msg-id
	// Message Identifier (msg_id)
	//
	// A (time-dependent) 64-bit number used uniquely to identify a message within a session.
	// Client message identifiers are divisible by 4,
	// server message identifiers modulo 4 yield 1 if the message is a response to a client message, and 3 otherwise.
	// Client message identifiers must increase monotonically (within a single session),
	// the same as server message identifiers, and must approximately equal unixtime*2^32.
	// This way, a message identifier points to the approximate moment in time the message was created.
	// A message is rejected over 300 seconds after it is created or 30 seconds
	// before it is created (this is needed to protect from replay attacks).
	// In this situation,
	// it must be re-sent with a different identifier (or placed in a container with a higher identifier).
	// The identifier of a message container must be strictly greater than those of its nested messages.
	//
	// Important: to counter replay-attacks the lower 32 bits of msg_id passed
	// by the client must not be empty and must present a fractional
	// part of the time point when the message was created.
	//
	if msgId%4 != 0 {
		errorCode = 18
	}

	// TODO(@benqi): other error code??

	if errorCode != 0 {
		badMsgNotification := &mtproto.TLBadMsgNotification{Data2: &mtproto.BadMsgNotification_Data{
			BadMsgId:    msgId,
			BadMsgSeqno: seqNo,
			ErrorCode:   errorCode,
		}}
		// glog.Info("badMsgNotification - ", badMsgNotification)
		// _ = badMsgNotification
		c.sendToClient(connID, md, 0, false, badMsgNotification.To_BadMsgNotification())
		return false
	}

	return true
}

func (c *clientSessionHandler) onNewSessionCreated(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64) {
	glog.Info("onNewSessionCreated - request data: ", msgId)
	newSessionCreated := &mtproto.TLNewSessionCreated{Data2: &mtproto.NewSession_Data{
		FirstMsgId: msgId,
		UniqueId:   c.uniqueId,
		ServerSalt: c.salt,
	}}

	// c.sendToClient(connID, md, 0, true, newSessionCreated)
	c.pendingMessages = append(c.pendingMessages, makePendingMessage(0, true, newSessionCreated))
	// TODO(@benqi): if not receive new_session_created confirm, will resend the message.
}

func (c *clientSessionHandler) onCloseSession() {
	// TODO(@benqi): remove queue???
}

func (c *clientSessionHandler) notifyMsgsStateReq() {
	// TODO(@benqi):
}

func (c *clientSessionHandler) notifyMsgsAllInfo() {
	// TODO(@benqi):
}

func (c *clientSessionHandler) notifyMsgDetailedInfo() {
	// TODO(@benqi):

	// Extended Voluntary Communication of Status of One Message
	//
	// Normally used by the server to respond to the receipt of a duplicate msg_id,
	// especially if a response to the message has already been generated and the response is large.
	// If the response is small, the server may re-send the answer itself instead.
	// This message can also be used as a notification instead of resending a large message.
	//
	// msg_detailed_info#276d3ec6 msg_id:long answer_msg_id:long bytes:int status:int = MsgDetailedInfo;
	// msg_new_detailed_info#809db6df answer_msg_id:long bytes:int status:int = MsgDetailedInfo;
	//
	// The second version is used to notify of messages that were created on the server
	// not in response to an RPC query (such as notifications of new messages)
	// and were transmitted to the client some time ago, but not acknowledged.
	//
	// Currently, status is always zero. This may change in future.
	//
	// This message does not require an acknowledgment.
	//
}

func (c *clientSessionHandler) notifyMsgResendAnsSeq() {
	// TODO(@benqi):

	// Explicit Request to Re-Send Answers
	//
	// msg_resend_ans_req#8610baeb msg_ids:Vector long = MsgResendReq;
	//
	// The remote party immediately responds by re-sending answers to the requested messages,
	// normally using the same connection that was used to transmit the query.
	// MsgsStateInfo is returned for all messages requested
	// as if the MsgResendReq query had been a MsgsStateReq query as well.
	//
}

func (c *clientSessionHandler) onMessageData(connID ClientConnID, md *zproto.ZProtoMetadata, messages []*mtproto.TLMessage2) {
	// glog.Info("onMessageData - ", messages)
	//if c.connType == UNKNOWN {
	//	connType := getConnectionType2(messages)
	//	if connType != UNKNOWN {
	//		c.connType = connType
	//	}
	//}
	//
	//if c.connType == GENERIC && c.manager.AuthUserId != 0 /* || c.connType == PUSH*/ {
	//	setUserOnline(1, connID, c.manager.authKeyId, c.sessionId, c.manager.AuthUserId)
	//
	//	//if c.manager.AuthUserId != 0 {
	//	//	for _, m := range messages {
	//	//		if !checkWithoutLogin(m.Object) {
	//	//			c.manager.AuthUserId = getCacheUserID(c.manager.authKeyId)
	//	//		}
	//	//	}
	//	//}
	//}

	var (
		hasRpcRequest bool
		hasHttpWait   bool
		ok            bool
	)

	for _, message := range messages {
		// glog.Info("onMessageData - ", message)

		if message.Object == nil {
			continue
		}

		switch message.Object.(type) {
		case *mtproto.TLRpcDropAnswer: // 所有链接都有可能
			rpcDropAnswer, _ := message.Object.(*mtproto.TLRpcDropAnswer)
			c.onRpcDropAnswer(connID, md, message.MsgId, message.Seqno, rpcDropAnswer)
		case *mtproto.TLGetFutureSalts: // GENERIC
			getFutureSalts, _ := message.Object.(*mtproto.TLGetFutureSalts)
			c.onGetFutureSalts(connID, md, message.MsgId, message.Seqno, getFutureSalts)
		case *mtproto.TLHttpWait: // 未知
			c.onHttpWait(connID, md, message.MsgId, message.Seqno, message.Object)
			hasHttpWait = true
			c.isUpdates = true
		case *mtproto.TLPing: // android未用
			ping, _ := message.Object.(*mtproto.TLPing)
			c.onPing(connID, md, message.MsgId, message.Seqno, ping)
		case *mtproto.TLPingDelayDisconnect: // PUSH和GENERIC
			ping, _ := message.Object.(*mtproto.TLPingDelayDisconnect)
			c.onPingDelayDisconnect(connID, md, message.MsgId, message.Seqno, ping)
		case *mtproto.TLDestroySession: // GENERIC
			destroySession, _ := message.Object.(*mtproto.TLDestroySession)
			c.onDestroySession(connID, md, message.MsgId, message.Seqno, destroySession)
		case *mtproto.TLMsgsAck: // 所有链接都有可能
			msgsAck, _ := message.Object.(*mtproto.TLMsgsAck)
			c.onMsgsAck(connID, md, message.MsgId, message.Seqno, msgsAck)
			// TODO(@benqi): check c.isUpdates
		case *mtproto.TLMsgsStateReq: // android未用
			c.onMsgsStateReq(connID, md, message.MsgId, message.Seqno, message.Object)
		case *mtproto.TLMsgsStateInfo: // android未用
			c.onMsgsStateInfo(connID, md, message.MsgId, message.Seqno, message.Object)
		case *mtproto.TLMsgsAllInfo: // android未用
			c.onMsgsAllInfo(connID, md, message.MsgId, message.Seqno, message.Object)
		case *mtproto.TLMsgResendReq: // 都有可能
			c.onMsgResendReq(connID, md, message.MsgId, message.Seqno, message.Object)
		case *mtproto.TLMsgDetailedInfo: // 都有可能
			// glog.Error("client side msg: ", object)
		case *mtproto.TLMsgNewDetailedInfo: // 都有可能
			// glog.Error("client side msg: ", object)
		case *mtproto.TLContestSaveDeveloperInfo: // 未知
			contestSaveDeveloperInfo, _ := message.Object.(*mtproto.TLContestSaveDeveloperInfo)
			c.onContestSaveDeveloperInfo(connID, md, message.MsgId, message.Seqno, contestSaveDeveloperInfo)
		case *TLInvokeAfterMsgExt: // 未知
			invokeAfterMsgExt, _ := message.Object.(*TLInvokeAfterMsgExt)
			// c.onRpcRequest(md, message.MsgId, message.Seqno, invokeAfterMsgExt.Query)
			ok = c.onInvokeAfterMsgExt(connID, md, message.MsgId, message.Seqno, invokeAfterMsgExt)
			if ok && !hasRpcRequest {
				hasRpcRequest = ok
			}
		case *TLInvokeAfterMsgsExt: // 未知
			invokeAfterMsgsExt, _ := message.Object.(*TLInvokeAfterMsgsExt)
			// c.onRpcRequest(md, message.MsgId, message.Seqno, invokeAfterMsgsExt.Query)
			ok = c.onInvokeAfterMsgsExt(connID, md, message.MsgId, message.Seqno, invokeAfterMsgsExt)
			if ok && !hasRpcRequest {
				hasRpcRequest = ok
			}
		case *TLInitConnectionExt: // 都有可能
			initConnectionExt, _ := message.Object.(*TLInitConnectionExt)
			ok = c.onInitConnectionEx(connID, md, message.MsgId, message.Seqno, initConnectionExt)
			if ok && !hasRpcRequest {
				hasRpcRequest = ok
			}
		case *TLInvokeWithoutUpdatesExt:
			invokeWithoutUpdatesExt, _ := message.Object.(*TLInvokeWithoutUpdatesExt)
			ok = c.onInvokeWithoutUpdatesExt(connID, md, message.MsgId, message.Seqno, invokeWithoutUpdatesExt)
			if ok && !hasRpcRequest {
				hasRpcRequest = ok
			}
		default:
			ok = c.onRpcRequest(connID, md, message.MsgId, message.Seqno, message.Object)
			if ok && !hasRpcRequest {
				hasRpcRequest = ok
			}
		}
	}

	//if hasHttpWait {
	//	if !hasRpcRequest {
	//		if len(c.pendingMessages) > 0 {
	//			// no rpc request and has pendingMessages, send to client
	//		} else {
	//			// http_wait
	//		}
	//	} else {
	//		// receive wait rpc result and send pending
	//	}
	//} else {
	//
	//}
	//	if !hasRpcRequest && len(c.pendingMessages) > 0 {
	//	} else {
	//
	//	}
	//}

	if c.isUpdates {
		c.manager.setUserOnline(c.sessionId, connID)
	}
	// c.isUpdates = true
	// subscribe
	// c.manager.updatesSession.SubscribeUpdates(c, connID)

	if connID.connType == mtproto.TRANSPORT_TCP {
		if c.isUpdates {
			c.manager.updatesSession.SubscribeUpdates(c, connID)
			// c.manager.setUserOnline(c.sessionId, connID)
			// c.manager.updatesSession.SubscribeUpdates(c, connID)
		}
		c.sendPendingMessagesToClient(connID, md, c.pendingMessages)
		c.pendingMessages = []*pendingMessage{}
	} else {
		if !hasRpcRequest {
			if len(c.pendingMessages) > 0 {
				c.sendPendingMessagesToClient(connID, md, c.pendingMessages)
				c.pendingMessages = []*pendingMessage{}
			} else {
				c.manager.updatesSession.SubscribeUpdates(c, connID)
				//if !hasHttpWait {
				//	// TODO(@benqi): close http
				//} else {
				//	// c.manager.setUserOnline(c.sessionId, connID)
				//	c.manager.updatesSession.SubscribeUpdates(c, connID)
				//}
			}
		} else {
			// wait
		}
	}

	_ = hasHttpWait

	if len(c.rpcMessages) > 0 {
		c.manager.rpcQueue.Push(&rpcApiMessages{connID: connID, md: md, sessionId: c.sessionId, rpcMessages: c.rpcMessages})
		c.rpcMessages = []*networkApiMessage{}
	}
}

//============================================================================================
func (c *clientSessionHandler) onPing(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, ping *mtproto.TLPing) {
	glog.Infof("onPing - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		logger.JsonDebugData(ping))

	pong := &mtproto.TLPong{Data2: &mtproto.Pong_Data{
		MsgId:  msgId,
		PingId: ping.PingId,
	}}

	// c.sendToClient(connID, md, 0, false, pong)
	c.pendingMessages = append(c.pendingMessages, makePendingMessage(0, false, pong))

	c.closeDate = time.Now().Unix() + kDefaultPingTimeout + kPingAddTimeout
}

func (c *clientSessionHandler) onPingDelayDisconnect(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, pingDelayDisconnect *mtproto.TLPingDelayDisconnect) {
	glog.Infof("onPing - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		logger.JsonDebugData(pingDelayDisconnect))

	pong := &mtproto.TLPong{Data2: &mtproto.Pong_Data{
		MsgId:  msgId,
		PingId: pingDelayDisconnect.PingId,
	}}

	// c.sendToClient(connID, md, 0, false, pong)
	c.pendingMessages = append(c.pendingMessages, makePendingMessage(0, false, pong))

	c.closeDate = time.Now().Unix() + int64(pingDelayDisconnect.DisconnectDelay) + kPingAddTimeout
}

func (c *clientSessionHandler) onMsgsAck(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request *mtproto.TLMsgsAck) {
	glog.Infof("onMsgsAck - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		logger.JsonDebugData(request))

	for _, id := range request.GetMsgIds() {
		// reqMsgId := msgId
		for e := c.apiMessages.Front(); e != nil; e = e.Next() {
			v, _ := e.Value.(*networkApiMessage)
			if v.rpcMsgId == id {
				v.state = kNetworkMessageStateAck
				glog.Info("onMsgsAck - networkSyncMessage change kNetworkMessageStateAck")
			}
		}

		//for e := c.syncMessages.Front(); e != nil; e = e.Next() {
		//	v, _ := e.Value.(*networkSyncMessage)
		//	if v.update.MsgId == id {
		//		v.state = kNetworkMessageStateAck
		//		glog.Info("onMsgsAck - networkSyncMessage change kNetworkMessageStateAck")
		//		// TODO(@benqi): update pts, qts, seq etc...
		//	}
		//}
	}
}

func (c *clientSessionHandler) onHttpWait(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request mtproto.TLObject) {
	glog.Infof("onHttpWait - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		logger.JsonDebugData(request))

	c.isUpdates = true
	// c.manager.setUserOnline(c.sessionId, connID)
	// c.manager.updatesSession.SubscribeHttpUpdates(c, connID)
}

func (c *clientSessionHandler) onMsgsStateReq(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request mtproto.TLObject) {
	glog.Infof("onMsgsStateReq - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		logger.JsonDebugData(request))

	// Request for Message Status Information
	//
	// If either party has not received information on the status of its outgoing messages for a while,
	// it may explicitly request it from the other party:
	//
	// msgs_state_req#da69fb52 msg_ids:Vector long = MsgsStateReq;
	// The response to the query contains the following information:
	//
	// Informational Message regarding Status of Messages
	// msgs_state_info#04deb57d req_msg_id:long info:string = MsgsStateInfo;
	//
	// Here, info is a string that contains exactly one byte of message status for
	// each message from the incoming msg_ids list:
	//
	// 1 = nothing is known about the message (msg_id too low, the other party may have forgotten it)
	// 2 = message not received (msg_id falls within the range of stored identifiers; however,
	// 	   the other party has certainly not received a message like that)
	// 3 = message not received (msg_id too high; however, the other party has certainly not received it yet)
	// 4 = message received (note that this response is also at the same time a receipt acknowledgment)
	// +8 = message already acknowledged
	// +16 = message not requiring acknowledgment
	// +32 = RPC query contained in message being processed or processing already complete
	// +64 = content-related response to message already generated
	// +128 = other party knows for a fact that message is already received
	//
	// This response does not require an acknowledgment.
	// It is an acknowledgment of the relevant msgs_state_req, in and of itself.
	//
	// Note that if it turns out suddenly that the other party does not have a message
	// that looks like it has been sent to it, the message can simply be re-sent.
	// Even if the other party should receive two copies of the message at the same time,
	// the duplicate will be ignored. (If too much time has passed,
	// and the original msg_id is not longer valid, the message is to be wrapped in msg_copy).
	//
}

func (c *clientSessionHandler) onInitConnectionEx(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request *TLInitConnectionExt) bool {
	glog.Infof("onInitConnection - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)
	// glog.Infof("onInitConnection - request: %s", request.String())
	// auth_session_client.BindAuthKeyUser()
	uploadInitConnection(c.manager.authKeyId, c.manager.Layer, md.ClientAddr, request)
	return c.onRpcRequest(connID, md, msgId, seqNo, request.Query)
}

func (c *clientSessionHandler) onMsgResendReq(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request mtproto.TLObject) {
	glog.Infof("onMsgResendReq - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)

	// Explicit Request to Re-Send Messages
	//
	// msg_resend_req#7d861a08 msg_ids:Vector long = MsgResendReq;
	//
	// The remote party immediately responds by re-sending the requested messages,
	// normally using the same connection that was used to transmit the query.
	// If at least one message with requested msg_id does not exist or has already been forgotten,
	// or has been sent by the requesting party (known from parity),
	// MsgsStateInfo is returned for all messages requested
	// as if the MsgResendReq query had been a MsgsStateReq query as well.
	//
}

func (c *clientSessionHandler) onMsgsStateInfo(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request mtproto.TLObject) {
	glog.Infof("onMsgsStateInfo - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)
}

func (c *clientSessionHandler) onMsgsAllInfo(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request mtproto.TLObject) {
	glog.Infof("onMsgsAllInfo - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)

	// Voluntary Communication of Status of Messages
	//
	// Either party may voluntarily inform the other party of the status of the messages transmitted by the other party.
	//
	// msgs_all_info#8cc0d131 msg_ids:Vector long info:string = MsgsAllInfo
	//
	// All message codes known to this party are enumerated,
	// with the exception of those for which the +128 and the +16 flags are set.
	// However, if the +32 flag is set but not +64, then the message status will still be communicated.
	//
	// This message does not require an acknowledgment.
	//
}

func (c *clientSessionHandler) onDestroySession(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request *mtproto.TLDestroySession) {
	glog.Infof("onDestroySession - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)

	// Request to Destroy Session
	//
	// Used by the client to notify the server that it may
	// forget the data from a different session belonging to the same user (i. e. with the same auth_key_id).
	// The result of this being applied to the current session is undefined.
	//
	// destroy_session#e7512126 session_id:long = DestroySessionRes;
	// destroy_session_ok#e22045fc session_id:long = DestroySessionRes;
	// destroy_session_none#62d350c9 session_id:long = DestroySessionRes;
	//

	if request.SessionId == c.sessionId {
		// The result of this being applied to the current session is undefined.
		glog.Error("the result of this being applied to the current session is undefined.")

		// TODO(@benqi): handle error???
		return
	}

	if _, ok := c.manager.sessions[request.SessionId]; ok {
		destroySessionOk := &mtproto.TLDestroySessionOk{Data2: &mtproto.DestroySessionRes_Data{
			SessionId: request.SessionId,
		}}
		// c.sendToClient(connID, md, 0, false, destroySessionOk.To_DestroySessionRes())
		c.pendingMessages = append(c.pendingMessages, makePendingMessage(0, false, destroySessionOk.To_DestroySessionRes()))

		delete(c.manager.sessions, request.SessionId)

		// TODO(@benqi): saved destroyed session???
	} else {
		destroySessionNone := &mtproto.TLDestroySessionOk{Data2: &mtproto.DestroySessionRes_Data{
			SessionId: request.SessionId,
		}}
		// c.sendToClient(connID, md, 0, false, destroySessionNone.To_DestroySessionRes())
		c.pendingMessages = append(c.pendingMessages, makePendingMessage(0, false, destroySessionNone.To_DestroySessionRes()))
	}
}

func (c *clientSessionHandler) onGetFutureSalts(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request *mtproto.TLGetFutureSalts) {
	glog.Infof("onGetFutureSalts - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)

	salts, _ := GetOrInsertSaltList(c.manager.authKeyId, int(request.Num))
	futureSalts := &mtproto.TLFutureSalts{Data2: &mtproto.FutureSalts_Data{
		ReqMsgId: msgId,
		Now:      int32(time.Now().Unix()),
		Salts:    salts,
	}}

	glog.Info("onGetFutureSalts - reply data: ", futureSalts)
	// c.sendToClient(connID, md, 0, false, futureSalts)
	c.pendingMessages = append(c.pendingMessages, makePendingMessage(0, false, futureSalts))
}

// sendToClient:
// 	rpc_answer_unknown#5e2ad36e = RpcDropAnswer;
// 	rpc_answer_dropped_running#cd78e586 = RpcDropAnswer;
// 	rpc_answer_dropped#a43ad8b7 msg_id:long seq_no:int bytes:int = RpcDropAnswer;
func (c *clientSessionHandler) onRpcDropAnswer(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request *mtproto.TLRpcDropAnswer) {
	glog.Infof("onRpcDropAnswer - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%v}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)

	rpcAnswer := &mtproto.RpcDropAnswer{Data2: &mtproto.RpcDropAnswer_Data{}}

	var found = false
	for e := c.apiMessages.Front(); e != nil; e = e.Next() {
		v, _ := e.Value.(*networkApiMessage)
		if v.rpcRequest.MsgId == request.ReqMsgId {
			if v.state == kNetworkMessageStateReceived {
				rpcAnswer.Constructor = mtproto.TLConstructor_CRC32_rpc_answer_dropped
				rpcAnswer.Data2.MsgId = request.ReqMsgId
				// TODO(@benqi): set seqno and bytes
				// rpcAnswer.Data2.SeqNo = 0
				// rpcAnswer.Data2.Bytes = 0
			} else if v.state == kNetworkMessageStateInvoked {
				rpcAnswer.Constructor = mtproto.TLConstructor_CRC32_rpc_answer_dropped_running
			} else {
				rpcAnswer.Constructor = mtproto.TLConstructor_CRC32_rpc_answer_unknown
			}
			found = true
			break
		}
	}

	if !found {
		rpcAnswer.Constructor = mtproto.TLConstructor_CRC32_rpc_answer_unknown
	}

	// android client code:
	/*
		 if (notifyServer) {
			TL_rpc_drop_answer *dropAnswer = new TL_rpc_drop_answer();
			dropAnswer->req_msg_id = request->messageId;
			sendRequest(dropAnswer, nullptr, nullptr, RequestFlagEnableUnauthorized | RequestFlagWithoutLogin | RequestFlagFailOnServerErrors, request->datacenterId, request->connectionType, true);
		 }
	*/

	// and both of these responses require an acknowledgment from the client.
	// c.sendToClient(connID, md, 0, true, &mtproto.TLRpcResult{ReqMsgId: msgId, Result: rpcAnswer})
	c.pendingMessages = append(c.pendingMessages, makePendingMessage(0, true, &mtproto.TLRpcResult{ReqMsgId: msgId, Result: rpcAnswer}))

}

func (c *clientSessionHandler) onContestSaveDeveloperInfo(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request *mtproto.TLContestSaveDeveloperInfo) {
	// contestSaveDeveloperInfo, _ := request.(*mtproto.TLContestSaveDeveloperInfo)
	glog.Infof("onContestSaveDeveloperInfo - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%v}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)

	// TODO(@benqi): 实现scontestSaveDeveloperInfo处理逻辑
	// r := &mtproto.TLTrue{}
	// c.sendToClient(connID, md, false, &mtproto.TLTrue{})

	// _ = r
}

func (c *clientSessionHandler) onInvokeAfterMsgExt(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request *TLInvokeAfterMsgExt) bool {
	glog.Infof("onInvokeAfterMsgExt - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%v}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)

	//		if invokeAfterMsg.GetQuery() == nil {
	//			glog.Errorf("invokeAfterMsg Query is nil, query: {%v}", invokeAfterMsg)
	//			return
	//		}
	//
	//		dbuf := mtproto.NewDecodeBuf(invokeAfterMsg.Query)
	//		query := dbuf.Object()
	//		if query == nil {
	//			glog.Errorf("Decode query error: %s", hex.EncodeToString(invokeAfterMsg.Query))
	//			return
	//		}
	//
	//		var found = false
	//		for j := 0; j < i; j++ {
	//			if messages[j].MsgId == invokeAfterMsg.MsgId {
	//				messages[i].Object = query
	//				found = true
	//				break
	//			}
	//		}
	//
	//		if !found {
	//			for j := i + 1; j < len(messages); j++ {
	//				if messages[j].MsgId == invokeAfterMsg.MsgId {
	//					// c.messages[i].Object = query
	//					messages[i].Object = query
	//					found = true
	//					messages = append(messages, messages[i])
	//
	//					// set messages[i] = nil, will ignore this.
	//					messages[i] = nil
	//					break
	//				}
	//			}
	//		}
	//
	//		if !found {
	//			// TODO(@benqi): backup message, wait.
	//
	//			messages[i].Object = query
	//		}

	return c.onRpcRequest(connID, md, msgId, seqNo, request.Query)
}

func (c *clientSessionHandler) onInvokeAfterMsgsExt(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request *TLInvokeAfterMsgsExt) bool {
	glog.Infof("onInvokeAfterMsgsExt - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%v}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		request)
	//		if invokeAfterMsgs.GetQuery() == nil {
	//			glog.Errorf("invokeAfterMsgs Query is nil, query: {%v}", invokeAfterMsgs)
	//			return
	//		}
	//
	//		dbuf := mtproto.NewDecodeBuf(invokeAfterMsgs.Query)
	//		query := dbuf.Object()
	//		if query == nil {
	//			glog.Errorf("Decode query error: %s", hex.EncodeToString(invokeAfterMsgs.Query))
	//			return
	//		}
	//
	//		if len(invokeAfterMsgs.MsgIds) == 0 {
	//			// TODO(@benqi): invalid msgIds, ignore??
	//
	//			messages[i].Object = query
	//		} else {
	//			var maxMsgId = invokeAfterMsgs.MsgIds[0]
	//			for j := 1; j < len(invokeAfterMsgs.MsgIds); j++ {
	//				if maxMsgId > invokeAfterMsgs.MsgIds[j] {
	//					maxMsgId = invokeAfterMsgs.MsgIds[j]
	//				}
	//			}
	//
	//
	//			var found = false
	//			for j := 0; j < i; j++ {
	//				if messages[j].MsgId == maxMsgId {
	//					messages[i].Object = query
	//					found = true
	//					break
	//				}
	//			}
	//
	//			if !found {
	//				for j := i + 1; j < len(messages); j++ {
	//					if messages[j].MsgId == maxMsgId {
	//						// c.messages[i].Object = query
	//						messages[i].Object = query
	//						found = true
	//						messages = append(messages, messages[i])
	//
	//						// set messages[i] = nil, will ignore this.
	//						messages[i] = nil
	//						break
	//					}
	//				}
	//			}
	//
	//			if !found {
	//				// TODO(@benqi): backup message, wait.
	//
	//				messages[i].Object = query
	//			}

	return c.onRpcRequest(connID, md, msgId, seqNo, request.Query)
}

func (c *clientSessionHandler) onInvokeWithoutUpdatesExt(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, request *TLInvokeWithoutUpdatesExt) bool {
	glog.Infof("onInvokeWithoutUpdatesExt - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		reflect.TypeOf(request))

	return c.onRpcRequest(connID, md, msgId, seqNo, request.Query)
}

func (c *clientSessionHandler) onRpcRequest(connID ClientConnID, md *zproto.ZProtoMetadata, msgId int64, seqNo int32, object mtproto.TLObject) bool {
	glog.Infof("onRpcRequest - request data: {sess: %s, conn_id: %s, md: %s, msg_id: %d, seq_no: %d, request: {%s}}",
		c,
		connID,
		md,
		msgId,
		seqNo,
		reflect.TypeOf(object))

	// TODO(@benqi): sync AuthUserId??
	requestMessage := &mtproto.TLMessage2{
		MsgId:  msgId,
		Seqno:  seqNo,
		Object: object,
	}

	// reqMsgId := msgId
	for e := c.apiMessages.Front(); e != nil; e = e.Next() {
		//v, _ := e.Value.(*networkApiMessage)
		//if v.rpcRequest.MsgId == msgId {
		//	if v.state >= kNetworkMessageStateInvoked {
		//		// c.pendingMessages = append(c.pendingMessages, makePendingMessage(v.rpcMsgId, true, v.rpcResult))
		//		return false
		//	}
		//}
	}

	if c.manager.AuthUserId == 0 {
		if !checkRpcWithoutLogin(object) {
			c.manager.AuthUserId = getCacheUserID(c.manager.authKeyId)
			if c.manager.AuthUserId == 0 {
				glog.Error("not found authUserId by authKeyId: ", c.manager.authKeyId)
				return false
			}
		}
	}

	// updates
	if checkRpcUpdatesType(object) {
		// c.manager.setUserOnline(c.sessionId, connID)
		glog.Infof("onRpcRequest - isUpdate: {connID: {%v}, rpcMethod: {%T}}", connID, object)
		c.isUpdates = true
		// c.manager.updatesSession.SubscribeUpdates(c, connID)

		// subscribe
		// c.manager.updatesSession.SubscribeUpdates(c, connID)
	}

	apiMessage := &networkApiMessage{
		date:       time.Now().Unix(),
		rpcRequest: requestMessage,
		state:      kNetworkMessageStateReceived,
	}
	glog.Info("onRpcRequest - ", apiMessage)
	// c.apiMessages = append(c.apiMessages, apiMessage)
	c.apiMessages.PushBack(apiMessage)

	c.rpcMessages = append(c.rpcMessages, apiMessage)
	// c.manager.rpcQueue.Push(&rpcApiMessage{connID: connID, sessionId: c.sessionId, rpcMessage: apiMessage})

	return true
}

// 客户端连接事件
func (c *clientSessionHandler) onSessionClientConnected() {
	//c.clientSession = &clientSession{conn, sessionID}
	if c.clientState == kStateOffline {
		glog.Infof("onSessionClientConnected: ", c)
		c.clientState = kStateOnline
		c.closeSessionDate = 0
		c.closeDate = time.Now().Unix() + kDefaultPingTimeout + kPingAddTimeout
		//if c.synced && c.connType == GENERIC{
		//	// TODO(@benqi): push sync data
		//	syncReq := &mtproto.NewUpdatesRequest{
		//		AuthKeyId: c.manager.authKeyId,
		//		UserId:    c.manager.AuthUserId,
		//	}
		//
		//	updates, err := c.manager.syncRpcClient.GetNewUpdatesData(context.Background(), syncReq)
		//	if err != nil {
		//		glog.Error(err)
		//		// return nil, false
		//	} else {
		//		glog.Info("getNewUpdatesData: ", updates)
		//		if len(updates.GetData2().Updates) > 0 {
		//			c.onSyncData(c.clientConnID, &mtproto.ZProtoMetadata{}, updates)
		//		}
		//	}
		//}
	}
}

func (c *clientSessionHandler) onCloseSessionClient() {
	if c.clientState == kStateOnline {
		glog.Infof("onCloseSessionClient: ", c)
		c.clientState = kStateOffline
		c.closeSessionDate = time.Now().Unix() + 3600
	}
}
