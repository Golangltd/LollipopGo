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

type TmpPasswordsDO struct {
	Id           int32  `db:"id"`
	AuthId       int64  `db:"auth_id"`
	UserId       int32  `db:"user_id"`
	PasswordHash string `db:"password_hash"`
	Period       int32  `db:"period"`
	TmpPassword  string `db:"tmp_password"`
	ValidUntil   int32  `db:"valid_until"`
	CreatedAt    string `db:"created_at"`
}
