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

type UserPasswordsDO struct {
	Id          int64  `db:"id"`
	UserId      int32  `db:"user_id"`
	ServerSalt  string `db:"server_salt"`
	Hash        string `db:"hash"`
	Salt        string `db:"salt"`
	Hint        string `db:"hint"`
	Email       string `db:"email"`
	HasRecovery int8   `db:"has_recovery"`
	Code        string `db:"code"`
	CodeExpired int32  `db:"code_expired"`
	Attempts    int32  `db:"attempts"`
	State       int8   `db:"state"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}
