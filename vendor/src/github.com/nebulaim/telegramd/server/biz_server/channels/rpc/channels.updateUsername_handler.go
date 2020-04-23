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
	"github.com/nebulaim/telegramd/biz/core/username"
	"github.com/nebulaim/telegramd/baselib/base"
	base2 "github.com/nebulaim/telegramd/biz/base"
)

// channels.updateUsername#3514b3de channel:InputChannel username:string = Bool;
func (s *ChannelsServiceImpl) ChannelsUpdateUsername(ctx context.Context, request *mtproto.TLChannelsUpdateUsername) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.updateUsername#3514b3de - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	if request.GetUsername() != "" {
		if len(request.Username) < username.MIN_USERNAME_LEN || !base.IsAlNumString(request.Username) || base.IsNumber(request.Username[0]) {
			err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_USERNAME_INVALID)
			glog.Error("account.updateUsername#3e0bdd7c - format error: ", err)
			return nil, err
		}
	}

	channel := request.GetChannel().To_InputChannel()

	// TODO(@benqi): check channel_id and access_hash
	// channelId := request.GetChannel().GetData2().ChannelId
	channelLogic, _ := s.ChannelModel.NewChannelLogicById(channel.GetChannelId())
	err := channelLogic.UpdateUsername(md.UserId, request.GetUsername(), func(channelId int32, username2 string) bool {
		existed := s.UsernameModel.CheckChannelUsername(channelId, username2)
		if existed == username.USERNAME_EXISTED_NOTME {
			err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_USERNAME_OCCUPIED)
			glog.Error("account.updateUsername#3e0bdd7c - format error: ", err)
			return false
		}

		s.UsernameModel.UpdateUsernameByPeer(base2.PEER_CHANNEL, channelId, username2)
		return true
	})

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	glog.Infof("channels.updateUsername#3514b3de - reply: {true}")
	return mtproto.ToBool(true), nil
}
