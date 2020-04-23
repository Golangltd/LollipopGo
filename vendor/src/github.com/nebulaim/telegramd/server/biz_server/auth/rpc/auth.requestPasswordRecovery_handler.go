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

// Forgot password?

// auth.requestPasswordRecovery#d897bc66 = auth.PasswordRecovery;
func (s *AuthServiceImpl) AuthRequestPasswordRecovery(ctx context.Context, request *mtproto.TLAuthRequestPasswordRecovery) (*mtproto.Auth_PasswordRecovery, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("auth.requestPasswordRecovery#d897bc66 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	passwordLogic, err := s.AccountModel.MakePasswordData(md.UserId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	passwordRecovery, err := passwordLogic.RequestPasswordRecovery()
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	glog.Infof("auth.requestPasswordRecovery#d897bc66 - reply: %s\n", logger.JsonDebugData(passwordRecovery))
	return passwordRecovery, nil
}
