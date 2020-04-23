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
	"golang.org/x/net/context"
	"time"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

// phone.receivedCall#17d54f61 peer:InputPhoneCall = Bool;
func (s *PhoneServiceImpl) PhoneReceivedCall(ctx context.Context, request *mtproto.TLPhoneReceivedCall) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("phone.receivedCall#17d54f61 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

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

	/////////////////////////////////////////////////////////////////////////////////
	updatesData := update2.NewUpdatesLogic(md.UserId)
	// 1. add phoneCallRequested
	updatePhoneCall := &mtproto.TLUpdatePhoneCall{Data2: &mtproto.Update_Data{
		PhoneCall: callSession.ToPhoneCallWaiting(callSession.AdminId, int32(time.Now().Unix())).To_PhoneCall(),
	}}
	updatesData.AddUpdate(updatePhoneCall.To_Update())
	// 2. add users
	updatesData.AddUsers(s.UserModel.GetUsersBySelfAndIDList(callSession.AdminId, []int32{md.UserId, callSession.AdminId}))
	// 3. sync
	sync_client.GetSyncClient().PushUpdates(callSession.AdminId, updatesData.ToUpdates())

	/////////////////////////////////////////////////////////////////////////////////
	glog.Infof("phone.receivedCall#17d54f61 - reply {true}")
	return mtproto.ToBool(true), nil
}
