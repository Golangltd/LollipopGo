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
	"github.com/nebulaim/telegramd/biz/base"
)

// channels.editAbout#13e27f1e channel:InputChannel about:string = Bool;
func (s *ChannelsServiceImpl) ChannelsEditAbout(ctx context.Context, request *mtproto.TLChannelsEditAbout) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.editAbout#13e27f1e - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): check channel_id and access_hash

	// channelId := request.GetChannel().GetData2().GetChannelId()
	peer := &base.PeerUtil{
		PeerType: base.PEER_CHANNEL,
		PeerId:   request.GetChannel().GetData2().GetChannelId(),
	}

	channelLogic, _ := s.ChannelModel.NewChannelLogicById(peer.PeerId)
	err := channelLogic.EditAbout(md.UserId, request.GetAbout())

	reply := mtproto.ToBool(err == nil)
	glog.Infof("channels.editAbout#13e27f1e - reply: {%s}", logger.JsonDebugData(reply))
	return reply, nil
}
