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
)

// channels.leaveChannel#f836aa95 channel:InputChannel = Updates;
func (s *ChannelsServiceImpl) ChannelsLeaveChannel(ctx context.Context, request *mtproto.TLChannelsLeaveChannel) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("ChannelsLeaveChannel - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): check channel_id and access_hash
	channelId := request.GetChannel().GetData2().GetChannelId()

	channelLogic, _ := s.ChannelModel.NewChannelLogicById(channelId)
	channelLogic.LeaveChannel(md.UserId)

	syncUpdates := updates.NewUpdatesLogic(md.UserId)
	syncUpdates.AddChat(channelLogic.ToChannel(md.UserId))

	sync_client.GetSyncClient().SyncChannelUpdatesNotMe(channelId, md.UserId, md.AuthId, syncUpdates.ToUpdates())

	reply := syncUpdates.ToUpdates()
	glog.Infof("channels.editBanned#bfd915cd - reply: {%s}", logger.JsonDebugData(reply))
	return reply, nil
}
