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
)

// messages.editChatAdmin#a9e69f2e chat_id:int user_id:InputUser is_admin:Bool = Bool;
func (s *MessagesServiceImpl) MessagesEditChatAdmin(ctx context.Context, request *mtproto.TLMessagesEditChatAdmin) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.editChatAdmin#a9e69f2e - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
		userId  int32
		isAdmin = mtproto.FromBool(request.GetIsAdmin())
		err     error
	)

	switch request.GetUserId().GetConstructor() {
	case mtproto.TLConstructor_CRC32_inputUser:
		// TODO(@benqi): check user_id's access_hash
		userId = request.GetUserId().GetData2().GetUserId()
	default:
		err = mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error("messages.editChatAdmin#a9e69f2e - invalid user_id, err: ", err)
		return nil, err
	}

	chatLogic, err := s.ChatModel.NewChatLogicById(request.ChatId)
	if err != nil {
		glog.Error("messages.editChatAdmin#a9e69f2e - error: ", err)
		return nil, err
	}

	err = chatLogic.EditChatAdmin(md.UserId, userId, isAdmin)
	if err != nil {
		glog.Error("messages.editChatAdmin#a9e69f2e - error: ", err)
		return nil, err
	}

	updateChatParticipants := &mtproto.TLUpdateChatParticipants{Data2: &mtproto.Update_Data{
		Participants: chatLogic.GetChatParticipants().To_ChatParticipants(),
	}}

	idList := chatLogic.GetChatParticipantIdList()
	for _, id := range idList {
		updates := update2.NewUpdatesLogic(md.UserId)
		updates.AddUpdate(updateChatParticipants.To_Update())
		updates.AddUsers(s.UserModel.GetUsersBySelfAndIDList(id, idList))
		updates.AddChat(chatLogic.ToChat(md.UserId))
		// sync_client.GetSyncClient().PushToUserUpdatesData(id, updates.ToUpdates())
	}

	glog.Infof("messages.editChatAdmin#a9e69f2e - reply: {true}")
	return mtproto.ToBool(true), nil
}
