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
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/service/contact/biz/dal/dataobject"
)

type UserContactsDAO struct {
	db *sqlx.DB
}

func NewUserContactsDAO(db *sqlx.DB) *UserContactsDAO {
	return &UserContactsDAO{db}
}

// insert into user_contacts(owner_user_id, contact_user_id, contact_phone, contact_first_name, contact_last_name, mutual, date2) values (:owner_user_id, :contact_user_id, :contact_phone, :contact_first_name, :contact_last_name, :mutual, :date2)
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) Insert(do *dataobject.UserContactsDO) int64 {
	var query = "insert into user_contacts(owner_user_id, contact_user_id, contact_phone, contact_first_name, contact_last_name, mutual, date2) values (:owner_user_id, :contact_user_id, :contact_phone, :contact_first_name, :contact_last_name, :mutual, :date2)"
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

// select id, owner_user_id, contact_user_id, contact_phone, contact_first_name, contact_last_name, mutual, is_deleted from user_contacts where owner_user_id = :owner_user_id and contact_user_id = :contact_user_id
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) SelectUserContact(owner_user_id int32, contact_user_id int32) *dataobject.UserContactsDO {
	var query = "select id, owner_user_id, contact_user_id, contact_phone, contact_first_name, contact_last_name, mutual, is_deleted from user_contacts where owner_user_id = ? and contact_user_id = ?"
	rows, err := dao.db.Queryx(query, owner_user_id, contact_user_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectUserContact(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.UserContactsDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectUserContact(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectUserContact(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// select id, owner_user_id, contact_user_id, contact_phone, contact_first_name, contact_last_name, mutual, is_deleted from user_contacts where owner_user_id = :owner_user_id
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) SelectAllUserContacts(owner_user_id int32) []dataobject.UserContactsDO {
	var query = "select id, owner_user_id, contact_user_id, contact_phone, contact_first_name, contact_last_name, mutual, is_deleted from user_contacts where owner_user_id = ?"
	rows, err := dao.db.Queryx(query, owner_user_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectAllUserContacts(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.UserContactsDO
	for rows.Next() {
		v := dataobject.UserContactsDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectAllUserContacts(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectAllUserContacts(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select id, owner_user_id, contact_user_id, contact_phone, contact_first_name, contact_last_name, mutual, is_deleted from user_contacts where owner_user_id = :owner_user_id and is_deleted = 0 order by contact_user_id asc
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) SelectUserContacts(owner_user_id int32) []dataobject.UserContactsDO {
	var query = "select id, owner_user_id, contact_user_id, contact_phone, contact_first_name, contact_last_name, mutual, is_deleted from user_contacts where owner_user_id = ? and is_deleted = 0 order by contact_user_id asc"
	rows, err := dao.db.Queryx(query, owner_user_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectUserContacts(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.UserContactsDO
	for rows.Next() {
		v := dataobject.UserContactsDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectUserContacts(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectUserContacts(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select contact_user_id from user_contacts where owner_user_id = :owner_user_id and is_blocked = 1 and is_deleted = 0 order by id asc limit :limit
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) SelectBlockedList(owner_user_id int32, limit int32) []dataobject.UserContactsDO {
	var query = "select contact_user_id from user_contacts where owner_user_id = ? and is_blocked = 1 and is_deleted = 0 order by id asc limit ?"
	rows, err := dao.db.Queryx(query, owner_user_id, limit)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectBlockedList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.UserContactsDO
	for rows.Next() {
		v := dataobject.UserContactsDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectBlockedList(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectBlockedList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// update user_contacts set contact_first_name = :contact_first_name, contact_last_name = :contact_last_name, is_deleted = 0 where id = :id
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) UpdateContactNameByID(contact_first_name string, contact_last_name string, id int64) int64 {
	var query = "update user_contacts set contact_first_name = ?, contact_last_name = ?, is_deleted = 0 where id = ?"
	r, err := dao.db.Exec(query, contact_first_name, contact_last_name, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateContactNameByID(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateContactNameByID(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update user_contacts set contact_user_id = :contact_user_id, mutual = :mutual where contact_user_id = 0 and owner_user_id = :owner_user_id and contact_phone = :contact_phone
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) UpdateContactUserId(contact_user_id int32, mutual int8, owner_user_id int32, contact_phone string) int64 {
	var query = "update user_contacts set contact_user_id = ?, mutual = ? where contact_user_id = 0 and owner_user_id = ? and contact_phone = ?"
	r, err := dao.db.Exec(query, contact_user_id, mutual, owner_user_id, contact_phone)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateContactUserId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateContactUserId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update user_contacts set mutual = :mutual where contact_user_id != 0 and (owner_user_id = :owner_user_id and contact_user_id = :contact_user_id)
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) UpdateMutual(mutual int8, owner_user_id int32, contact_user_id int32) int64 {
	var query = "update user_contacts set mutual = ? where contact_user_id != 0 and (owner_user_id = ? and contact_user_id = ?)"
	r, err := dao.db.Exec(query, mutual, owner_user_id, contact_user_id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateMutual(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateMutual(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update user_contacts set mutual = :mutual where id = :id
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) UpdateMutualByID(mutual int8, id int64) int64 {
	var query = "update user_contacts set mutual = ? where id = ?"
	r, err := dao.db.Exec(query, mutual, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateMutualByID(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateMutualByID(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update user_contacts set is_blocked = :is_blocked where contact_user_id != 0 and (owner_user_id = :owner_user_id and contact_user_id = :contact_user_id)
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) UpdateBlock(is_blocked int8, owner_user_id int32, contact_user_id int32) int64 {
	var query = "update user_contacts set is_blocked = ? where contact_user_id != 0 and (owner_user_id = ? and contact_user_id = ?)"
	r, err := dao.db.Exec(query, is_blocked, owner_user_id, contact_user_id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateBlock(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateBlock(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update user_contacts set is_deleted = 1, mutual = 0 where contact_user_id != 0 and (owner_user_id = :owner_user_id and contact_user_id in (:id_list))
// TODO(@benqi): sqlmap
func (dao *UserContactsDAO) DeleteContacts(owner_user_id int32, id_list []int32) int64 {
	var q = "update user_contacts set is_deleted = 1, mutual = 0 where contact_user_id != 0 and (owner_user_id = ? and contact_user_id in (?))"
	query, a, err := sqlx.In(q, owner_user_id, id_list)
	r, err := dao.db.Exec(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in DeleteContacts(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in DeleteContacts(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}
