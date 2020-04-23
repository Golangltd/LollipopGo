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

type ChatsDAO struct {
	db *sqlx.DB
}

func NewChatsDAO(db *sqlx.DB) *ChatsDAO {
	return &ChatsDAO{db}
}

// insert into chats(creator_user_id, access_hash, random_id, participant_count, title, `date`) values (:creator_user_id, :access_hash, :random_id, :participant_count, :title, :date)
// TODO(@benqi): sqlmap
func (dao *ChatsDAO) Insert(do *dataobject.ChatsDO) int64 {
	var query = "insert into chats(creator_user_id, access_hash, random_id, participant_count, title, `date`) values (:creator_user_id, :access_hash, :random_id, :participant_count, :title, :date)"
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

// select id, creator_user_id, access_hash, participant_count, title, photo_id, admins_enabled, deactivated, version, `date` from chats where id = :id
// TODO(@benqi): sqlmap
func (dao *ChatsDAO) Select(id int32) *dataobject.ChatsDO {
	var query = "select id, creator_user_id, access_hash, participant_count, title, photo_id, admins_enabled, deactivated, version, `date` from chats where id = ?"
	rows, err := dao.db.Queryx(query, id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in Select(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.ChatsDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in Select(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in Select(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// update chats set title = :title, `date` = :date, version = version + 1 where id = :id
// TODO(@benqi): sqlmap
func (dao *ChatsDAO) UpdateTitle(title string, date int32, id int32) int64 {
	var query = "update chats set title = ?, `date` = ?, version = version + 1 where id = ?"
	r, err := dao.db.Exec(query, title, date, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateTitle(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateTitle(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// select id, access_hash, participant_count, title, photo_id, admins_enabled, deactivated, version, `date` from chats where id in (:idList)
// TODO(@benqi): sqlmap
func (dao *ChatsDAO) SelectByIdList(idList []int32) []dataobject.ChatsDO {
	var q = "select id, access_hash, participant_count, title, photo_id, admins_enabled, deactivated, version, `date` from chats where id in (?)"
	query, a, err := sqlx.In(q, idList)
	rows, err := dao.db.Queryx(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.ChatsDO
	for rows.Next() {
		v := dataobject.ChatsDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByIdList(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// update chats set participant_count = :participant_count, `date` = :date, version = version + 1 where id = :id
// TODO(@benqi): sqlmap
func (dao *ChatsDAO) UpdateParticipantCount(participant_count int32, date int32, id int32) int64 {
	var query = "update chats set participant_count = ?, `date` = ?, version = version + 1 where id = ?"
	r, err := dao.db.Exec(query, participant_count, date, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateParticipantCount(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateParticipantCount(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update chats set photo_id = :photo_id, `date` = :date, version = version + 1 where id = :id
// TODO(@benqi): sqlmap
func (dao *ChatsDAO) UpdatePhotoId(photo_id int64, date int32, id int32) int64 {
	var query = "update chats set photo_id = ?, `date` = ?, version = version + 1 where id = ?"
	r, err := dao.db.Exec(query, photo_id, date, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdatePhotoId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdatePhotoId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update chats set admins_enabled = :admins_enabled, `date` = :date, version = version + 1 where id = :id
// TODO(@benqi): sqlmap
func (dao *ChatsDAO) UpdateAdminsEnabled(admins_enabled int8, date int32, id int32) int64 {
	var query = "update chats set admins_enabled = ?, `date` = ?, version = version + 1 where id = ?"
	r, err := dao.db.Exec(query, admins_enabled, date, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateAdminsEnabled(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateAdminsEnabled(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update chats set `date` = :date, version = version + 1 where id = :id
// TODO(@benqi): sqlmap
func (dao *ChatsDAO) UpdateVersion(date int32, id int32) int64 {
	var query = "update chats set `date` = ?, version = version + 1 where id = ?"
	r, err := dao.db.Exec(query, date, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateVersion(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateVersion(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}
