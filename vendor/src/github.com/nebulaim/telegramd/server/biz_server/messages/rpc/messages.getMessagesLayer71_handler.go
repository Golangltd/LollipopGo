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

package rpc

import (
    "github.com/golang/glog"
    "github.com/nebulaim/telegramd/proto/mtproto"
    "golang.org/x/net/context"
    "github.com/nebulaim/telegramd/biz/core/message"
    "github.com/nebulaim/telegramd/baselib/grpc_util"
    "github.com/nebulaim/telegramd/baselib/logger"
)

// messages.getMessages#4222fa74 id:Vector<int> = messages.Messages;
func (s *MessagesServiceImpl) MessagesGetMessagesLayer71(ctx context.Context, request *mtproto.TLMessagesGetMessagesLayer71) (*mtproto.Messages_Messages, error) {
    md := grpc_util.RpcMetadataFromIncoming(ctx)
    glog.Infof("messages.getMessages#4222fa74 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

    var idList = request.GetId()
    messages := s.MessageModel.GetUserMessagesByMessageIdList(md.UserId, idList)
    userIdList, chatIdList, _ := message.PickAllIDListByMessages(messages)
    userList := s.UserModel.GetUsersBySelfAndIDList(md.UserId, userIdList)
    chatList := s.ChatModel.GetChatListBySelfAndIDList(md.UserId, chatIdList)

    messagesMessages := &mtproto.TLMessagesMessages{Data2: &mtproto.Messages_Messages_Data{
        Messages: messages,
        Users:    userList,
        Chats:    chatList,
    }}

    glog.Infof("messages.getMessages#4222fa74 - reply: %s", logger.JsonDebugData(messagesMessages))
    return messagesMessages.To_Messages_Messages(), nil
}
