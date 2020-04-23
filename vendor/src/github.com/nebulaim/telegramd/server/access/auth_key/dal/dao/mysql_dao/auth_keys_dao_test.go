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

package mysql_dao

import (
	"fmt"
	"github.com/nebulaim/telegramd/baselib/mysql_client"
	"github.com/nebulaim/telegramd/server/access/auth_key/dal/dataobject"
	"testing"
)

func init() {
	mysqlConfig := mysql_client.MySQLConfig{
		Name:   "immaster",
		DSN:    "root:@/nebulaim?charset=utf8",
		Active: 5,
		Idle:   2,
	}
	mysql_client.InstallMysqlClientManager([]mysql_client.MySQLConfig{mysqlConfig})
	// InstallMysqlDAOManager(mysql_client.GetMysqlClientManager())
}

func TestCheckExists(t *testing.T) {
	authKeysDAO := NewAuthKeysDAO(mysql_client.GetMysqlClient("immaster"))
	do := &dataobject.AuthKeysDO{
		AuthId: 2,
		Body:   "123",
	}

	fmt.Println(authKeysDAO.Insert(do))
}
