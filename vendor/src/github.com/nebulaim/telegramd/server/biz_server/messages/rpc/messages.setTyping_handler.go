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
	"github.com/nebulaim/telegramd/biz/base"
	"github.com/nebulaim/telegramd/proto/mtproto"
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"golang.org/x/net/context"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
	"time"
)

// messages.setTyping#a3825e50 peer:InputPeer action:SendMessageAction = Bool;
func (s *MessagesServiceImpl) MessagesSetTyping(ctx context.Context, request *mtproto.TLMessagesSetTyping) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.setTyping#a3825e50 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	peer := base.FromInputPeer(request.GetPeer())
	if peer.PeerType == base.PEER_SELF || peer.PeerType == base.PEER_USER {
		typing := &mtproto.TLUpdateUserTyping{Data2: &mtproto.Update_Data{
			UserId: md.UserId,
			Action: request.GetAction(),
		}}

		updates := &mtproto.TLUpdates{Data2: &mtproto.Updates_Data{
			Updates: []*mtproto.Update{typing.To_Update()},
			Users:   []*mtproto.User{},
			Chats:   []*mtproto.Chat{},
			Seq:     0,
			Date:    int32(time.Now().Unix()),
		}}
		sync_client.GetSyncClient().PushUpdates(peer.PeerId, updates.To_Updates())
	} else {
		// 其他的不需要推送
	}

	glog.Info("messages.setTyping#a3825e50 - reply: {true}")
	return mtproto.ToBool(true), nil
}
