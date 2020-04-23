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

type UnregisteredContactsDAO struct {
	db *sqlx.DB
}

func NewUnregisteredContactsDAO(db *sqlx.DB) *UnregisteredContactsDAO {
	return &UnregisteredContactsDAO{db}
}

// insert into unregistered_contacts(importer_user_id, contact_phone, contact_first_name, contact_last_name) values (:importer_user_id, :contact_phone, :contact_first_name, :contact_last_name)
// TODO(@benqi): sqlmap
func (dao *UnregisteredContactsDAO) Insert(do *dataobject.UnregisteredContactsDO) int64 {
	var query = "insert into unregistered_contacts(importer_user_id, contact_phone, contact_first_name, contact_last_name) values (:importer_user_id, :contact_phone, :contact_first_name, :contact_last_name)"
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

// select id, importer_user_id, contact_phone, contact_first_name, contact_last_name from unregistered_contacts where importer_user_id = :importer_user_id and is_deleted = 0 and contact_phone in (:phone_list)
// TODO(@benqi): sqlmap
func (dao *UnregisteredContactsDAO) SelectContacts(importer_user_id int32, phone_list []string) []dataobject.UnregisteredContactsDO {
	var q = "select id, importer_user_id, contact_phone, contact_first_name, contact_last_name from unregistered_contacts where importer_user_id = ? and is_deleted = 0 and contact_phone in (?)"
	query, a, err := sqlx.In(q, importer_user_id, phone_list)
	rows, err := dao.db.Queryx(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectContacts(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.UnregisteredContactsDO
	for rows.Next() {
		v := dataobject.UnregisteredContactsDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectContacts(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectContacts(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// update unregistered_contacts set contact_first_name = :contact_first_name, contact_last_name = :contact_last_name where id = :id
// TODO(@benqi): sqlmap
func (dao *UnregisteredContactsDAO) UpdateContactName(contact_first_name string, contact_last_name string, id int32) int64 {
	var query = "update unregistered_contacts set contact_first_name = ?, contact_last_name = ? where id = ?"
	r, err := dao.db.Exec(query, contact_first_name, contact_last_name, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateContactName(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateContactName(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update unregistered_contacts set is_deleted = 1 where id in (:id_list)
// TODO(@benqi): sqlmap
func (dao *UnregisteredContactsDAO) DeleteContacts(id_list []int64) int64 {
	var q = "update unregistered_contacts set is_deleted = 1 where id in (?)"
	query, a, err := sqlx.In(q, id_list)
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
