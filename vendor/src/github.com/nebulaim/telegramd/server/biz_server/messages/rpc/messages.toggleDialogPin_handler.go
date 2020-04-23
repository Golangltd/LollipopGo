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
	"golang.org/x/net/context"
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"time"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

// messages.toggleDialogPin#a731e257 flags:# pinned:flags.0?true peer:InputDialogPeer = Bool;
func (s *MessagesServiceImpl) MessagesToggleDialogPin(ctx context.Context, request *mtproto.TLMessagesToggleDialogPin) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.toggleDialogPin#3289be6a - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	peer := base.FromInputPeer2(md.UserId, request.GetPeer().GetData2().GetPeer())

	if peer.PeerType == base.PEER_EMPTY {
		glog.Error("empty peer")
		return mtproto.ToBool(false), nil
	}

	// TODO(@benqi): check access_hash
	dialogLogic := s.DialogModel.MakeDialogLogic(md.UserId, peer.PeerType, peer.PeerId)
	dialogLogic.ToggleDialogPin(request.GetPinned())

	// sync other sessions
	updateDialogPinned := &mtproto.TLUpdateDialogPinned{Data2: &mtproto.Update_Data{
		Pinned:  request.GetPinned(),
		Peer_39: peer.ToPeer(),
	}}
	updates := &mtproto.TLUpdates{Data2: &mtproto.Updates_Data{
		Updates: []*mtproto.Update{updateDialogPinned.To_Update()},
		Users:   []*mtproto.User{},
		Chats:   []*mtproto.Chat{},
		Seq:     0,
		Date:    int32(time.Now().Unix()),
	}}

	sync_client.GetSyncClient().SyncUpdatesNotMe(md.UserId, md.AuthId, updates.To_Updates())

	glog.Info("messages.toggleDialogPin#a731e257 - reply {true}")
	return mtproto.ToBool(true), nil
}
