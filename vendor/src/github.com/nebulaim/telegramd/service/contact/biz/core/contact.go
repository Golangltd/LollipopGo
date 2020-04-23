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

package core

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/mysql_client"
	"github.com/nebulaim/telegramd/service/contact/biz/dal/dao/mysql_dao"
	"github.com/nebulaim/telegramd/service/contact/proto"
)

// var modelInstance *contactModel

type contactsDAO struct {
	*mysql_dao.UnregisteredContactsDAO
	*mysql_dao.UserContactsDAO
	*mysql_dao.PopularContactsDAO
}

type ContactModel struct {
	dao *contactsDAO
}

func InitContactModel(dbName string) (*ContactModel, error) {
	// mysql_dao.UnregisteredContactsDAO{}
	dbClient := mysql_client.GetMysqlClient(dbName)
	if dbClient == nil {
		err := fmt.Errorf("invalid dbName: %s", dbName)
		glog.Error(err)
		return nil, err
	}

	m := &ContactModel{dao: &contactsDAO{
		UnregisteredContactsDAO: mysql_dao.NewUnregisteredContactsDAO(dbClient),
		UserContactsDAO:         mysql_dao.NewUserContactsDAO(dbClient),
		PopularContactsDAO:      mysql_dao.NewPopularContactsDAO(dbClient),
	}}
	return m, nil
}

func (m *ContactModel) makeContactLogic(userId int32) *contactLogic {
	return &contactLogic{selfUserID: userId, dao: m.dao}
}

func (m *ContactModel) ImportContacts(selfUserId int32, contacts []*contact.InputContactData) []*contact.ImportedContactData {
	if len(contacts) == 0 {
		glog.Errorf("phoneContacts not empty.")
		return []*contact.ImportedContactData{}
		// return
	}

	logic := m.makeContactLogic(selfUserId)
	if len(contacts) == 1 {
		return []*contact.ImportedContactData{logic.importContact(contacts[0])}
	} else {
		// sync phone book
		return logic.importContacts(contacts)
	}
}

func (m *ContactModel) DeleteContact(selfUserId, contactUserId int32) *contact.DeleteResult {
	logic := m.makeContactLogic(selfUserId)
	return logic.deleteContact(contactUserId)
}

func (m *ContactModel) DeleteContacts(selfUserId int32, contactUserIdList []int32) []*contact.DeleteResult {
	logic := m.makeContactLogic(selfUserId)
	return logic.deleteContacts(contactUserIdList)
}

func (m *ContactModel) BlockUser(selfUserId, id int32) bool {
	logic := m.makeContactLogic(selfUserId)
	return logic.blockUser(id)
}

func (m *ContactModel) UnBlockUser(selfUserId, id int32) bool {
	logic := m.makeContactLogic(selfUserId)
	return logic.unBlockUser(id)
}

func (m *ContactModel) CheckContactAndMutual(selfUserId, id int32) (bool, bool) {
	logic := m.makeContactLogic(selfUserId)
	return logic.checkContactAndMutual(id)
}
