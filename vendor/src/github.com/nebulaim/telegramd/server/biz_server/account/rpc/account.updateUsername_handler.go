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
	"github.com/nebulaim/telegramd/server/sync/sync_client"
	"github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/biz/core/username"
	base2 "github.com/nebulaim/telegramd/biz/base"
)

// account.updateUsername#3e0bdd7c username:string = User;
func (s *AccountServiceImpl) AccountUpdateUsername(ctx context.Context, request *mtproto.TLAccountUpdateUsername) (*mtproto.User, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("account.updateUsername#3e0bdd7c - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))


	username2 := request.GetUsername()
	if username2 != "" {
		if len(request.Username) < username.MIN_USERNAME_LEN || !base.IsAlNumString(request.Username) || base.IsNumber(request.Username[0]) {
			err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_USERNAME_INVALID)
			glog.Error("account.updateUsername#3e0bdd7c - format error: ", err)
			return nil, err
		}

		existed := s.UsernameModel.CheckAccountUsername(md.UserId, request.Username)
		if existed == username.USERNAME_EXISTED_NOTME {
			err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_USERNAME_OCCUPIED)
			glog.Error("account.updateUsername#3e0bdd7c - format error: ", err)
			return nil, err
		}
	}

	// affected
	s.UsernameModel.UpdateUsernameByPeer(base2.PEER_USER, md.GetUserId(), request.GetUsername())

	user := s.UserModel.GetUserById(md.UserId, md.UserId)
	// 要考虑到数据库主从同步问题
	user.SetUsername(request.GetUsername())

	// sync to other sessions
	// updateUserName#a7332b73 user_id:int first_name:string last_name:string username:string = Update;
	updateUserName := &mtproto.TLUpdateUserName{Data2: &mtproto.Update_Data{
		UserId:    md.UserId,
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Username:  request.GetUsername(),
	}}

	syncUpdates := updates.NewUpdatesLogicByUpdate(md.UserId, updateUserName.To_Update())
	sync_client.GetSyncClient().SyncUpdatesNotMe(md.UserId, md.AuthId, syncUpdates.ToUpdateShort())

	//pushUpdates := syncUpdates
	//sync_client.GetSyncClient().PushUpdates(md.UserId, updateUserName.To_Update())

	// TODO(@benqi): push to other contacts

	glog.Infof("account.updateUsername#3e0bdd7c - reply: %s", logger.JsonDebugData(user))
	return user.To_User(), nil
}
