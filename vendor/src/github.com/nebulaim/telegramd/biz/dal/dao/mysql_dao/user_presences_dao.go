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

type UserPresencesDAO struct {
	db *sqlx.DB
}

func NewUserPresencesDAO(db *sqlx.DB) *UserPresencesDAO {
	return &UserPresencesDAO{db}
}

// insert into user_presences(user_id, last_seen_at, last_seen_auth_key_id, created_at) values (:user_id, :last_seen_at, :last_seen_auth_key_id, :created_at)
// TODO(@benqi): sqlmap
func (dao *UserPresencesDAO) Insert(do *dataobject.UserPresencesDO) int64 {
	var query = "insert into user_presences(user_id, last_seen_at, last_seen_auth_key_id, created_at) values (:user_id, :last_seen_at, :last_seen_auth_key_id, :created_at)"
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

// update user_presences set last_seen_at = :last_seen_at, last_seen_auth_key_id = :last_seen_auth_key_id, version = version + 1 where user_id = :user_id
// TODO(@benqi): sqlmap
func (dao *UserPresencesDAO) UpdateLastSeen(last_seen_at int64, last_seen_auth_key_id int64, user_id int32) int64 {
	var query = "update user_presences set last_seen_at = ?, last_seen_auth_key_id = ?, version = version + 1 where user_id = ?"
	r, err := dao.db.Exec(query, last_seen_at, last_seen_auth_key_id, user_id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateLastSeen(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateLastSeen(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// select last_seen_at from user_presences where user_id = :user_id
// TODO(@benqi): sqlmap
func (dao *UserPresencesDAO) SelectByUserID(user_id int32) *dataobject.UserPresencesDO {
	var query = "select last_seen_at from user_presences where user_id = ?"
	rows, err := dao.db.Queryx(query, user_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByUserID(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.UserPresencesDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByUserID(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByUserID(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// select user_id, last_seen_at from user_presences where user_id in (:idList) order by user_id asc
// TODO(@benqi): sqlmap
func (dao *UserPresencesDAO) SelectByUserIDList(idList []int32) *dataobject.UserPresencesDO {
	var q = "select user_id, last_seen_at from user_presences where user_id in (?) order by user_id asc"
	query, a, err := sqlx.In(q, idList)
	rows, err := dao.db.Queryx(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByUserIDList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.UserPresencesDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByUserIDList(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByUserIDList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}
