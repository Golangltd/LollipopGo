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
)

// account.getNotifySettings#12b3ad31 peer:InputNotifyPeer = PeerNotifySettings;
func (s *AccountServiceImpl) AccountGetNotifySettings(ctx context.Context, request *mtproto.TLAccountGetNotifySettings) (*mtproto.PeerNotifySettings, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("account.getNotifySettings#12b3ad31 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
		settings *mtproto.PeerNotifySettings
	)

	switch request.GetPeer().GetConstructor() {
	case mtproto.TLConstructor_CRC32_inputNotifyPeer:
		peer := base.FromInputNotifyPeer(request.GetPeer())
		settings = s.AccountModel.GetNotifySettings(md.UserId, peer)
	case mtproto.TLConstructor_CRC32_inputNotifyUsers,
		mtproto.TLConstructor_CRC32_inputNotifyChats:

		peerSettings := &mtproto.TLPeerNotifySettings{Data2: &mtproto.PeerNotifySettings_Data{
			ShowPreviews: mtproto.ToBool(true),
			Silent:       mtproto.ToBool(false),
			MuteUntil:    0,
			Sound:        "default",
		}}
		settings = peerSettings.To_PeerNotifySettings()
	}

	glog.Infof("account.getNotifySettings#12b3ad31 - reply: %s", logger.JsonDebugData(settings))
	return settings, nil
}
