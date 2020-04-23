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
	"github.com/nebulaim/telegramd/biz/base"
	"github.com/nebulaim/telegramd/biz/core/update"
)

// messages.importChatInvite#6c50051c hash:string = Updates;
func (s *MessagesServiceImpl) MessagesImportChatInvite(ctx context.Context, request *mtproto.TLMessagesImportChatInvite) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.importChatInvite#6c50051c - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	channelLogic, err := s.ChannelModel.NewChannelLogicByLink(request.GetHash())
	if err != nil {
		// TODO(@benqi): do chat checkChatInvite
		glog.Errorf("messages.importChatInvite#6c50051c - error: {%v}", err)
		return nil, err
	}

	updateChannel := &mtproto.TLUpdateChannel{Data2: &mtproto.Update_Data{
		ChannelId: channelLogic.GetChannelId(),
	}}

	// TODO(@benqi): importChatInvite
	channelLogic.InviteToChannel(channelLogic.CreatorUserId, md.UserId)
	s.DialogModel.InsertOrChannelUpdateDialog(md.UserId, base.PEER_CHANNEL, channelLogic.GetChannelId())

	syncUpdates := updates.NewUpdatesLogic(md.UserId)
	syncUpdates.AddUpdate(updateChannel.To_Update())
	syncUpdates.AddChat(channelLogic.ToChannel(md.UserId))

	sync_client.GetSyncClient().PushChannelUpdates(channelLogic.GetChannelId(), md.UserId, syncUpdates.ToUpdates())

	replyUpdates := updates.NewUpdatesLogic(md.UserId)
	replyUpdates.AddChat(channelLogic.ToChannel(md.UserId))

	reply := replyUpdates.ToUpdates()
	glog.Infof("channels.inviteToChannel#199f3a6c - reply: %s", logger.JsonDebugData(reply))
	return reply, nil

}
