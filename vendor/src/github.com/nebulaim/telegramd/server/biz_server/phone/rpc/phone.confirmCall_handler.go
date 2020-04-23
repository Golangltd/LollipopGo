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
	update2 "github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/proto/mtproto"
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"golang.org/x/net/context"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

// phone.confirmCall#2efe1722 peer:InputPhoneCall g_a:bytes key_fingerprint:long protocol:PhoneCallProtocol = phone.PhoneCall;
func (s *PhoneServiceImpl) PhoneConfirmCall(ctx context.Context, request *mtproto.TLPhoneConfirmCall) (*mtproto.Phone_PhoneCall, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("phone.confirmCall#2efe1722 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	//// TODO(@benqi): check peer
	peer := request.GetPeer().To_InputPhoneCall()

	callSession, err := s.PhoneCallModel.MakePhoneCallLogcByLoad(peer.GetId())
	if err != nil {
		glog.Errorf("invalid peer: {%v}, err: %v", peer, err)
		return nil, err
	}
	// if peer.GetAccessHash() != callSession.AdminAccessHash {
	// 	err = fmt.Errorf("invalid peer: {%v}", peer)
	// 	glog.Errorf("invalid peer: {%v}", peer)
	// 	return nil, err
	// }

	// TODO(@benqi): callSession.SetGA() ???
	callSession.GA = request.GetGA()

	/////////////////////////////////////////////////////////////////////////////////
	updatesData := update2.NewUpdatesLogic(md.UserId)
	// 1. add phoneCallRequested
	updatePhoneCall := &mtproto.TLUpdatePhoneCall{Data2: &mtproto.Update_Data{
		PhoneCall: callSession.ToPhoneCall(callSession.ParticipantId, request.GetKeyFingerprint(), s.RelayIp).To_PhoneCall(),
	}}
	updatesData.AddUpdate(updatePhoneCall.To_Update())
	// 2. add users
	updatesData.AddUsers(s.UserModel.GetUsersBySelfAndIDList(callSession.ParticipantId, []int32{md.UserId, callSession.ParticipantId}))
	// 3. sync
	sync_client.GetSyncClient().PushUpdates(callSession.ParticipantId, updatesData.ToUpdates())

	/////////////////////////////////////////////////////////////////////////////////
	// 2. reply
	phoneCall := &mtproto.TLPhonePhoneCall{Data2: &mtproto.Phone_PhoneCall_Data{
		PhoneCall: callSession.ToPhoneCall(md.UserId, request.GetKeyFingerprint(), s.RelayIp).To_PhoneCall(),
		Users:     s.UserModel.GetUsersBySelfAndIDList(md.UserId, []int32{md.UserId, callSession.ParticipantId}),
	}}

	glog.Infof("phone.confirmCall#2efe1722 - reply: {%v}", phoneCall)
	return phoneCall.To_Phone_PhoneCall(), nil
}
