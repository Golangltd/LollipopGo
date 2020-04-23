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
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

// channels.inviteToChannel#199f3a6c channel:InputChannel users:Vector<InputUser> = Updates;
func (s *ChannelsServiceImpl) ChannelsInviteToChannel(ctx context.Context, request *mtproto.TLChannelsInviteToChannel) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.inviteToChannel#199f3a6c - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	if request.Channel.Constructor == mtproto.TLConstructor_CRC32_inputChannelEmpty {
		// TODO(@benqi): chatUser不能是inputUser和inputUserSelf
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error("channels.exportInvite#c7560885 - error: ", err, "; InputPeer invalid")
		return nil, err
	}

	channelLogic, err := s.ChannelModel.NewChannelLogicById(request.GetChannel().GetData2().GetChannelId())
	if err != nil {
		glog.Error("channels.inviteToChannel#199f3a6c - error: ", err)
		return nil, err
	}

	updateChannel := &mtproto.TLUpdateChannel{Data2: &mtproto.Update_Data{
		ChannelId: channelLogic.GetChannelId(),
	}}

	for _, u := range request.Users {
		if u.GetConstructor() == mtproto.TLConstructor_CRC32_inputUserEmpty ||
			u.GetConstructor() == mtproto.TLConstructor_CRC32_inputUserSelf {
			// TODO(@benqi): handle inputUserSelf
			continue
		}
		channelLogic.AddChannelUser(md.UserId, u.GetData2().GetUserId())

		psuhUpdates := update2.NewUpdatesLogic(u.GetData2().GetUserId())
		psuhUpdates.AddUpdate(updateChannel.To_Update())
		psuhUpdates.AddChat(channelLogic.ToChannel(u.GetData2().GetUserId()))
		sync_client.GetSyncClient().PushUpdates(u.GetData2().GetUserId(), psuhUpdates.ToUpdates())
	}

	replyUpdates := update2.NewUpdatesLogic(md.UserId)
	replyUpdates.AddUpdate(updateChannel.To_Update())
	replyUpdates.AddChat(channelLogic.ToChannel(md.UserId))

	glog.Infof("channels.inviteToChannel#199f3a6c - reply: {%v}", replyUpdates)
	return replyUpdates.ToUpdates(), nil
}
