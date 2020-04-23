/*
 *  Copyright (c) 2018, https://github.com/nebulaim
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

// session.getUserId auth_key_id:long = Int64;
func (s *SessionServiceImpl) SessionGetUserId(ctx context.Context, request *mtproto.TLSessionGetUserId) (*mtproto.Int32, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("session.getUserId - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	userId := s.AuthSessionModel.GetAuthKeyUserId(request.GetAuthKeyId())
	reply := &mtproto.TLInt32{Data2: &mtproto.Int32_Data{
		V: userId,
	}}

	glog.Infof("session.getUserId - reply: {%s}", logger.JsonDebugData(reply))
	return reply.To_Int32(), nil
}
