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

package contact

import (
	"github.com/nebulaim/telegramd/biz/core"
	"github.com/nebulaim/telegramd/biz/dal/dao"
	"github.com/nebulaim/telegramd/biz/dal/dao/mysql_dao"
)

type ImportedContactData struct {
	UserId        int32
	Importers     int32
	MutualUpdated bool
}

type InputContactData struct {
	UserId    int32
	Phone     string
	FirstName string
	LastName  string
}

type DeleteResult struct {
	UserId int32
	State  int32
}

type contactsDAO struct {
	*mysql_dao.UserContactsDAO
	*mysql_dao.UsersDAO
	*mysql_dao.UnregisteredContactsDAO
	*mysql_dao.PopularContactsDAO
	*mysql_dao.ImportedContactsDAO
}

type ContactModel struct {
	dao *contactsDAO
}

func (m *ContactModel) InstallModel() {
	m.dao.UserContactsDAO = dao.GetUserContactsDAO(dao.DB_MASTER)
	m.dao.UsersDAO = dao.GetUsersDAO(dao.DB_MASTER)
	m.dao.UnregisteredContactsDAO = dao.GetUnregisteredContactsDAO(dao.DB_MASTER)
	m.dao.PopularContactsDAO = dao.GetPopularContactsDAO(dao.DB_MASTER)
	m.dao.ImportedContactsDAO = dao.GetImportedContactsDAO(dao.DB_MASTER)
}

func (m *ContactModel) RegisterCallback(cb interface{}) {
}

func (m *ContactModel) CheckContactAndMutualByUserId(selfId, contactId int32) (bool, bool) {
	do := m.dao.UserContactsDAO.SelectUserContact(selfId, contactId)
	if do == nil {
		return false, false
	} else {
		return true, do.Mutual == 1
	}
}

/*
func (m *ContactModel) makeContactLogic(userId int32) *contactLogic {
	return &contactLogic{selfUserID: userId, dao: m.dao}
}

func (m *ContactModel) ImportContacts(selfUserId int32, contacts []*InputContactData) []*ImportedContactData {
	if len(contacts) == 0 {
		glog.Errorf("phoneContacts not empty.")
		return []*ImportedContactData{}
		// return
	}

	logic := m.makeContactLogic(selfUserId)
	if len(contacts) == 1 {
		return []*ImportedContactData{logic.importContact(contacts[0])}
	} else {
		// sync phone book
		return logic.importContacts(contacts)
	}
}

func (m *ContactModel) DeleteContact(selfUserId, contactUserId int32) *DeleteResult {
	logic := m.makeContactLogic(selfUserId)
	return logic.deleteContact(contactUserId)
}

func (m *ContactModel) DeleteContacts(selfUserId int32, contactUserIdList []int32) []*DeleteResult {
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
*/

// Impl ContactCallback
func (m *ContactModel) GetContactAndMutual(selfUserId, id int32) (bool, bool) {
	return m.CheckContactAndMutualByUserId(selfUserId, id)
}

func init() {
	core.RegisterCoreModel(&ContactModel{dao: &contactsDAO{}})
}
