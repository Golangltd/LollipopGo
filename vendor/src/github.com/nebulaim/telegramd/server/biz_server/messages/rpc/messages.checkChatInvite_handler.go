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

// messages.checkChatInvite#3eadb1bb hash:string = ChatInvite;
func (s *MessagesServiceImpl) MessagesCheckChatInvite(ctx context.Context, request *mtproto.TLMessagesCheckChatInvite) (*mtproto.ChatInvite, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.checkChatInvite#3eadb1bb - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
		chatInvite *mtproto.ChatInvite
	)

	channelLogic, err := s.ChannelModel.NewChannelLogicByLink(request.GetHash())
	if err == nil {
		chatInvite = channelLogic.ToChatInvite(md.UserId, func(idList []int32) []*mtproto.User {
			return s.UserModel.GetUsersBySelfAndIDList(md.UserId, idList)
		})
		glog.Infof("messages.checkChatInvite#3eadb1bb - reply: {%s}", logger.JsonDebugData(chatInvite))
		return chatInvite, nil
	}

	// TODO(@benqi): do chat checkChatInvite
	glog.Errorf("messages.checkChatInvite#3eadb1bb - error: {%v}", err)
	return nil, err
}
