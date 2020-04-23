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
	"github.com/nebulaim/telegramd/proto/mtproto"
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"golang.org/x/net/context"
	"time"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

// account.updateStatus#6628562c offline:Bool = Bool;
func (s *AccountServiceImpl) AccountUpdateStatus(ctx context.Context, request *mtproto.TLAccountUpdateStatus) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("account.updateStatus#6628562c - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var status *mtproto.UserStatus

	offline := mtproto.FromBool(request.GetOffline())
	if offline {
		// pc端：离开应用程序激活状态（点击其他应用程序）
		statusOffline := &mtproto.TLUserStatusOffline{Data2: &mtproto.UserStatus_Data{
			WasOnline: int32(time.Now().Unix()),
		}}
		status = statusOffline.To_UserStatus()
	} else {
		// pc端：客户端应用程序激活（点击客户端窗口）
		now := time.Now().Unix()
		statusOnline := &mtproto.TLUserStatusOnline{Data2: &mtproto.UserStatus_Data{
			Expires: int32(now + 5*30),
		}}
		status = statusOnline.To_UserStatus()
		s.UserModel.UpdateUserStatus(md.UserId, now)
	}

	updateUserStatus := &mtproto.TLUpdateUserStatus{Data2: &mtproto.Update_Data{
		UserId: md.UserId,
		Status: status,
	}}
	updates := &mtproto.TLUpdateShort{Data2: &mtproto.Updates_Data{
		Update: updateUserStatus.To_Update(),
		Date:   int32(time.Now().Unix()),
	}}

	// push to other contacts.
	contactIDList := s.UserModel.GetContactUserIDList(md.UserId)
	for _, id := range contactIDList {
		_ = id
		sync_client.GetSyncClient().PushUpdates(id, updates.To_Updates())
	}

	glog.Infof("account.updateStatus#6628562c - reply: {true}")
	return mtproto.ToBool(true), nil
}
