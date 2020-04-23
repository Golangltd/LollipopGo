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
	"github.com/nebulaim/telegramd/biz/base"
	updates2 "github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/proto/mtproto"
	// "github.com/nebulaim/telegramd/server/sync/sync_client"
	"golang.org/x/net/context"
)

// Android client有三种场景会调用importContacts:
// 1. 导入通信录
// 2. 导入单个联系人
// contacts.importContacts#2c800be5 contacts:Vector<InputContact> = contacts.ImportedContacts;
func (s *ContactsServiceImpl) ContactsImportContacts(ctx context.Context, request *mtproto.TLContactsImportContacts) (*mtproto.Contacts_ImportedContacts, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("contacts.importContacts#2c800be5 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
		err              error
		importedContacts *mtproto.TLContactsImportedContacts
		region 			 string
		// := s.UserModel.GetCountryCodeByUser(md.UserId)
	)

	if len(request.Contacts) == 0 {
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error(err, ": contacts empty")
		return nil, err
	}

	importList := make([]*mtproto.InputContact_Data, 0, len(request.Contacts))
	importPhoneList := make([]string, 0, len(request.Contacts))
	for _, c2 := range request.GetContacts() {
		c := c2.To_InputPhoneContact()
		// if c.GetPhone() == ""
		// ignore first_name and last_name is emtpy
		if c.GetFirstName() == "" && c.GetLastName() == "" {
			continue
		}
		phone := c.GetPhone()
		pnumber, err := base.MakePhoneNumberUtil(phone, "")
		if err != nil {
			if region == "" {
				region = s.UserModel.GetCountryCodeByUser(md.UserId)
			}
			pnumber, err = base.MakePhoneNumberUtil(phone, region)
			if err != nil {
				continue
			}
		}
		phone = pnumber.GetNormalizeDigits()
		c.SetPhone(phone)
		importPhoneList = append(importPhoneList, phone)
		importList = append(importList, c.GetData2())
	}

	if len(importList) == 0 {
		glog.Warning("")
		importedContacts = &mtproto.TLContactsImportedContacts{ Data2: &mtproto.Contacts_ImportedContacts_Data{
			Imported:       []*mtproto.ImportedContact{},
			PopularInvites: []*mtproto.PopularContact{},
			RetryContacts:  []int64{},
		}}
		return importedContacts.To_Contacts_ImportedContacts(), nil
	}

	// 注意: 目前只支持导入1个并且已经注册的联系人!!!!
	inputContact := request.Contacts[0].To_InputPhoneContact()

	if inputContact.GetFirstName() == "" {
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_FIRSTNAME_INVALID)
		glog.Error(err, ": first_name is empty")
		return nil, err
	}

	pnumber, err := base.MakePhoneNumberUtil(inputContact.GetPhone(), "")
	if err != nil {
		region := s.UserModel.GetCountryCodeByUser(md.UserId)
		pnumber, err = base.MakePhoneNumberUtil(inputContact.GetPhone(), region)
	}

	// var contactUser
	// phone, err := base.CheckAndGetPhoneNumber(inputContact.GetPhone())
	if err != nil {
		glog.Error(err)
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_PHONE_CODE_INVALID)
		glog.Error(err, ": phone code invalid - ", inputContact.GetPhone())
		return nil, err
	}

	phoneNumber := pnumber.GetNormalizeDigits()
	contactUser := s.UserModel.GetUserByPhoneNumber(md.UserId, phoneNumber)
	if contactUser == nil {
		// 该手机号未注册，我们认为手机号出错
		//err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_PHONE_CODE_INVALID)
		//glog.Error(err, ": phone code invalid - ", inputContact.GetPhone())
		//return nil, err

		importedContacts = &mtproto.TLContactsImportedContacts{Data2: &mtproto.Contacts_ImportedContacts_Data{
			Imported:       []*mtproto.ImportedContact{},
			PopularInvites: []*mtproto.PopularContact{},
			RetryContacts:  []int64{},
			Users:          []*mtproto.User{},
		}}

		glog.Infof("contacts.importContacts#2c800be5 - reply: %s", logger.JsonDebugData(importedContacts))
		return importedContacts.To_Contacts_ImportedContacts(), nil
	}
	// contactUser.SetContact(true)
	// contactUser.SetMutualContact(true)
	contactLogic := s.ContactModel.MakeContactLogic(md.UserId)
	needUpdate := contactLogic.ImportContact(contactUser.GetId(), phoneNumber, inputContact.GetFirstName(), inputContact.GetLastName())
	// _ = needUpdate

	selfUpdates := updates2.NewUpdatesLogic(md.UserId)
	contactLink := &mtproto.TLUpdateContactLink{Data2: &mtproto.Update_Data{
		UserId:      contactUser.GetId(),
		MyLink:      mtproto.NewTLContactLinkContact().To_ContactLink(),
		ForeignLink: mtproto.NewTLContactLinkContact().To_ContactLink(),
	}}
	selfUpdates.AddUpdate(contactLink.To_Update())

	// 添加成功后是我的联系人
	contactUser.SetContact(true)
	selfUpdates.AddUser(contactUser.To_User())
	// TODO(@benqi): handle seq
	_ = selfUpdates
	// sync_client.GetSyncClient().PushToUserNotMeUpdatesData(md.AuthId, md.SessionId, md.UserId, selfUpdates.ToUpdates())

	// TODO(@benqi): 推给联系人逻辑需要再考虑考虑
	if needUpdate {
		// TODO(@benqi): push to contact user update contact link
		contactUpdates := updates2.NewUpdatesLogic(contactUser.GetId())
		contactLink2 := &mtproto.TLUpdateContactLink{Data2: &mtproto.Update_Data{
			UserId:      md.UserId,
			MyLink:      mtproto.NewTLContactLinkContact().To_ContactLink(),
			ForeignLink: mtproto.NewTLContactLinkContact().To_ContactLink(),
		}}
		contactUpdates.AddUpdate(contactLink2.To_Update())

		myUser := s.UserModel.GetUserById(contactUser.GetId(), md.UserId)
		contactUpdates.AddUser(myUser.To_User())
		// TODO(@benqi): handle seq
		_ = contactUpdates
		// sync_client.GetSyncClient().PushToUserUpdatesData(contactUser.GetId(), contactUpdates.ToUpdates())
	}

	////////////////////////////////////////////////////////////////////////////////////////
	imported := &mtproto.TLImportedContact{Data2: &mtproto.ImportedContact_Data{
		UserId:   contactUser.GetId(),
		ClientId: inputContact.GetClientId(),
	}}
	// contacts.importedContacts#77d01c3b imported:Vector<ImportedContact> popular_invites:Vector<PopularContact> retry_contacts:Vector<long> users:Vector<User> = contacts.ImportedContacts;
	importedContacts = &mtproto.TLContactsImportedContacts{Data2: &mtproto.Contacts_ImportedContacts_Data{
		Imported: []*mtproto.ImportedContact{imported.To_ImportedContact()},
		Users:    []*mtproto.User{contactUser.To_User()},
	}}

	glog.Infof("contacts.importContacts#2c800be5 - reply: %s", logger.JsonDebugData(importedContacts))
	return importedContacts.To_Contacts_ImportedContacts(), nil
}
