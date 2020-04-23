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
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/logger"
	update2 "github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/proto/mtproto"
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"golang.org/x/net/context"
	"time"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

// https://core.telegram.org/api/end-to-end/voice-calls
//
// B accepts the call on one of their devices,
// stores the received value of g_a_hash for this instance of the voice call creation protocol,
// chooses a random value of b, 1 < b < p-1, computes g_b:=power(g,b) mod p,
// performs all the required security checks, and invokes the phone.acceptCall method,
// which has a g_b:bytes field (among others), to be filled with the value of g_b itself (not its hash).
//
// The Server S sends an updatePhoneCall with the phoneCallDiscarded constructor to all other devices B has authorized,
// to prevent accepting the same call on any of the other devices. From this point on,
// the server S works only with that of B's devices which has invoked phone.acceptCall first.
//
// phone.acceptCall#3bd2b4a0 peer:InputPhoneCall g_b:bytes protocol:PhoneCallProtocol = phone.PhoneCall;
func (s *PhoneServiceImpl) PhoneAcceptCall(ctx context.Context, request *mtproto.TLPhoneAcceptCall) (*mtproto.Phone_PhoneCall, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("phone.acceptCall#3bd2b4a0 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	//// TODO(@benqi): check peer
	peer := request.GetPeer().To_InputPhoneCall()

	callSession, err := s.PhoneCallModel.MakePhoneCallLogcByLoad(peer.GetId())
	if err != nil {
		glog.Errorf("invalid peer: {%v}, err: %v", peer, err)
		return nil, err
	}
	if peer.GetAccessHash() != callSession.ParticipantAccessHash {
		err = fmt.Errorf("invalid peer: {%v}", peer)
		glog.Errorf("invalid peer: {%v}", peer)
		return nil, err
	}

	// cache g_b
	callSession.SetGB(request.GetGB())

	/////////////////////////////////////////////////////////////////////////////////
	updatesData := update2.NewUpdatesLogic(callSession.ParticipantId)
	// 1. add updateUserStatus
	//var status *mtproto.UserStatus
	statusOnline := &mtproto.TLUserStatusOnline{Data2: &mtproto.UserStatus_Data{
		Expires: int32(time.Now().Unix() + 5*30),
	}}
	// status = statusOnline.To_UserStatus()
	updateUserStatus := &mtproto.TLUpdateUserStatus{Data2: &mtproto.Update_Data{
		UserId: callSession.ParticipantId,
		Status: statusOnline.To_UserStatus(),
	}}
	updatesData.AddUpdate(updateUserStatus.To_Update())
	// 2. add phoneCallRequested
	phoneCallAccepted := mtproto.NewTLPhoneCallAccepted()
	phoneCallAccepted.SetId(callSession.Id)
	phoneCallAccepted.SetAccessHash(callSession.ParticipantAccessHash)
	phoneCallAccepted.SetDate(int32(callSession.Date))
	phoneCallAccepted.SetAdminId(callSession.AdminId)
	phoneCallAccepted.SetParticipantId(callSession.ParticipantId)
	phoneCallAccepted.SetGB(callSession.GB)
	phoneCallAccepted.SetProtocol(callSession.ToPhoneCallProtocol())
	updatePhoneCall := &mtproto.TLUpdatePhoneCall{Data2: &mtproto.Update_Data{
		//PhoneCall: callSession.ToPhoneCallRequested().To_PhoneCall(),
		PhoneCall: phoneCallAccepted.To_PhoneCall(),
	}}
	updatesData.AddUpdate(updatePhoneCall.To_Update())
	// 3. add users
	updatesData.AddUsers(s.UserModel.GetUsersBySelfAndIDList(callSession.AdminId, []int32{md.UserId, callSession.AdminId}))
	sync_client.GetSyncClient().PushUpdates(callSession.AdminId, updatesData.ToUpdates())

	/////////////////////////////////////////////////////////////////////////////////
	// 2. reply
	phoneCall := &mtproto.TLPhonePhoneCall{Data2: &mtproto.Phone_PhoneCall_Data{
		PhoneCall: callSession.ToPhoneCallWaiting(md.UserId, 0).To_PhoneCall(),
		Users:     s.UserModel.GetUsersBySelfAndIDList(md.UserId, []int32{md.UserId, callSession.AdminId}),
	}}

	glog.Infof("phone.acceptCall#3bd2b4a0 - reply: {%v}", phoneCall)
	return phoneCall.To_Phone_PhoneCall(), nil
}
