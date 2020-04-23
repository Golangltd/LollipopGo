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

package server

import (
	"context"
	"encoding/base64"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/server/access/auth_key/dal/dao"
	"github.com/nebulaim/telegramd/server/access/auth_key/dal/dao/mysql_dao"
)

type AuthKeyServiceImpl struct {
	*mysql_dao.AuthUsersDAO
	*mysql_dao.AuthKeysDAO
}

// rpc QueryAuthKey(AuthKeyRequest) returns (AuthKeyData);
func (s *AuthKeyServiceImpl) QueryAuthKey(ctx context.Context, request *mtproto.AuthKeyRequest) (*mtproto.AuthKeyData, error) {
	glog.Infof("auth_key.queryAuthKey - request: %s", logger.JsonDebugData(request))

	authKeyData := &mtproto.AuthKeyData{
		Result:    0,
		AuthKeyId: request.AuthKeyId,
	}

	// Check auth_key_id
	if request.AuthKeyId == 0 {
		authKeyData.Result = 1000
	}

	// TODO(@benqi): cache auth_key
	do, err := dao.GetAuthKeysDAO(dao.DB_MASTER).SelectByAuthId(request.AuthKeyId)
	if err != nil {
		glog.Error(err)
		authKeyData.Result = 1001
	} else {
		if do == nil {
			glog.Errorf("read keyData error: not find keyId = %d", request.AuthKeyId)
			authKeyData.Result = 1002
		} else {
			authKeyData.AuthKey, err = base64.RawStdEncoding.DecodeString(do.Body)
			if err != nil {
				glog.Errorf("read keyData error - keyId = %d, %v", request.AuthKeyId, err)
				authKeyData.Result = 1003
			}
		}
	}

	glog.Infof("queryAuthKey {auth_key_id: %d} ok.", request.AuthKeyId)
	return authKeyData, nil
}

// rpc QueryUserId(AuthKeyIdRequest) returns (UserIdResponse);
func (s *AuthKeyServiceImpl) QueryUserId(ctx context.Context, request *mtproto.AuthKeyIdRequest) (*mtproto.UserIdResponse, error) {
	glog.Infof("auth_key.queryUserId - request: %s", logger.JsonDebugData(request))

	userId := &mtproto.UserIdResponse{
		Result:    0,
		AuthKeyId: request.AuthKeyId,
	}

	// Check auth_key_id
	if request.AuthKeyId == 0 {
		userId.Result = 1000
	}

	// TODO(@benqi): cache auth_key
	do, err := dao.GetAuthUsersDAO(dao.DB_MASTER).SelectByAuthId(request.AuthKeyId)
	if err != nil {
		glog.Error(err)
		userId.Result = 1001
	} else {
		if do == nil {
			glog.Errorf("getUserId error: not find keyId = %d", request.AuthKeyId)
			userId.Result = 1002
		} else {
			userId.UserId = do.UserId
		}
	}

	glog.Infof("queryUserId {auth_key_id: %d} ok.", request.AuthKeyId)
	return userId, nil
}
