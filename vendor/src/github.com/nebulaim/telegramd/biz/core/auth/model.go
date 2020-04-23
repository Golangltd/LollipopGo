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

package auth

import (
	"github.com/nebulaim/telegramd/biz/core"
	"github.com/nebulaim/telegramd/biz/dal/dao"
	"github.com/nebulaim/telegramd/biz/dal/dao/mysql_dao"
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
)

type authsDAO struct {
	*mysql_dao.CommonDAO
	*mysql_dao.AuthUsersDAO
	*mysql_dao.AuthPhoneTransactionsDAO
}

type AuthModel struct {
	dao *authsDAO
}

func (m *AuthModel) InstallModel() {
	m.dao.CommonDAO = dao.GetCommonDAO(dao.DB_MASTER)
	m.dao.AuthUsersDAO = dao.GetAuthUsersDAO(dao.DB_MASTER)
	m.dao.AuthPhoneTransactionsDAO = dao.GetAuthPhoneTransactionsDAO(dao.DB_MASTER)
}

func (m *AuthModel) RegisterCallback(cb interface{}) {
}

func (m *AuthModel) CheckBannedByPhoneNumber(phoneNumber string) bool {
	params := map[string]interface{}{
		"phone": phoneNumber,
	}
	return m.dao.CommonDAO.CheckExists("banned", params)
}

func (m *AuthModel) CheckPhoneNumberExist(phoneNumber string) bool {
	params := map[string]interface{}{
		"phone": phoneNumber,
	}
	return m.dao.CommonDAO.CheckExists("users", params)
}

func (m *AuthModel) BindAuthKeyAndUser(authKeyId int64, userId int32) {
	do3 := m.dao.AuthUsersDAO.SelectByAuthId(authKeyId)
	if do3 == nil {
		do3 := &dataobject.AuthUsersDO{
			AuthId: authKeyId,
			UserId: userId,
		}
		m.dao.AuthUsersDAO.Insert(do3)
	}
}

func (m *AuthModel) MakeCodeData(authKeyId int64, phoneNumber string) *phoneCodeData {
	// TODO(@benqi): 独立出统一消息推送系统
	// 检查phpne是否存在，若存在是否在线决定是否通过短信发送或通过其他客户端发送
	// 透传AuthId，UserId，终端类型等
	// 检查满足条件的TransactionHash是否存在，可能的条件：
	//  1. is_deleted !=0 and now - created_at < 15 分钟
	//

	// sentCodeType, nextCodeType := makeCodeType(phoneRegistered, allowFlashCall, currentNumber)
	code := &phoneCodeData{
		authKeyId:   authKeyId,
		phoneNumber: phoneNumber,
		state:       kCodeStateNone,
		dataType:    kDBTypeCreate,
		dao:         m.dao,
	}
	return code
}

func (m *AuthModel) MakeCancelCodeData(authKeyId int64, phoneNumber, codeHash string) *phoneCodeData {
	code := &phoneCodeData{
		authKeyId:   authKeyId,
		codeHash:    codeHash,
		phoneNumber: phoneNumber,
		state:       kCodeStateNone,
		dataType:    kDBTypeDelete,
		dao:         m.dao,
	}
	return code
}

func (m *AuthModel) MakeCodeDataByHash(authKeyId int64, phoneNumber, codeHash string) *phoneCodeData {
	code := &phoneCodeData{
		authKeyId:   authKeyId,
		codeHash:    codeHash,
		phoneNumber: phoneNumber,
		state:       kCodeStateNone,
		dataType:    kDBTypeLoad,
		dao:         m.dao,
	}
	return code
}

func init() {
	core.RegisterCoreModel(&AuthModel{dao: &authsDAO{}})
}
