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
	"github.com/nebulaim/telegramd/baselib/base"
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"time"
)

//type accountData struct {
//	userId int32
//	banned bool
//}
//
//func (this *accountData) IsBanned() bool {
//	return this.banned
//}

type userData struct {
	*mtproto.TLUser
}

func (this *userData) ToUser() *mtproto.User {
	return this.TLUser.To_User()
}

func (m *UserModel) GetUsersBySelfAndIDList(selfUserId int32, userIdList []int32) (users []*mtproto.User) {
	if len(userIdList) == 0 {
		users = []*mtproto.User{}
	} else {
		// TODO(@benqi):  需要优化，makeUserDataByDO需要查询用户状态以及获取Mutual和Contact状态信息而导致多次查询
		userDOList := m.dao.UsersDAO.SelectUsersByIdList(userIdList)
		users = make([]*mtproto.User, 0, len(userDOList))
		for i := 0; i < len(userDOList); i++ {
			user := m.makeUserDataByDO(selfUserId, &userDOList[i])
			//
			//// TODO(@benqi): fill bot, photo, about...
			//user := &mtproto.TLUser{Data2: &mtproto.User_Data{
			//	Self:          selfUserId == userDO.Id,
			//	Id:            userDO.Id,
			//	AccessHash:    userDO.AccessHash,
			//	FirstName:     userDO.FirstName,
			//	LastName:      userDO.LastName,
			//	Username:      userDO.Username,
			//	Phone:         userDO.Phone,
			//	Contact:       true,
			//	MutualContact: true,
			//}}
			users = append(users, user.To_User())
		}
	}
	return
}

func (m *UserModel) GetUserFull(selfUserId int32, userId int32) (userFull *mtproto.TLUserFull) {
	//TODO(@benqi): 等Link和NotifySettings实现后再来完善
	//fullUser := &mtproto.TLUserFull{}
	//fullUser.PhoneCallsAvailable = true
	//fullUser.PhoneCallsPrivate = true
	//fullUser.About = "@Benqi"
	//fullUser.CommonChatsCount = 0
	//
	//switch request.Id.Payload.(type) {
	//case *mtproto.InputUser_InputUserSelf:
	//	// User
	//	userDO, _ := s.UsersDAO.SelectById(2)
	//	user := &mtproto.TLUser{}
	//	user.Self = true
	//	user.Contact = false
	//	user.Id = userDO.Id
	//	user.FirstName = userDO.FirstName
	//	user.LastName = userDO.LastName
	//	user.Username = userDO.Username
	//	user.AccessHash = userDO.AccessHash
	//	user.Phone = userDO.Phone
	//	fullUser.User = mtproto.MakeUser(user)
	//
	//	// Link
	//	link := &mtproto.TLContactsLink{}
	//	link.MyLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.ForeignLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.User = mtproto.MakeUser(user)
	//	fullUser.Link = mtproto.MakeContacts_Link(link)
	//case *mtproto.InputUser_InputUser:
	//	inputUser := request.Id.Payload.(*mtproto.InputUser_InputUser).InputUser
	//	// User
	//	userDO, _ := s.UsersDAO.SelectById(inputUser.UserId)
	//	user := &mtproto.TLUser{}
	//	user.Self = false
	//	user.Contact = true
	//	user.Id = userDO.Id
	//	user.FirstName = userDO.FirstName
	//	user.LastName = userDO.LastName
	//	user.Username = userDO.Username
	//	user.AccessHash = userDO.AccessHash
	//	user.Phone = userDO.Phone
	//	fullUser.User = mtproto.MakeUser(user)
	//
	//	// Link
	//	link := &mtproto.TLContactsLink{}
	//	link.MyLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.ForeignLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.User = mtproto.MakeUser(user)
	//	fullUser.Link = mtproto.MakeContacts_Link(link)
	//case *mtproto.InputUser_InputUserEmpty:
	//	// TODO(@benqi): BAD_REQUEST: 400
	//}
	//
	//// NotifySettings
	//peerNotifySettings := &mtproto.TLPeerNotifySettings{}
	//peerNotifySettings.ShowPreviews = true
	//peerNotifySettings.MuteUntil = 0
	//peerNotifySettings.Sound = "default"
	//fullUser.NotifySettings = mtproto.MakePeerNotifySettings(peerNotifySettings)
	return nil
}

func (m *UserModel) UpdateUserStatus(userId int32, lastSeenAt int64) {
	// now := time.Now().Unix()
	rows := m.dao.UserPresencesDAO.UpdateLastSeen(lastSeenAt, 0, userId)
	if rows == 0 {
		do := &dataobject.UserPresencesDO{
			UserId:            userId,
			LastSeenAt:        lastSeenAt,
			LastSeenAuthKeyId: 0,
			LastSeenIp:        "",
			CreatedAt:         base.NowFormatYMDHMS(),
		}
		m.dao.UserPresencesDAO.Insert(do)
	}
}

func (m *UserModel) GetUserStatus(userId int32) *mtproto.UserStatus {
	now := time.Now().Unix()
	do := m.dao.UserPresencesDAO.SelectByUserID(userId)
	if do == nil {
		return mtproto.NewTLUserStatusEmpty().To_UserStatus()
	}

	if now <= do.LastSeenAt+5*60 {
		status := &mtproto.TLUserStatusOnline{Data2: &mtproto.UserStatus_Data{
			Expires: int32(do.LastSeenAt + 5*30),
		}}
		return status.To_UserStatus()
	} else {
		status := &mtproto.TLUserStatusOffline{Data2: &mtproto.UserStatus_Data{
			WasOnline: int32(do.LastSeenAt),
		}}
		return status.To_UserStatus()
	}
}

func (m *UserModel) DeleteUser(userId int32, reason string) bool {
	affected := m.dao.UsersDAO.Delete(reason, base.NowFormatYMDHMS(), userId)
	return affected == 1
}
