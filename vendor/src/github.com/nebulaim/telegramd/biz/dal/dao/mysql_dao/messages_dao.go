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

type MessagesDAO struct {
	db *sqlx.DB
}

func NewMessagesDAO(db *sqlx.DB) *MessagesDAO {
	return &MessagesDAO{db}
}

// insert into message_boxes(user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, random_id, message_type, message_data, date2) values (:user_id, :user_message_box_id, :dialog_message_id, :sender_user_id, :message_box_type, :peer_type, :peer_id, :random_id, :message_type, :message_data, :date2)
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) Insert(do *dataobject.MessagesDO) int64 {
	var query = "insert into message_boxes(user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, random_id, message_type, message_data, date2) values (:user_id, :user_message_box_id, :dialog_message_id, :sender_user_id, :message_box_type, :peer_type, :peer_id, :random_id, :message_type, :message_data, :date2)"
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

// select user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = :user_id and deleted = 0 and user_message_box_id in (:idList) order by user_message_box_id desc
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectByMessageIdList(user_id int32, idList []int32) []dataobject.MessagesDO {
	var q = "select user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = ? and deleted = 0 and user_message_box_id in (?) order by user_message_box_id desc"
	query, a, err := sqlx.In(q, user_id, idList)
	rows, err := dao.db.Queryx(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByMessageIdList(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = :user_id and user_message_box_id = :user_message_box_id and deleted = 0 limit 1
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectByMessageId(user_id int32, user_message_box_id int32) *dataobject.MessagesDO {
	var query = "select user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = ? and user_message_box_id = ? and deleted = 0 limit 1"
	rows, err := dao.db.Queryx(query, user_id, user_message_box_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.MessagesDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByMessageId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectByMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = :user_id and peer_type = :peer_type and peer_id = :peer_id and user_message_box_id < :user_message_box_id and deleted = 0 order by user_message_box_id desc limit :limit
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectBackwardByPeerOffsetLimit(user_id int32, peer_type int8, peer_id int32, user_message_box_id int32, limit int32) []dataobject.MessagesDO {
	var query = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = ? and peer_type = ? and peer_id = ? and user_message_box_id < ? and deleted = 0 order by user_message_box_id desc limit ?"
	rows, err := dao.db.Queryx(query, user_id, peer_type, peer_id, user_message_box_id, limit)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectBackwardByPeerOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectBackwardByPeerOffsetLimit(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectBackwardByPeerOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = :user_id and ((sender_user_id = :user_id and peer_id = :peer_id) or (sender_user_id = :peer_id and peer_id = :user_id)) and peer_type = :peer_type and user_message_box_id < :user_message_box_id and deleted = 0 order by user_message_box_id desc limit :limit
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectBackwardByPeerUserOffsetLimit(user_id int32, peer_id int32, peer_type int8, user_message_box_id int32, limit int32) []dataobject.MessagesDO {
	var query = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = ? and ((sender_user_id = ? and peer_id = ?) or (sender_user_id = ? and peer_id = ?)) and peer_type = ? and user_message_box_id < ? and deleted = 0 order by user_message_box_id desc limit ?"
	rows, err := dao.db.Queryx(query, user_id, user_id, peer_id, peer_id, user_id, peer_type, user_message_box_id, limit)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectBackwardByPeerUserOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectBackwardByPeerUserOffsetLimit(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectBackwardByPeerUserOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = :user_id and peer_type = :peer_type and peer_id = :peer_id and user_message_box_id >= :user_message_box_id and deleted = 0 order by user_message_box_id asc limit :limit
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectForwardByPeerOffsetLimit(user_id int32, peer_type int8, peer_id int32, user_message_box_id int32, limit int32) []dataobject.MessagesDO {
	var query = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = ? and peer_type = ? and peer_id = ? and user_message_box_id >= ? and deleted = 0 order by user_message_box_id asc limit ?"
	rows, err := dao.db.Queryx(query, user_id, peer_type, peer_id, user_message_box_id, limit)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectForwardByPeerOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectForwardByPeerOffsetLimit(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectForwardByPeerOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = :user_id and ((sender_user_id = :user_id and peer_id = :peer_id) or (sender_user_id = :peer_id and peer_id = :user_id)) and peer_type = :peer_type and user_message_box_id >= :user_message_box_id and deleted = 0 order by user_message_box_id asc limit :limit
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectForwardByPeerUserOffsetLimit(user_id int32, peer_id int32, peer_type int8, user_message_box_id int32, limit int32) []dataobject.MessagesDO {
	var query = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id = ? and ((sender_user_id = ? and peer_id = ?) or (sender_user_id = ? and peer_id = ?)) and peer_type = ? and user_message_box_id >= ? and deleted = 0 order by user_message_box_id asc limit ?"
	rows, err := dao.db.Queryx(query, user_id, user_id, peer_id, peer_id, user_id, peer_type, user_message_box_id, limit)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectForwardByPeerUserOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectForwardByPeerUserOffsetLimit(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectForwardByPeerUserOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select user_message_box_id, message_box_type from message_boxes where user_id = :peerId and dialog_message_id = (select dialog_message_id from messages where user_id = :user_id and user_message_box_id = :user_message_box_id and deleted = 0 limit 1)
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectPeerMessageId(peerId int32, user_id int32, user_message_box_id int32) *dataobject.MessagesDO {
	var query = "select user_message_box_id, message_box_type from message_boxes where user_id = ? and dialog_message_id = (select dialog_message_id from messages where user_id = ? and user_message_box_id = ? and deleted = 0 limit 1)"
	rows, err := dao.db.Queryx(query, peerId, user_id, user_message_box_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectPeerMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.MessagesDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectPeerMessageId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectPeerMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, dialog_message_id, message_type from message_boxes where user_id != :user_id and dialog_message_id in (select dialog_message_id from messages where user_id = :user_id and user_message_box_id in (:idList)) and deleted = 0
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectPeerDialogMessageIdList(user_id int32, idList []int32) []dataobject.MessagesDO {
	var q = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, dialog_message_id, message_type from message_boxes where user_id != ? and dialog_message_id in (select dialog_message_id from messages where user_id = ? and user_message_box_id in (?)) and deleted = 0"
	query, a, err := sqlx.In(q, user_id, user_id, idList)
	rows, err := dao.db.Queryx(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectPeerDialogMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectPeerDialogMessageIdList(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectPeerDialogMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where dialog_message_id = (select dialog_message_id from messages where user_id = :user_id and user_message_box_id = :user_message_box_id) and deleted = 0
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectDialogMessageListByMessageId(user_id int32, user_message_box_id int32) []dataobject.MessagesDO {
	var query = "select user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where dialog_message_id = (select dialog_message_id from messages where user_id = ? and user_message_box_id = ?) and deleted = 0"
	rows, err := dao.db.Queryx(query, user_id, user_message_box_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectDialogMessageListByMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectDialogMessageListByMessageId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectDialogMessageListByMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id != :user_id and dialog_message_id = (select dialog_message_id from messages where user_id = :user_id and user_message_box_id = :user_message_box_id) and deleted = 0
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectPeerDialogMessageListByMessageId(user_id int32, user_message_box_id int32) []dataobject.MessagesDO {
	var query = "select user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from message_boxes where user_id != ? and dialog_message_id = (select dialog_message_id from messages where user_id = ? and user_message_box_id = ?) and deleted = 0"
	rows, err := dao.db.Queryx(query, user_id, user_id, user_message_box_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectPeerDialogMessageListByMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectPeerDialogMessageListByMessageId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectPeerDialogMessageListByMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// select user_message_box_id from message_boxes where user_id = :user_id order by user_message_box_id desc limit 2
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectLastTwoMessageId(user_id int32) *dataobject.MessagesDO {
	var query = "select user_message_box_id from message_boxes where user_id = ? order by user_message_box_id desc limit 2"
	rows, err := dao.db.Queryx(query, user_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectLastTwoMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.MessagesDO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectLastTwoMessageId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectLastTwoMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return do
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, dialog_message_id, message_type from message_boxes where user_id = :user_id and user_message_box_id in (:idList) and deleted = 0
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectDialogsByMessageIdList(user_id int32, idList []int32) []dataobject.MessagesDO {
	var q = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, dialog_message_id, message_type from message_boxes where user_id = ? and user_message_box_id in (?) and deleted = 0"
	query, a, err := sqlx.In(q, user_id, idList)
	rows, err := dao.db.Queryx(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectDialogsByMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectDialogsByMessageIdList(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectDialogsByMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}

// update message_boxes set deleted = 1 where user_id = :user_id and user_message_box_id in (:idList) and deleted = 0
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) DeleteMessagesByMessageIdList(user_id int32, idList []int32) int64 {
	var q = "update message_boxes set deleted = 1 where user_id = ? and user_message_box_id in (?) and deleted = 0"
	query, a, err := sqlx.In(q, user_id, idList)
	r, err := dao.db.Exec(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in DeleteMessagesByMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in DeleteMessagesByMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// update message_boxes set message_data = :message_data where user_id = :user_id and user_message_box_id = :user_message_box_id and deleted = 0
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) UpdateMessagesData(message_data string, user_id int32, user_message_box_id int32) int64 {
	var query = "update message_boxes set message_data = ? where user_id = ? and user_message_box_id = ? and deleted = 0"
	r, err := dao.db.Exec(query, message_data, user_id, user_message_box_id)

	if err != nil {
		errDesc := fmt.Sprintf("Exec in UpdateMessagesData(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	rows, err := r.RowsAffected()
	if err != nil {
		errDesc := fmt.Sprintf("RowsAffected in UpdateMessagesData(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return rows
}

// select user_message_box_id from message_boxes where user_id = :user_id and ((sender_user_id = :user_id and peer_id = :peer_id) or (sender_user_id = :peer_id and peer_id = :user_id)) and peer_type = :peer_type and deleted = 0
// TODO(@benqi): sqlmap
func (dao *MessagesDAO) SelectDialogMessageIdList(user_id int32, peer_id int32, peer_type int8) []dataobject.MessagesDO {
	var query = "select user_message_box_id from message_boxes where user_id = ? and ((sender_user_id = ? and peer_id = ?) or (sender_user_id = ? and peer_id = ?)) and peer_type = ? and deleted = 0"
	rows, err := dao.db.Queryx(query, user_id, user_id, peer_id, peer_id, user_id, peer_type)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectDialogMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.MessagesDO
	for rows.Next() {
		v := dataobject.MessagesDO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectDialogMessageIdList(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	err = rows.Err()
	if err != nil {
		errDesc := fmt.Sprintf("rows in SelectDialogMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	return values
}
