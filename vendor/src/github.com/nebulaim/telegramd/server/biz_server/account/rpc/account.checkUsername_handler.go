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
	"github.com/nebulaim/telegramd/baselib/base"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"golang.org/x/net/context"
	"github.com/nebulaim/telegramd/biz/core/username"
)

// account.checkUsername#2714d86c username:string = Bool;
func (s *AccountServiceImpl) AccountCheckUsername(ctx context.Context, request *mtproto.TLAccountCheckUsername) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("account.checkUsername#2714d86c - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// Check username format
	// You can choose a username on Telegram.
	// If you do, other people will be able to find
	// you by this username and contact you
	// without knowing your phone number.
	//
	// You can use a-z, 0-9 and underscores.
	// Minimum length is 5 characters.";
	//
	if len(request.Username) < username.MIN_USERNAME_LEN || !base.IsAlNumString(request.Username) || base.IsNumber(request.Username[0]) {
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_USERNAME_INVALID)
		glog.Error("account.checkUsername#2714d86c - format error: ", err)
		return nil, err
	} else {
		existed := s.UsernameModel.CheckAccountUsername(md.UserId, request.GetUsername())
		if existed == username.USERNAME_EXISTED_NOTME {
			err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_USERNAME_OCCUPIED)
			glog.Error("account.checkUsername#2714d86c - exists username: ", err)
			return nil, err
		}
	}

	glog.Infof("account.checkUsername#2714d86c - reply: {true}")
	return mtproto.ToBool(true), nil
}
