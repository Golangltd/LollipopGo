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

package account

import (
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
)

func (m *AccountModel) InsertReportData(userId, peerType, peerId, reason int32, text string) bool {
	do := &dataobject.ReportsDO{
		UserId:   userId,
		PeerType: peerType,
		PeerId:   peerId,
		Reason:   int8(reason),
		Content:  text,
	}
	do.Id = m.dao.ReportsDAO.Insert(do)
	return do.Id > 0
}
