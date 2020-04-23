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
	"github.com/nebulaim/telegramd/service/channel/biz/dal/dataobject"
)

type ChannelParticipantsDAO struct {
	db *sqlx.DB
}

func NewChannelParticipantsDAO(db *sqlx.DB) *ChannelParticipantsDAO {
	return &ChannelParticipantsDAO{db}
}

// insert into channel_participants(channel_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state) values (:channel_id, :user_id, :participant_type, :inviter_user_id, :invited_at, :joined_at, :state)
// TODO(@benqi): sqlmap
func (dao *ChannelParticipantsDAO) Insert(do *dataobject.ChannelParticipantsDO) int64 {
	var query = "insert into channel_participants(channel_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state) values (:channel_id, :user_id, :participant_type, :inviter_user_id, :invited_at, :joined_at, :state)"
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

// select id, channel_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state from channel_participants where channel_id = :channel_id
// TODO(@benqi): sqlmap
func (dao *ChannelParticipantsDAO) SelectByChannelId(channel_id int32) []dataobject.ChannelParticipantsDO {
	var query = "select id, channel_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state from channel_participants where channel_id = ?"
	rows, err := dao.db.Queryx(query, channel_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByChannelId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.ChannelParticipantsDO
	for rows.Next() {
		v := dataobject.ChannelParticipantsDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByChannelId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByChannelId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select id, channel_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state from channel_participants where channel_id = :channel_id and user_id in (:idList)
// TODO(@benqi): sqlmap
func (dao *ChannelParticipantsDAO) SelectByUserIdList(channel_id int32, idList []int32) []dataobject.ChannelParticipantsDO {
	var q = "select id, channel_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state from channel_participants where channel_id = ? and user_id in (?)"
	query, a, err := sqlx.In(q, channel_id, idList)
	rows, err := dao.db.Queryx(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByUserIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.ChannelParticipantsDO
	for rows.Next() {
		v := dataobject.ChannelParticipantsDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByUserIdList(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByUserIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select id, channel_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state from channel_participants where channel_id = :channel_id and user_id = :user_id
// TODO(@benqi): sqlmap
func (dao *ChannelParticipantsDAO) SelectByUserId(channel_id int32, user_id int32) *dataobject.ChannelParticipantsDO {
	var query = "select id, channel_id, user_id, participant_type, inviter_user_id, invited_at, joined_at, state from channel_participants where channel_id = ? and user_id = ?"
	rows, err := dao.db.Queryx(query, channel_id, user_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByUserId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.ChannelParticipantsDO{}
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

// update channel_participants set state = 1 where id = :id
// TODO(@benqi): sqlmap
func (dao *ChannelParticipantsDAO) DeleteChannelUser(id int32) int64 {
	var query = "update channel_participants set state = 1 where id = ?"
	r, err := dao.db.Exec(query, id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in DeleteChannelUser(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in DeleteChannelUser(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update channel_participants set inviter_user_id = :inviter_user_id, invited_at = :invited_at, joined_at = :joined_at, state = 0 where id = :id
// TODO(@benqi): sqlmap
func (dao *ChannelParticipantsDAO) Update(inviter_user_id int32, invited_at int32, joined_at int32, id int32) int64 {
	var query = "update channel_participants set inviter_user_id = ?, invited_at = ?, joined_at = ?, state = 0 where id = ?"
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

// update channel_participants set participant_type = :participant_type where id = :id
// TODO(@benqi): sqlmap
func (dao *ChannelParticipantsDAO) UpdateParticipantType(participant_type int8, id int32) int64 {
	var query = "update channel_participants set participant_type = ? where id = ?"
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
