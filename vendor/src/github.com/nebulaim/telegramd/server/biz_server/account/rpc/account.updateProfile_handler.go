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
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"golang.org/x/net/context"
)

// account.updateProfile#78515775 flags:# first_name:flags.0?string last_name:flags.1?string about:flags.2?string = User;
func (s *AccountServiceImpl) AccountUpdateProfile(ctx context.Context, request *mtproto.TLAccountUpdateProfile) (*mtproto.User, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("account.updateProfile#78515775 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): Check first_name and last_name invalid. has err: FIRSTNAME_INVALID, LASTNAME_INVALID

	// Check format
	// about长度<70并且可以为emtpy
	// first_name必须有值
	if len(request.FirstName) > 0 && len(request.About) > 0 {
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_FIRSTNAME_INVALID)
		glog.Error(err)
		return nil, err
	}

	user := s.UserModel.GetUserById(md.UserId, md.UserId)

	if len(request.FirstName) > 0 {
		s.AccountModel.UpdateFirstAndLastName(md.UserId, request.FirstName, request.LastName)

		// return new first_name and last_name.
		user.SetFirstName(request.FirstName)
		user.SetLastName(request.LastName)
	} else {
		s.AccountModel.UpdateAbout(md.UserId, request.About)
	}

	// sync to other sessions
	// updateUserName#a7332b73 user_id:int first_name:string last_name:string username:string = Update;
	updateUserName := &mtproto.TLUpdateUserName{Data2: &mtproto.Update_Data{
		UserId:    md.UserId,
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Username:  user.GetUsername(),
	}}
	_ = updateUserName
	// sync_client.GetSyncClient().PushToUserUpdateShortData(md.UserId, updateUserName.To_Update())
	// TODO(@benqi): push to other contacts

	glog.Infof("account.updateProfile#78515775 - reply: {%v}", user)
	return user.To_User(), nil
}
