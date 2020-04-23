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
	"github.com/nebulaim/telegramd/biz/core/channel"
	"github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

// channels.editBanned#bfd915cd channel:InputChannel user_id:InputUser banned_rights:ChannelBannedRights = Updates;
func (s *ChannelsServiceImpl) ChannelsEditBanned(ctx context.Context, request *mtproto.TLChannelsEditBanned) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("ChannelsEditBanned - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): check channel_id and access_hash
	channelId := request.GetChannel().GetData2().ChannelId

	// TODO(@benqi): check user_id and access_hasn
	// userId  := request.GetUserId()
	peerUserId := request.GetUserId().GetData2().UserId

	// bannedRights := request.GetBannedRights()
	// bannedRights := channel.MakeChannelBannedRights(request.GetBannedRights().To_ChannelBannedRights())

	channelLogic, _ := s.ChannelModel.NewChannelLogicById(channelId)
	channelLogic.EditBanned(md.UserId, request.GetUserId().GetData2().UserId, request.GetBannedRights())

	replyUpdates := updates.NewUpdatesLogic(md.UserId)
	replyUpdates.AddChat(channelLogic.ToChannel(md.UserId))

	pushUpdates := updates.NewUpdatesLogicByUpdate(peerUserId, channel.MakeUpdateChannel(channelId))
	pushUpdates.AddChat(channelLogic.ToChannel(peerUserId))

	//if bannedRights.IsForbidden() {
	//	pushUpdates.AddChat(channelLogic.ToChannelForbidden())
	//} else {
	//	pushUpdates.AddChat(channelLogic.ToChannel(peerUserId))
	//}
	//
	sync_client.GetSyncClient().PushChannelUpdates(channelId, peerUserId, pushUpdates.ToUpdates())

	reply := replyUpdates.ToUpdates()
	glog.Infof("channels.editBanned#bfd915cd - reply: {%s}", logger.JsonDebugData(reply))
	return reply, nil
}
