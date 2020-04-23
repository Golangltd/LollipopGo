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

// channels.readHistory#cc104937 channel:InputChannel max_id:int = Bool;
func (s *ChannelsServiceImpl) ChannelsReadHistory(ctx context.Context, request *mtproto.TLChannelsReadHistory) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.readHistory#cc104937 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	if request.GetChannel().GetConstructor() == mtproto.TLConstructor_CRC32_inputChannelEmpty {
		glog.Infof("channels.readHistory#cc104937 - reply: {false}")
		return mtproto.ToBool(false), nil
	}

	// TODO(@benqi): check access_hash
	channelId := request.GetChannel().GetData2().GetChannelId()

	channelLogic, _ := s.ChannelModel.NewChannelLogicById(channelId)
	channelLogic.ReadOutboxHistory(md.UserId, request.GetMaxId())

	syncUpdates := updates.NewUpdatesLogic(md.UserId)
	updateReadChannelInbox := &mtproto.TLUpdateReadChannelInbox{Data2: &mtproto.Update_Data{
		ChannelId: channelId,
		MaxId:     request.GetMaxId(),
	}}
	syncUpdates.AddUpdate(updateReadChannelInbox.To_Update())
	sync_client.GetSyncClient().SyncChannelUpdatesNotMe(channelId, md.UserId, md.AuthId, syncUpdates.ToUpdates())

	glog.Infof("channels.readHistory#cc104937 - reply: {true}")
	return mtproto.ToBool(true), nil
}
