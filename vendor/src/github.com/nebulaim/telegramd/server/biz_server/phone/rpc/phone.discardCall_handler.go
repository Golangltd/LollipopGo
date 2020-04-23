/*
 *  Copyright (c) 2017, https://github.com/nebulaim
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

package rpc

import (
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/biz/base"
	// "github.com/nebulaim/telegramd/biz/core"
	update2 "github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/proto/mtproto"
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"golang.org/x/net/context"
	// "time"
)

// phone.discardCall#78d413a6 peer:InputPhoneCall duration:int reason:PhoneCallDiscardReason connection_id:long = Updates;
func (s *PhoneServiceImpl) PhoneDiscardCall(ctx context.Context, request *mtproto.TLPhoneDiscardCall) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("phone.discardCall#78d413a6 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	//// TODO(@benqi): check peer
	peer := request.GetPeer().To_InputPhoneCall()

	callSession, err := s.PhoneCallModel.MakePhoneCallLogcByLoad(peer.GetId())
	if err != nil {
		glog.Errorf("invalid peer: {%v}, err: %v", peer, err)
		return nil, err
	}

	phoneCallDiscarded := &mtproto.TLPhoneCallDiscarded{Data2: &mtproto.PhoneCall_Data{
		Id:        callSession.Id,
		NeedDebug: true,
		Reason:    request.GetReason(),
		Duration:  request.GetDuration(),
	}}

	// var toId int32
	// = md.UserId
	// if md.UserId == callSession.AdminId {
	//	toId = callSession.ParticipantId
	// } else {
	//	toId = callSession.AdminId
	// }

	// glog.Info("toId: ", toId)

	/////////////////////////////////////////////////////////////////////////////////
	// updatesData := update2.NewUpdatesLogic(md.UserId)
	adminUpdatesData := update2.NewUpdatesLogic(callSession.AdminId)
	participantUpdatesData := update2.NewUpdatesLogic(callSession.ParticipantId)

	// 1. add phoneCallRequested
	updatePhoneCall := &mtproto.TLUpdatePhoneCall{Data2: &mtproto.Update_Data{
		PhoneCall: phoneCallDiscarded.To_PhoneCall(),
	}}
	adminUpdatesData.AddUpdate(updatePhoneCall.To_Update())
	participantUpdatesData.AddUpdate(updatePhoneCall.To_Update())

	// add message service
	action := &mtproto.TLMessageActionPhoneCall{Data2: &mtproto.MessageAction_Data{
		CallId:   callSession.Id,
		Reason:   request.GetReason(),
		Duration: request.GetDuration(),
	}}
	_ = action
	peer2 := &base.PeerUtil{
		PeerType: base.PEER_USER,
		PeerId:   callSession.ParticipantId,
	}

	_ = peer2
	// message := &mtproto.TLMessageService{Data2: &mtproto.Message_Data{
	//	Out:    true,
	//	Date:   int32(time.Now().Unix()),
	//	FromId: callSession.AdminId,
	//	ToId:   peer2.ToPeer(),
	//	Action: action.To_MessageAction(),
	//}}

	//randomId := core.GetUUID()
	//outbox := s.MessageModel.CreateMessageOutboxByNew(callSession.AdminId, peer2, randomId, message.To_Message(), func(messageId int32) {
	//	s.UserModel.CreateOrUpdateByOutbox(callSession.AdminId, peer2.PeerType, peer2.PeerId, messageId, false, false)
	//})
	//inboxList, _ := outbox.InsertMessageToInbox(callSession.AdminId, peer2, func(inBoxUserId, messageId int32) {
	//	s.UserModel.CreateOrUpdateByInbox(inBoxUserId, base.PEER_USER, peer2.PeerId, messageId, false)
	//})
	//
	//adminUpdatesData.AddUpdateNewMessage(outbox.Message)
	//participantUpdatesData.AddUpdateNewMessage(inboxList[0].Message)

	// 2. add users
	adminUpdatesData.AddUsers(s.UserModel.GetUsersBySelfAndIDList(callSession.AdminId, []int32{callSession.AdminId, callSession.ParticipantId}))
	participantUpdatesData.AddUsers(s.UserModel.GetUsersBySelfAndIDList(callSession.ParticipantId, []int32{callSession.AdminId, callSession.ParticipantId}))

	// TODO(@benqi): Add updateReadHistoryInbox
	// 3. sync
	//sync_client.GetSyncClient().PushToUserUpdatesData(toId, updatesData.ToUpdates())
	// sync_client.GetSyncClient().PushToUserUpdatesData(callSession.AdminId, adminUpdatesData.ToUpdates())
	// sync_client.GetSyncClient().PushToUserUpdatesData(callSession.ParticipantId, participantUpdatesData.ToUpdates())

	/////////////////////////////////////////////////////////////////////////////////
	replyUpdatesData := update2.NewUpdatesLogic(md.UserId)
	replyUpdatesData.AddUpdate(updatePhoneCall.To_Update())
	//
	//if md.UserId == callSession.AdminId {
	//	replyUpdatesData.AddUpdateNewMessage(outbox.Message)
	//} else {
	//	replyUpdatesData.AddUpdateNewMessage(inboxList[0].Message)
	//}
	// 2. add users
	replyUpdatesData.AddUsers(s.UserModel.GetUsersBySelfAndIDList(md.UserId, []int32{callSession.AdminId, callSession.ParticipantId}))

	glog.Infof("phone.discardCall#78d413a6 - reply {%s}", logger.JsonDebugData(replyUpdatesData))
	return replyUpdatesData.ToUpdates(), nil
}
