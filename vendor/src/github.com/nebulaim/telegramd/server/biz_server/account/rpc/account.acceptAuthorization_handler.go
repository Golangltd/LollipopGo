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
    "github.com/nebulaim/telegramd/proto/mtproto"
    "golang.org/x/net/context"
    "fmt"
    "github.com/nebulaim/telegramd/baselib/grpc_util"
    "github.com/nebulaim/telegramd/baselib/logger"
)

// account.acceptAuthorization#e7027c94 bot_id:int scope:string public_key:string value_hashes:Vector<SecureValueHash> credentials:SecureCredentialsEncrypted = Bool;
func (s *AccountServiceImpl) AccountAcceptAuthorization(ctx context.Context, request *mtproto.TLAccountAcceptAuthorization) (*mtproto.Bool, error) {
    md := grpc_util.RpcMetadataFromIncoming(ctx)
    glog.Infof("AccountAcceptAuthorization - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

    // TODO(@benqi): Impl AccountAcceptAuthorization logic

    return nil, fmt.Errorf("Not impl AccountAcceptAuthorization")
}
