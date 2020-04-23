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

package user_client

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/mysql_client"
	"github.com/nebulaim/telegramd/biz/dal/dao/mysql_dao"
	"github.com/nebulaim/telegramd/proto/mtproto"
)

type localUserFacade struct {
	// TODO(@benqi): add user cache
	*mysql_dao.UsersDAO
}

func localUserFacadeInstance() UserFacade {
	return &localUserFacade{}
}

func NewLocalUserApi(dbName string) (*localUserFacade, error) {
	var err error

	dbClient := mysql_client.GetMysqlClient(dbName)
	if dbClient == nil {
		err = fmt.Errorf("invalid dbName: %s", dbName)
		glog.Error(err)
		return nil, err
	}

	return &localUserFacade{UsersDAO: mysql_dao.NewUsersDAO(dbClient)}, nil
}

func (c *localUserFacade) Initialize(config string) error {
	glog.Info("localUserApi - Initialize config: ", config)

	var err error

	dbName := config
	dbClient := mysql_client.GetMysqlClient(dbName)
	if dbClient == nil {
		err = fmt.Errorf("invalid dbName: %s", dbName)
		glog.Error(err)
	}
	c.UsersDAO = mysql_dao.NewUsersDAO(dbClient)

	return err
}

func (c *localUserFacade) GetUser(id int32) (*mtproto.User, error) {
	return nil, nil
}

func (c *localUserFacade) GetUserList(idList []int32) ([]*mtproto.User, error) {
	return nil, nil
}

func (c *localUserFacade) GetUserByPhoneNumber(phone string) (*mtproto.User, error) {
	return nil, nil
}

func init() {
	Register("local", localUserFacadeInstance)
}
