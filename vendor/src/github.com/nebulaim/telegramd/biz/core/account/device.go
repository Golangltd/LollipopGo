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
)

//const (
//	TOKEN_TYPE_APNS = 1
//	TOKEN_TYPE_GCM = 2
//	TOKEN_TYPE_MPNS = 3
//	TOKEN_TYPE_SIMPLE_PUSH = 4
//	TOKEN_TYPE_UBUNTU_PHONE = 5
//	TOKEN_TYPE_BLACKBERRY = 6
//	// Android里使用
//	TOKEN_TYPE_INTERNAL_PUSH = 7
//)

func (m *AccountModel) RegisterDevice(authKeyId int64, userId int32, tokenType int8, token string) bool {
	do := m.dao.DevicesDAO.SelectByToken(tokenType, token)
	if do == nil {
		do = &dataobject.DevicesDO{
			AuthKeyId: authKeyId,
			UserId:    userId,
			TokenType: tokenType,
			Token:     token,
		}
		do.Id = m.dao.DevicesDAO.Insert(do)
	} else {
		m.dao.DevicesDAO.UpdateStateById(0, do.Id)
	}

	return true
}

func (m *AccountModel) UnRegisterDevice(tokenType int8, token string) bool {
	m.dao.DevicesDAO.UpdateStateByToken(int8(1), tokenType, token)
	return true
}
