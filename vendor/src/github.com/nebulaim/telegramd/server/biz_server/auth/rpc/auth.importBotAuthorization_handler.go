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
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"golang.org/x/net/context"
)

// auth.importBotAuthorization#67a3ff2c flags:int api_id:int api_hash:string bot_auth_token:string = auth.Authorization;
func (s *AuthServiceImpl) AuthImportBotAuthorization(ctx context.Context, request *mtproto.TLAuthImportBotAuthorization) (*mtproto.Auth_Authorization, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("AuthImportBotAuthorization - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): Impl AuthImportBotAuthorization logic

	return nil, fmt.Errorf("Not impl AuthImportBotAuthorization")
}
