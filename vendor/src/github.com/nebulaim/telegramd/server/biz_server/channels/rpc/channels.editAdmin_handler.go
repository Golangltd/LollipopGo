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
	"github.com/nebulaim/telegramd/biz/core/channel"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

// channels.editAdmin#20b88214 channel:InputChannel user_id:InputUser admin_rights:ChannelAdminRights = Updates;
func (s *ChannelsServiceImpl) ChannelsEditAdmin(ctx context.Context, request *mtproto.TLChannelsEditAdmin) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("ChannelsEditAdmin - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): check channel_id and access_hash

	channelId := request.GetChannel().GetData2().GetChannelId()
	editAdminUserId := request.GetUserId().GetData2().GetUserId()

	channelLogic, _ := s.ChannelModel.NewChannelLogicById(channelId)
	err := channelLogic.EditAdminRights(md.UserId, request.GetUserId().GetData2().GetUserId(), request.GetAdminRights())
	if err != nil {
		glog.Errorf("channels.editTitle#566decd0 - error: %s", err)
		return nil, err
	}

	// reply
	replyUpdates := updates.NewUpdatesLogic(md.UserId)
	replyUpdates.AddChat(channelLogic.ToChannel(md.UserId))

	// push
	pushUpdates := updates.NewUpdatesLogic(editAdminUserId)
	pushUpdates.AddUpdate(channel.MakeUpdateChannel(channelId))
	pushUpdates.AddChat(channelLogic.ToChannel(editAdminUserId))
	sync_client.GetSyncClient().PushChannelUpdates(channelId, editAdminUserId, pushUpdates.ToUpdates())

	reply := replyUpdates.ToUpdates()
	glog.Infof("channels.editTitle#566decd0 - reply: {%s}", logger.JsonDebugData(reply))
	return reply, nil
}
