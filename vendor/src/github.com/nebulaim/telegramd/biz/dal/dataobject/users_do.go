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

type UsersDO struct {
	Id             int32  `db:"id"`
	AccessHash     int64  `db:"access_hash"`
	FirstName      string `db:"first_name"`
	LastName       string `db:"last_name"`
	Username       string `db:"username"`
	Phone          string `db:"phone"`
	CountryCode    string `db:"country_code"`
	Bio            string `db:"bio"`
	About          string `db:"about"`
	State          int32  `db:"state"`
	IsBot          int8   `db:"is_bot"`
	Banned         int64  `db:"banned"`
	BannedReason   string `db:"banned_reason"`
	AccountDaysTtl int32  `db:"account_days_ttl"`
	Photos         string `db:"photos"`
	Deleted        int8   `db:"deleted"`
	DeletedReason  string `db:"deleted_reason"`
	CreatedAt      string `db:"created_at"`
	UpdatedAt      string `db:"updated_at"`
	BannedAt       string `db:"banned_at"`
	DeletedAt      string `db:"deleted_at"`
}
