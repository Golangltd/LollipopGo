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
	"github.com/nebulaim/telegramd/service/auth_session/client"
)

// auth.logOut#5717da40 = Bool;
func (s *AuthServiceImpl) AuthLogOut(ctx context.Context, request *mtproto.TLAuthLogOut) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("auth.logOut#5717da40 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// unbind auth_key and user_id
	auth_session_client.UnbindAuthKeyUser(md.AuthId, md.UserId)

	glog.Info("auth.logOut#5717da40 - reply: {true}")
	return mtproto.ToBool(true), nil
}
