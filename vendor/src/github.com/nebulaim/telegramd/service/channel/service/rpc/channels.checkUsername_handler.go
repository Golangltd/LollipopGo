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

// channels.checkUsername#10e6bd2c channel:InputChannel username:string = Bool;
func (s *ChannelsServiceImpl) ChannelsCheckUsername(ctx context.Context, request *mtproto.TLChannelsCheckUsername) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.checkUsername#10e6bd2c - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var checked bool
	if request.GetChannel().GetConstructor() == mtproto.TLConstructor_CRC32_inputChannelEmpty {
		checked = false
	} else {
		checked = s.ChannelModel.CheckChannelUserName(request.GetChannel().GetData2().GetChannelId(), request.GetUsername())
	}

	glog.Infof("channels.checkUsername#10e6bd2c - reply: {%v}", checked)
	return mtproto.ToBool(checked), nil
}
