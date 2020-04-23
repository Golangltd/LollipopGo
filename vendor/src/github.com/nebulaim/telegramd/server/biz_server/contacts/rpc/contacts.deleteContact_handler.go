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
	updates2 "github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/proto/mtproto"
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"golang.org/x/net/context"
)

// contacts.deleteContact#8e953744 id:InputUser = contacts.Link;
func (s *ContactsServiceImpl) ContactsDeleteContact(ctx context.Context, request *mtproto.TLContactsDeleteContact) (*mtproto.Contacts_Link, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("contacts.deleteContact#8e953744 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
		deleteId int32
		id       = request.Id
	)

	switch id.GetConstructor() {
	case mtproto.TLConstructor_CRC32_inputUserSelf:
		deleteId = md.UserId
	case mtproto.TLConstructor_CRC32_inputUser:
		// Check access hash
		if ok := s.UserModel.CheckAccessHashByUserId(id.GetData2().GetUserId(), id.GetData2().GetAccessHash()); !ok {
			// TODO(@benqi): Add ACCESS_HASH_INVALID codes
			err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
			glog.Error(err, ": is access_hash error")
			return nil, err
		}

		deleteId = id.GetData2().GetUserId()
		// TODO(@benqi): contact exist
	default:
		// mtproto.TLConstructor_CRC32_inputUserEmpty:
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error(err, ": is inputUserEmpty")
		return nil, err
	}

	// selfUser := user2.GetUserById(md.UserId, md.UserId)
	deleteUser := s.UserModel.GetUserById(md.UserId, deleteId)

	contactLogic := s.ContactModel.MakeContactLogic(md.UserId)
	needUpdate := contactLogic.DeleteContact(deleteId, deleteUser.GetMutualContact())

	selfUpdates := updates2.NewUpdatesLogic(md.UserId)
	contactLink := &mtproto.TLUpdateContactLink{Data2: &mtproto.Update_Data{
		UserId:      deleteId,
		MyLink:      mtproto.NewTLContactLinkHasPhone().To_ContactLink(),
		ForeignLink: mtproto.NewTLContactLinkHasPhone().To_ContactLink(),
	}}
	selfUpdates.AddUpdate(contactLink.To_Update())
	selfUpdates.AddUser(deleteUser.To_User())
	// TODO(@benqi): handle seq
	// sync_client.GetSyncClient().PushToUserUpdatesData(md.UserId, selfUpdates.ToUpdates())

	// TODO(@benqi): 推给联系人逻辑需要再考虑考虑
	if needUpdate {
		// TODO(@benqi): push to contact user update contact link
		contactUpdates := updates2.NewUpdatesLogic(deleteUser.GetId())
		contactLink2 := &mtproto.TLUpdateContactLink{Data2: &mtproto.Update_Data{
			UserId:      md.UserId,
			MyLink:      mtproto.NewTLContactLinkContact().To_ContactLink(),
			ForeignLink: mtproto.NewTLContactLinkContact().To_ContactLink(),
		}}
		contactUpdates.AddUpdate(contactLink2.To_Update())

		selfUser := s.UserModel.GetUserById(md.UserId, md.UserId)
		contactUpdates.AddUser(selfUser.To_User())
		// TODO(@benqi): handle seq
		// sync_client.GetSyncClient().PushToUserUpdatesData(deleteId, contactUpdates.ToUpdates())
	}

	////////////////////////////////////////////////////////////////////////////////////////
	contactsLink := &mtproto.TLContactsLink{Data2: &mtproto.Contacts_Link_Data{
		MyLink:      mtproto.NewTLContactLinkHasPhone().To_ContactLink(),
		ForeignLink: mtproto.NewTLContactLinkHasPhone().To_ContactLink(),
		User:        s.UserModel.GetUserById(md.UserId, md.UserId).To_User(),
	}}

	glog.Infof("contacts.deleteContact#8e953744 - reply: %s", logger.JsonDebugData(contactsLink))
	return contactsLink.To_Contacts_Link(), nil
}
