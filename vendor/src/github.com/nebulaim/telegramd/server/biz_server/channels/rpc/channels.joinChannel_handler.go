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
	"github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
	"github.com/nebulaim/telegramd/biz/base"
)

// channels.joinChannel#24b524c5 channel:InputChannel = Updates;
func (s *ChannelsServiceImpl) ChannelsJoinChannel(ctx context.Context, request *mtproto.TLChannelsJoinChannel) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.joinChannel#24b524c5 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	if request.Channel.Constructor == mtproto.TLConstructor_CRC32_inputChannelEmpty {
		// TODO(@benqi): chatUser不能是inputUser和inputUserSelf
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error("channels.exportInvite#c7560885 - error: ", err, "; InputPeer invalid")
		return nil, err
	}

	channel := request.GetChannel().To_InputChannel()
	channelLogic, err := s.ChannelModel.NewChannelLogicById(channel.GetChannelId())
	if err != nil {
		// TODO(@benqi): do chat checkChatInvite
		glog.Errorf("messages.importChatInvite#6c50051c - error: {%v}", err)
		return nil, err
	}

	updateChannel := &mtproto.TLUpdateChannel{Data2: &mtproto.Update_Data{
		ChannelId: channelLogic.GetChannelId(),
	}}

	channelLogic.JoinChannel(md.UserId)
	s.DialogModel.InsertOrChannelUpdateDialog(md.UserId, base.PEER_CHANNEL, channelLogic.GetChannelId())

	syncUpdates := updates.NewUpdatesLogic(md.UserId)
	syncUpdates.AddUpdate(updateChannel.To_Update())
	syncUpdates.AddChat(channelLogic.ToChannel(md.UserId))

	sync_client.GetSyncClient().PushChannelUpdates(channelLogic.GetChannelId(), md.UserId, syncUpdates.ToUpdates())

	replyUpdates := updates.NewUpdatesLogic(md.UserId)
	replyUpdates.AddChat(channelLogic.ToChannel(md.UserId))

	reply := replyUpdates.ToUpdates()
	glog.Infof("channels.joinChannel#24b524c5 - reply: %s", logger.JsonDebugData(reply))
	return reply, nil
}
