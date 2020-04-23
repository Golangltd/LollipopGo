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
	// "fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"golang.org/x/net/context"
)

// messages.uploadMedia#519bc2b1 peer:InputPeer media:InputMedia = MessageMedia;
func (s *MessagesServiceImpl) MessagesUploadMedia(ctx context.Context, request *mtproto.TLMessagesUploadMedia) (*mtproto.MessageMedia, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.uploadMedia#519bc2b1 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	messageMedia := s.makeMediaByInputMedia(md.AuthId, request.GetMedia())

	//// TODO(@benqi): Impl MessagesUploadMedia logic
	//return nil, fmt.Errorf("Not impl MessagesUploadMedia")

	glog.Infof("messages.uploadMedia#519bc2b1 - reply: %s", logger.JsonDebugData(messageMedia))
	return messageMedia, nil
}
