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
	// "github.com/nebulaim/telegramd/biz_server/sync_client"
	// peer2 "github.com/nebulaim/telegramd/biz/core/peer"
)

// account.resetNotifySettings#db7e1747 = Bool;
func (s *AccountServiceImpl) AccountResetNotifySettings(ctx context.Context, request *mtproto.TLAccountResetNotifySettings) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("account.resetNotifySettings#db7e1747 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	s.AccountModel.ResetNotifySettings(md.UserId)

	// TODO(@benqi): update notify setting
	/*
		 Android client source:
			} else if (update instanceof TLRPC.TL_updateNotifySettings) {
				TLRPC.TL_updateNotifySettings updateNotifySettings = (TLRPC.TL_updateNotifySettings) update;
				if (update.notify_settings instanceof TLRPC.TL_peerNotifySettings && updateNotifySettings.peer instanceof TLRPC.TL_notifyPeer) {
		           ......
		        }
		    }
	*/

	//peer := &peer2.PeerUtil{}
	//peer.PeerType = peer2.PEER_ALL
	//update := mtproto.NewTLUpdateNotifySettings()
	//update.SetPeer(peer.ToNotifyPeer())
	//updateSettings := mtproto.NewTLPeerNotifySettings()
	//updateSettings.SetShowPreviews(true)
	//updateSettings.SetSilent(false)
	//updateSettings.SetMuteUntil(0)
	//updateSettings.SetSound("default")
	//update.SetNotifySettings(updateSettings.To_PeerNotifySettings())
	//
	//sync_client.GetSyncClient().PushToUserMeOneUpdateData(md.AuthId, md.SessionId, md.UserId, update.To_Update())

	glog.Infof("account.resetNotifySettings#db7e1747 - reply: {true}")
	return mtproto.ToBool(true), nil
}
