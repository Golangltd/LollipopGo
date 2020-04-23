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

package user

import (
	"encoding/json"
)

// type profileData *ProfilePhotoIds

func MakeProfilePhotoData(jsonData string) *ProfilePhotoIds {
	if jsonData == "" {
		return &ProfilePhotoIds{}
	}
	data2 := &ProfilePhotoIds{}
	err := json.Unmarshal([]byte(jsonData), data2)
	if err != nil {
		return &ProfilePhotoIds{}
	}
	return data2
}

func (m *ProfilePhotoIds) AddPhotoId(id int64) {
	idList := make([]int64, 0, len(m.IdList))
	idList = append(idList, id)
	idList = append(idList, m.IdList...)
	m.IdList = idList
	m.Default = id
}

func (m *ProfilePhotoIds) RemovePhotoId(id int64) int64 {
	if len(m.IdList) <= 1 {
		m.IdList = []int64{}
		m.Default = 0
	} else {
		if id == m.Default {
			id = m.IdList[1]
			m.IdList = m.IdList[1:]
		} else {
			for i, j := range m.IdList {
				if j == id {
					m.IdList = append(m.IdList[:i], m.IdList[i+1:]...)
				}
			}
		}
	}
	return m.Default
}

func (m *ProfilePhotoIds) ToJson() string {
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}
