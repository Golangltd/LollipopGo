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

package dataobject

type MessagesDO struct {
	Id               int32  `db:"id"`
	UserId           int32  `db:"user_id"`
	UserMessageBoxId int32  `db:"user_message_box_id"`
	DialogMessageId  int64  `db:"dialog_message_id"`
	SenderUserId     int32  `db:"sender_user_id"`
	MessageBoxType   int8   `db:"message_box_type"`
	PeerType         int8   `db:"peer_type"`
	PeerId           int32  `db:"peer_id"`
	RandomId         int64  `db:"random_id"`
	MessageType      int8   `db:"message_type"`
	MessageData      string `db:"message_data"`
	Date2            int32  `db:"date2"`
	Deleted          int8   `db:"deleted"`
	CreatedAt        string `db:"created_at"`
	UpdatedAt        string `db:"updated_at"`
}
