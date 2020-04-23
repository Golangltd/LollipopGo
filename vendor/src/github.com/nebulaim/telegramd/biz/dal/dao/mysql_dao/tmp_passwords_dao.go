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
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
	"github.com/nebulaim/telegramd/proto/mtproto"
)

type TmpPasswordsDAO struct {
	db *sqlx.DB
}

func NewTmpPasswordsDAO(db *sqlx.DB) *TmpPasswordsDAO {
	return &TmpPasswordsDAO{db}
}

// insert into devices(auth_id, user_id, password_hash, period, tmp_password, valid_until) values (:auth_id, :user_id, :password_hash, :period, :tmp_password, :valid_until)
// TODO(@benqi): sqlmap
func (dao *TmpPasswordsDAO) Insert(do *dataobject.TmpPasswordsDO) int64 {
	var query = "insert into devices(auth_id, user_id, password_hash, period, tmp_password, valid_until) values (:auth_id, :user_id, :password_hash, :period, :tmp_password, :valid_until)"
	r, err := dao.db.NamedExec(query, do)
	if err != nil {
		errDesc := fmt.Sprintf("NamedExec in Insert(%v), error: %v", do, err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	id, err := r.LastInsertId()
	if err != nil {
		errDesc := fmt.Sprintf("LastInsertId in Insert(%v)_error: %v", do, err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}
	return id
}
