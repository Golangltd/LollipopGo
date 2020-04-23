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

// users.getUsers#d91a548 id:Vector<InputUser> = Vector<User>;
func (s *UsersServiceImpl) UsersGetUsers(ctx context.Context, request *mtproto.TLUsersGetUsers) (*mtproto.Vector_User, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("users.getUsers#d91a548 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): Impl UsersGetUsers logic
	userList := &mtproto.Vector_User{
		Datas: make([]*mtproto.User, 0, len(request.Id)),
	}

	for _, inputUser := range request.Id {
		switch inputUser.GetConstructor() {
		case mtproto.TLConstructor_CRC32_inputUserSelf:
			userData := s.UserModel.GetUserById(md.GetUserId(), md.GetUserId())
			userList.Datas = append(userList.Datas, userData.To_User())
		case mtproto.TLConstructor_CRC32_inputUser:
			userData := s.UserModel.GetUserById(md.GetUserId(), inputUser.GetData2().GetUserId())
			userList.Datas = append(userList.Datas, userData.To_User())
		case mtproto.TLConstructor_CRC32_inputUserEmpty:
		}
	}

	glog.Infof("users.getUsers#d91a548 - reply: ", logger.JsonDebugData(userList))
	return userList, nil
}
