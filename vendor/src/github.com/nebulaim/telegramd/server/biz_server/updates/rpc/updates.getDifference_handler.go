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
	"golang.org/x/net/context"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
	"github.com/nebulaim/telegramd/biz/core/message"
	"fmt"
)

// updates.getDifference#25939651 flags:# pts:int pts_total_limit:flags.0?int date:int qts:int = updates.Difference;
func (s *UpdatesServiceImpl) UpdatesGetDifference(ctx context.Context, request *mtproto.TLUpdatesGetDifference) (*mtproto.Updates_Difference, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("updates.getDifference#25939651 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	difference, err := sync_client.GetSyncClient().SyncGetDifference(md.AuthId, md.UserId, request.GetPts())
	if err != nil {
		glog.Error("sync.getDifference error - ", err)
		return nil, err
	}

	switch difference.GetConstructor() {
	case mtproto.TLConstructor_CRC32_updates_differenceEmpty:
	case mtproto.TLConstructor_CRC32_updates_difference:
		diff := difference.To_UpdatesDifference()

		userIdList, chatIdList, _ := message.PickAllIDListByMessages(diff.GetNewMessages())
		userList := s.UserModel.GetUsersBySelfAndIDList(md.UserId, userIdList)
		chatList := s.ChatModel.GetChatListBySelfAndIDList(md.UserId, chatIdList)
		diff.SetUsers(userList)
		diff.SetChats(chatList)

	default:
		err = fmt.Errorf("not impl")
	}

	glog.Infof("updates.getDifference#25939651 - reply: %s", logger.JsonDebugData(difference))
	return difference, nil
}
