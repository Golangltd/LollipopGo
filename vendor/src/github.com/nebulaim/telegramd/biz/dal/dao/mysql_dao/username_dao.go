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

type UsernameDAO struct {
	db *sqlx.DB
}

func NewUsernameDAO(db *sqlx.DB) *UsernameDAO {
	return &UsernameDAO{db}
}

// insert into username(peer_type, peer_id, username) values (:peer_type, :peer_id, :username)
// TODO(@benqi): sqlmap
func (dao *UsernameDAO) Insert(do *dataobject.UsernameDO) int64 {
	var query = "insert into username(peer_type, peer_id, username) values (:peer_type, :peer_id, :username)"
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

// select peer_type, peer_id, username from username where peer_type = :peer_type and peer_id = :peer_id
// TODO(@benqi): sqlmap
func (dao *UsernameDAO) SelectByPeer(peer_type int8, peer_id int32) *dataobject.UsernameDO {
	var query = "select peer_type, peer_id, username from username where peer_type = ? and peer_id = ?"
	rows, err := dao.db.Queryx(query, peer_type, peer_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByPeer(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.UsernameDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByPeer(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByPeer(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// select peer_type, peer_id, username from username where peer_type = 2 and peer_id = :peer_id
// TODO(@benqi): sqlmap
func (dao *UsernameDAO) SelectByUserId(peer_id int32) *dataobject.UsernameDO {
	var query = "select peer_type, peer_id, username from username where peer_type = 2 and peer_id = ?"
	rows, err := dao.db.Queryx(query, peer_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByUserId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.UsernameDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByUserId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByUserId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// select peer_type, peer_id, username from username where peer_type = 4 and peer_id = :peer_id
// TODO(@benqi): sqlmap
func (dao *UsernameDAO) SelectByChannelId(peer_id int32) *dataobject.UsernameDO {
	var query = "select peer_type, peer_id, username from username where peer_type = 4 and peer_id = ?"
	rows, err := dao.db.Queryx(query, peer_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByChannelId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.UsernameDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByChannelId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByChannelId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// select peer_type, peer_id, username from username where username = :username
// TODO(@benqi): sqlmap
func (dao *UsernameDAO) SelectByUsername(username string) *dataobject.UsernameDO {
	var query = "select peer_type, peer_id, username from username where username = ?"
	rows, err := dao.db.Queryx(query, username)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByUsername(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.UsernameDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByUsername(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByUsername(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// update username set username = :username where peer_type = :peer_type and peer_id = :peer_id
// TODO(@benqi): sqlmap
func (dao *UsernameDAO) UpdateUsername(username string, peer_type int8, peer_id int32) int64 {
	var query = "update username set username = ? where peer_type = ? and peer_id = ?"
	r, err := dao.db.Exec(query, username, peer_type, peer_id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateUsername(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateUsername(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}
