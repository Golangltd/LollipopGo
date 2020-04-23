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
)

// messages.getFullChat#3b831c66 chat_id:int = messages.ChatFull;
func (s *MessagesServiceImpl) MessagesGetFullChat(ctx context.Context, request *mtproto.TLMessagesGetFullChat) (*mtproto.Messages_ChatFull, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.getFullChat#3b831c66 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): chat_id is channel

	chatLogic, err := s.ChatModel.NewChatLogicById(request.GetChatId())
	if err != nil {
		glog.Error("messages.getFullChat#3b831c66 - error: ", err)
		return nil, err
	}

	idList := chatLogic.GetChatParticipantIdList()
	messagesChatFull := &mtproto.TLMessagesChatFull{Data2: &mtproto.Messages_ChatFull_Data{
		FullChat: s.ChatModel.GetChatFullBySelfId(md.UserId, chatLogic).To_ChatFull(),
		Chats:    []*mtproto.Chat{chatLogic.ToChat(md.UserId)},
		Users:    s.UserModel.GetUsersBySelfAndIDList(md.UserId, idList),
	}}

	glog.Infof("messages.getFullChat#3b831c66 - reply: %s", logger.JsonDebugData(messagesChatFull))
	return messagesChatFull.To_Messages_ChatFull(), nil
}
