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

package auth_session

import (
	"github.com/nebulaim/telegramd/service/auth_session/biz/dal/dao/mysql_dao"
	"github.com/nebulaim/telegramd/baselib/mysql_client"
	"github.com/golang/glog"
	"fmt"
	"github.com/nebulaim/telegramd/service/auth_session/biz/dal/dataobject"
	"encoding/base64"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"time"
)

type authSessionDAO struct {
	*mysql_dao.AuthKeysDAO
	*mysql_dao.AuthOpLogsDAO
	*mysql_dao.AuthsDAO
	*mysql_dao.AuthUsersDAO
}

type AuthSessionModel struct {
	dao authSessionDAO
}

func NewAuthSessionModel(dbName, redisName string) *AuthSessionModel {
	db := mysql_client.GetMysqlClient(dbName)
	if db == nil {
		glog.Fatal("not found db: ", dbName)
	}

	m := &AuthSessionModel{dao: authSessionDAO{
		AuthKeysDAO:   mysql_dao.NewAuthKeysDAO(db),
		AuthOpLogsDAO: mysql_dao.NewAuthOpLogsDAO(db),
		AuthsDAO:      mysql_dao.NewAuthsDAO(db),
		AuthUsersDAO:  mysql_dao.NewAuthUsersDAO(db),
	}}
	return m
}

func (m *AuthSessionModel) QueryAuthKey(authKeyId int64) (*mtproto.AuthKeyInfo, error) {
	do := m.dao.AuthKeysDAO.SelectByAuthKeyId(authKeyId)
	if do == nil {
		err := fmt.Errorf("not find key - keyId = %d", authKeyId)
		return nil, err
	}
	authKey, err := base64.RawStdEncoding.DecodeString(do.Body)
	if err != nil {
		glog.Errorf("read keyData error - keyId = %d, %v", authKeyId, err)
		return nil, err
	}
	keyInfo := &mtproto.TLAuthKeyInfo{Data2: &mtproto.AuthKeyInfo_Data{
		AuthKeyId: authKeyId,
		AuthKey:   authKey,
	}}

	// TODO(@benqi): get salt
	return keyInfo.To_AuthKeyInfo(), nil
}

func (m *AuthSessionModel) InsertAuthKey(authKeyId int64, authKey []byte, salt *mtproto.TLFutureSalt) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("storage auth_key error: auth_key_id = %d", authKeyId)
		}
	}()

	do := &dataobject.AuthKeysDO{
		AuthKeyId: authKeyId,
		Body:      base64.RawStdEncoding.EncodeToString(authKey),
	}

	do.Id = int32(m.dao.AuthKeysDAO.Insert(do))

	// TODO(@benqi): cache salt
	return nil
}

func (m *AuthSessionModel) GetApiLayer(authKeyId int64) int32 {
	do := m.dao.AuthsDAO.SelectLayer(authKeyId)
	if do == nil {
		glog.Errorf("not find layer - keyId = %d", authKeyId)
		return 0
	}
	return do.Layer
}

func (m *AuthSessionModel) GetAuthKeyUserId(authKeyId int64) int32 {
	do := m.dao.AuthUsersDAO.Select(authKeyId)
	if do == nil {
		glog.Errorf("not find user - keyId = %d", authKeyId)
		return 0
	}
	return do.UserId
}

func (m *AuthSessionModel) BindAuthKeyUser(authKeyId int64, userId int32) bool {
	now := int32(time.Now().Unix())
	authUsersDO := &dataobject.AuthUsersDO{
		AuthKeyId:   authKeyId,
		UserId:      userId,
		DateCreated: now,
		DateActived: now,
	}
	m.dao.AuthUsersDAO.InsertOrUpdates(authUsersDO)
	return true
}

func (m *AuthSessionModel) UnbindAuthUser(authKeyId int64, userId int32) bool {
	m.dao.AuthUsersDAO.Delete(authKeyId, userId)
	return true
}

func (m *AuthSessionModel) SetClientSessionInfo(session *mtproto.TLClientSessionInfo) bool {
	do := &dataobject.AuthsDO{
		AuthKeyId:      session.GetAuthKeyId(),
		Layer:          session.GetLayer(),
		ApiId:          session.GetApiId(),
		DeviceModel:    session.GetDeviceModel(),
		SystemVersion:  session.GetSystemVersion(),
		AppVersion:     session.GetAppVersion(),
		SystemLangCode: session.GetSystemLangCode(),
		LangPack:       session.GetLangPack(),
		LangCode:       session.GetLangCode(),
		ClientIp:       session.GetIp(),
	}
	m.dao.AuthsDAO.InsertOrUpdate(do)
	return true
}
