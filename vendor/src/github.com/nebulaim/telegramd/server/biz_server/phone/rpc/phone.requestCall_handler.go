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
	"golang.org/x/net/context"
	"time"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

// phone.requestCall#5b95b3d4 user_id:InputUser random_id:int g_a_hash:bytes protocol:PhoneCallProtocol = phone.PhoneCall;
func (s *PhoneServiceImpl) PhoneRequestCall(ctx context.Context, request *mtproto.TLPhoneRequestCall) (*mtproto.Phone_PhoneCall, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("phone.requestCall#5b95b3d4 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
		err           error
		participantId int32
	)

	switch request.GetUserId().GetConstructor() {
	case mtproto.TLConstructor_CRC32_inputUser:
		// TODO(@benqi): Check access_hash
		participantId = request.GetUserId().GetData2().GetUserId()
	default:
		err = mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error("inputUser is empty or self, err: ", err)
		return nil, err
	}

	callSession := s.PhoneCallModel.NewPhoneCallLogic(md.UserId, participantId, request.GetGAHash(), request.GetProtocol().To_PhoneCallProtocol())

	/////////////////////////////////////////////////////////////////////////////////
	updatesData := update2.NewUpdatesLogic(md.UserId)
	// 1. add updateUserStatus
	//var status *mtproto.UserStatus
	statusOnline := &mtproto.TLUserStatusOnline{Data2: &mtproto.UserStatus_Data{
		Expires: int32(time.Now().Unix() + 5*30),
	}}
	// status = statusOnline.To_UserStatus()
	updateUserStatus := &mtproto.TLUpdateUserStatus{Data2: &mtproto.Update_Data{
		UserId: md.UserId,
		Status: statusOnline.To_UserStatus(),
	}}
	updatesData.AddUpdate(updateUserStatus.To_Update())
	// 2. add phoneCallRequested
	updatePhoneCall := &mtproto.TLUpdatePhoneCall{Data2: &mtproto.Update_Data{
		PhoneCall: callSession.ToPhoneCallRequested().To_PhoneCall(),
	}}
	updatesData.AddUpdate(updatePhoneCall.To_Update())
	// 3. add users
	updatesData.AddUsers(s.UserModel.GetUsersBySelfAndIDList(participantId, []int32{md.UserId, participantId}))
	sync_client.GetSyncClient().PushUpdates(participantId, updatesData.ToUpdates())

	/////////////////////////////////////////////////////////////////////////////////
	// 2. reply
	phoneCall := &mtproto.TLPhonePhoneCall{Data2: &mtproto.Phone_PhoneCall_Data{
		PhoneCall: callSession.ToPhoneCallWaiting(md.UserId, 0).To_PhoneCall(),
		Users:     s.UserModel.GetUsersBySelfAndIDList(md.UserId, []int32{md.UserId, participantId}),
	}}

	glog.Infof("phone.requestCall#5b95b3d4 - reply: {%v}", phoneCall)
	return phoneCall.To_Phone_PhoneCall(), nil
}
