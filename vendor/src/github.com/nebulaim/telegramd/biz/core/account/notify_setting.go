/*
 *  Copyright (c) 2017, https://github.com/nebulaim
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
	base2 "github.com/nebulaim/telegramd/baselib/base"
	"github.com/nebulaim/telegramd/biz/base"
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
	// "github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/proto/mtproto"
)

func (m *AccountModel) GetNotifySettings(userId int32, peer *base.PeerUtil) *mtproto.PeerNotifySettings {
	do := m.dao.UserNotifySettingsDAO.SelectByPeer(userId, int8(peer.PeerType), peer.PeerId)

	// var mute_until int32 = 0
	if do == nil {
		settings := &mtproto.TLPeerNotifySettings{Data2: &mtproto.PeerNotifySettings_Data{
			//ShowPreviews: mtproto.ToBool(false),
			//Silent:       mtproto.ToBool(false),
			MuteUntil:      1,
			Sound:        "default",
		}}
		return settings.To_PeerNotifySettings()
	} else {
		settings := mtproto.NewTLPeerNotifySettings()
		if do.ShowPreviews == 1 {
			settings.SetShowPreviews(mtproto.ToBool(true))
		}
		if do.Silent == 1 {
			settings.SetSilent(mtproto.ToBool(true))
		}
		if do.MuteUntil == 0 {
			settings.SetMuteUntil(1)
		} else {
			settings.SetMuteUntil(do.MuteUntil)
		}
		if do.Sound == "" {
			settings.SetSound("default")
		} else {
			settings.SetSound(do.Sound)
		}
		return settings.To_PeerNotifySettings()
	}
}

func (m *AccountModel) SetNotifySettings(userId int32, peer *base.PeerUtil, settings *mtproto.TLInputPeerNotifySettings) {
	var (
		showPreviews = base2.BoolToInt8(mtproto.FromBool(settings.GetShowPreviews()))
		silent       = base2.BoolToInt8(mtproto.FromBool(settings.GetSilent()))
	)

	do := m.dao.UserNotifySettingsDAO.SelectByPeer(userId, int8(peer.PeerType), peer.PeerId)
	if do == nil {
		do = &dataobject.UserNotifySettingsDO{
			UserId:       userId,
			PeerType:     int8(peer.PeerType),
			PeerId:       peer.PeerId,
			ShowPreviews: showPreviews,
			Silent:       silent,
			MuteUntil:    settings.GetMuteUntil(),
			Sound:        settings.GetSound(),
		}
		m.dao.UserNotifySettingsDAO.Insert(do)
	} else {
		m.dao.UserNotifySettingsDAO.UpdateByPeer(showPreviews, silent, settings.GetMuteUntil(), settings.GetSound(), 0, userId, int8(peer.PeerType), peer.PeerId)
	}
}

func (m *AccountModel) ResetNotifySettings(userId int32) {
	m.dao.UserNotifySettingsDAO.DeleteNotAll(userId)
	do := m.dao.UserNotifySettingsDAO.SelectByAll(userId)
	if do == nil {
		do = &dataobject.UserNotifySettingsDO{}
		do.UserId = userId
		// do.PeerType = base.PEER_ALL
		do.PeerId = 0
		do.ShowPreviews = 1
		do.Silent = 0
		do.MuteUntil = 0
		m.dao.UserNotifySettingsDAO.Insert(do)
	} else {
		m.dao.UserNotifySettingsDAO.UpdateByPeer(1, 0, 0, "default", 0, userId, base.PEER_USER, 0)
	}
}
