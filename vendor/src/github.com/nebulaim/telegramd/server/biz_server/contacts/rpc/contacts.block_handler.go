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
	// updates2 "github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/proto/mtproto"
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"golang.org/x/net/context"
)

// contacts.block#332b49fc id:InputUser = Bool;
func (s *ContactsServiceImpl) ContactsBlock(ctx context.Context, request *mtproto.TLContactsBlock) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("contacts.block#332b49fc - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
		blockId int32
		id      = request.Id
	)

	switch id.GetConstructor() {
	case mtproto.TLConstructor_CRC32_inputUserSelf:
		blockId = md.UserId
	case mtproto.TLConstructor_CRC32_inputUser:
		// Check access hash
		if ok := s.UserModel.CheckAccessHashByUserId(id.GetData2().GetUserId(), id.GetData2().GetAccessHash()); !ok {
			// TODO(@benqi): Add ACCESS_HASH_INVALID codes
			err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
			glog.Error(err, ": is access_hash error")
			return nil, err
		}

		blockId = id.GetData2().GetUserId()
		// TODO(@benqi): contact exist
	default:
		// mtproto.TLConstructor_CRC32_inputUserEmpty:
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error(err, ": is inputUserEmpty")
		return nil, err
	}

	contactLogic := s.ContactModel.MakeContactLogic(md.UserId)
	blocked := contactLogic.BlockUser(blockId)

	if blocked {
		// Sync unblocked: updateUserBlocked
		//updateUserBlocked := &mtproto.TLUpdateUserBlocked{Data2: &mtproto.Update_Data{
		//	UserId:  blockId,
		//	Blocked: mtproto.ToBool(true),
		//}}
		//
		//blockedUpdates := updates2.NewUpdatesLogic(md.UserId)
		//blockedUpdates.AddUpdate(updateUserBlocked.To_Update())
		//blockedUpdates.AddUser(s.UserModel.GetUserById(md.UserId, blockId).To_User())

		// TODO(@benqi): handle seq
		// sync_client.GetSyncClient().SyncUpdatesData(md.AuthId, md.SessionId, blockId, blockedUpdates.ToUpdates())
	}

	// Blocked会影响收件箱
	glog.Infof("contacts.block#332b49fc - reply: {%v}", blocked)
	return mtproto.ToBool(blocked), nil
}
