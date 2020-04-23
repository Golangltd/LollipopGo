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

type ChatParticipantsDAO struct {
	db *sqlx.DB
}

func NewChatParticipantsDAO(db *sqlx.DB) *ChatParticipantsDAO {
	return &ChatParticipantsDAO{db}
}

// insert into chat_participants(chat_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state) values (:chat_id, :user_id, :participant_type, :inviter_user_id, :invited_at, :joined_at, :state)
// TODO(@benqi): sqlmap
func (dao *ChatParticipantsDAO) Insert(do *dataobject.ChatParticipantsDO) int64 {
	var query = "insert into chat_participants(chat_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state) values (:chat_id, :user_id, :participant_type, :inviter_user_id, :invited_at, :joined_at, :state)"
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

// select id, chat_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state from chat_participants where chat_id = :chat_id
// TODO(@benqi): sqlmap
func (dao *ChatParticipantsDAO) SelectByChatId(chat_id int32) []dataobject.ChatParticipantsDO {
	var query = "select id, chat_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state from chat_participants where chat_id = ?"
	rows, err := dao.db.Queryx(query, chat_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByChatId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.ChatParticipantsDO
	for rows.Next() {
		v := dataobject.ChatParticipantsDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByChatId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByChatId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// update chat_participants set state = 1 where id = :id
// TODO(@benqi): sqlmap
func (dao *ChatParticipantsDAO) DeleteChatUser(id int32) int64 {
	var query = "update chat_participants set state = 1 where id = ?"
	r, err := dao.db.Exec(query, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in DeleteChatUser(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in DeleteChatUser(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update chat_participants set inviter_user_id = :inviter_user_id, invited_at = :invited_at, joined_at = :joined_at, state = 0 where id = :id
// TODO(@benqi): sqlmap
func (dao *ChatParticipantsDAO) Update(inviter_user_id int32, invited_at int32, joined_at int32, id int32) int64 {
	var query = "update chat_participants set inviter_user_id = ?, invited_at = ?, joined_at = ?, state = 0 where id = ?"
	r, err := dao.db.Exec(query, inviter_user_id, invited_at, joined_at, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in Update(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in Update(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update chat_participants set participant_type = :participant_type where id = :id
// TODO(@benqi): sqlmap
func (dao *ChatParticipantsDAO) UpdateParticipantType(participant_type int8, id int32) int64 {
	var query = "update chat_participants set participant_type = ? where id = ?"
	r, err := dao.db.Exec(query, participant_type, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateParticipantType(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateParticipantType(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}
