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

package account

import (
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
	"github.com/golang/glog"
)

// not found, return 0
func (m *AccountModel) CheckUsername(username string) int32 {
	do := m.dao.UsernameDAO.SelectByUsername(username)
	if do == nil {
		return 0
	}
	return do.PeerId
}

func (m *AccountModel) UpdateUsernameByUserId(id int32, username string) bool {
	usernameDO := m.dao.UsernameDAO.SelectByUsername(username)
	if usernameDO == nil {
		usernameDO = &dataobject.UsernameDO{
			PeerType: 4,
			PeerId:   id,
			Username: username,
		}
		m.dao.UsernameDAO.Insert(usernameDO)
	} else {
		if usernameDO.PeerType == 2 && usernameDO.PeerId == id {
			m.dao.UsernameDAO.UpdateUsername(username, int8(2), id)
		} else {
			glog.Error("username existed.")
			return false
		}
	}
	return true
}

func (m *AccountModel) UpdateFirstAndLastName(id int32, firstName, lastName string) int64 {
	return m.dao.UsersDAO.UpdateFirstAndLastName(firstName, lastName, id)
}

func (m *AccountModel) UpdateAbout(id int32, about string) int64 {
	return m.dao.UsersDAO.UpdateAbout(about, id)
}
