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

type DevicesDAO struct {
	db *sqlx.DB
}

func NewDevicesDAO(db *sqlx.DB) *DevicesDAO {
	return &DevicesDAO{db}
}

// insert into devices(auth_key_id, user_id, token_type, token) values (:auth_key_id, :user_id, :token_type, :token)
// TODO(@benqi): sqlmap
func (dao *DevicesDAO) Insert(do *dataobject.DevicesDO) int64 {
	var query = "insert into devices(auth_key_id, user_id, token_type, token) values (:auth_key_id, :user_id, :token_type, :token)"
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

// select id, auth_key_id, user_id, token_type, token from devices where token_type = :token_type and token = :token
// TODO(@benqi): sqlmap
func (dao *DevicesDAO) SelectByToken(token_type int8, token string) *dataobject.DevicesDO {
	var query = "select id, auth_key_id, user_id, token_type, token from devices where token_type = ? and token = ?"
	rows, err := dao.db.Queryx(query, token_type, token)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByToken(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.DevicesDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByToken(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByToken(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// select id, auth_key_id, user_id, token_type, token from devices where token_type = :token_type and token = :token and state = 1
// TODO(@benqi): sqlmap
func (dao *DevicesDAO) SelectListById(token_type int8, token string) []dataobject.DevicesDO {
	var query = "select id, auth_key_id, user_id, token_type, token from devices where token_type = ? and token = ? and state = 1"
	rows, err := dao.db.Queryx(query, token_type, token)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectListById(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.DevicesDO
	for rows.Next() {
		v := dataobject.DevicesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectListById(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectListById(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// update devices set state = :state where id = :id
// TODO(@benqi): sqlmap
func (dao *DevicesDAO) UpdateStateById(state int8, id int64) int64 {
	var query = "update devices set state = ? where id = ?"
	r, err := dao.db.Exec(query, state, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateStateById(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateStateById(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update devices set state = :state where token_type = :token_type and token = :token
// TODO(@benqi): sqlmap
func (dao *DevicesDAO) UpdateStateByToken(state int8, token_type int8, token string) int64 {
	var query = "update devices set state = ? where token_type = ? and token = ?"
	r, err := dao.db.Exec(query, state, token_type, token)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateStateByToken(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateStateByToken(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}
