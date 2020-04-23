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

// messages.getMessagesViews#c4c8a55d peer:InputPeer id:Vector<int> increment:Bool = Vector<int>;
func (s *MessagesServiceImpl) MessagesGetMessagesViews(ctx context.Context, request *mtproto.TLMessagesGetMessagesViews) (*mtproto.VectorInt, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.getMessagesViews#c4c8a55d - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var viewsList []int32

	if request.GetPeer().GetConstructor() != mtproto.TLConstructor_CRC32_inputPeerChannel {
		viewsList = []int32{}
	} else {
		// TODO(@benqi): push updateChannelMessageViews??
		channelId := request.GetPeer().GetData2().GetChannelId()
		increment :=  mtproto.FromBool(request.GetIncrement())

		viewsList = s.MessageModel.GetChannelMessagesViews(channelId, request.GetId(), increment)
		if increment {
			s.MessageModel.IncrementChannelMessagesViews(channelId, request.GetId())
		}
	}

	views := &mtproto.VectorInt{
		Datas: viewsList,
	}

	glog.Infof("messages.getMessagesViews#c4c8a55d - reply: %s", logger.JsonDebugData(views))
	return views, nil
}
