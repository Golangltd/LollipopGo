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

// messages.reportEncryptedSpam#4b0c8c0f peer:InputEncryptedChat = Bool;
func (s *MessagesServiceImpl) MessagesReportEncryptedSpam(ctx context.Context, request *mtproto.TLMessagesReportEncryptedSpam) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("MessagesReportEncryptedSpam - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// peer := helper.FromInputPeer(request.GetPeer())
	//
	// if peer.PeerType == helper.PEER_USER || peer.PeerType == helper.PEER_CHAT {
	//	// TODO(@benqi): 入库
	// }

	glog.Info("MessagesReportEncryptedSpam - reply: {true}")
	return mtproto.ToBool(true), nil
}
