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

package user

import (
	"github.com/nebulaim/telegramd/proto/mtproto"
)

func (m *UserModel) GetContactUserIDList(userId int32) []int32 {
	contactsDOList := m.dao.UserContactsDAO.SelectUserContacts(userId)
	idList := make([]int32, 0, len(contactsDOList))

	for _, do := range contactsDOList {
		idList = append(idList, do.ContactUserId)
	}
	return idList
}

func (m *UserModel) GetStatusList(selfId int32) []*mtproto.ContactStatus {
	//doList := m.dao.UserContactsDAO.SelectUserContacts(selfId)
	//
	//contactIdList := make([]int32, 0, len(doList))
	//for _, do := range doList {
	//	contactIdList = append(contactIdList, do.ContactUserId)
	//}
	//return nil

	// TODO(@benqi): impl
	return []*mtproto.ContactStatus{}
}
