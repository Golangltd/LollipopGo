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
	"encoding/hex"
	"github.com/nebulaim/telegramd/baselib/crypto"
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"math/rand"
	"time"
)

func CheckUserAccessHash(id int32, hash int64) bool {
	return true
}

func (m *UserModel) CheckPhoneNumberExist(phoneNumber string) bool {
	return nil != m.dao.UsersDAO.SelectByPhoneNumber(phoneNumber)
}

func makeUserStatusOnline() *mtproto.UserStatus {
	now := time.Now().Unix()
	status := &mtproto.UserStatus{
		Constructor: mtproto.TLConstructor_CRC32_userStatusOnline,
		Data2: &mtproto.UserStatus_Data{
			// WasOnline: int32(now),
			Expires: int32(now + 60),
		},
	}
	return status
}

func (m *UserModel) makeUserDataByDO(selfId int32, do *dataobject.UsersDO) *userData {
	if do == nil {
		return nil
	} else {
		var (
			status                 *mtproto.UserStatus
			photo                  *mtproto.UserProfilePhoto
			phone                  string
			contact, mutualContact bool
			isSelf                 = selfId == do.Id
		)

		if isSelf {
			status = makeUserStatusOnline()
			contact = true
			mutualContact = true
			phone = do.Phone
		} else {
			status = m.GetUserStatus(do.Id)
			contact, mutualContact = m.contactCallback.GetContactAndMutual(selfId, do.Id)
			// if contact {
			phone = do.Phone
			// }
		}

		photoId := m.GetDefaultUserPhotoID(do.Id)
		if photoId == 0 {
			photo = mtproto.NewTLUserProfilePhotoEmpty().To_UserProfilePhoto()
		} else {
			photo = m.photoCallback.GetUserProfilePhoto(photoId)
			//sizeList, _ := nbfs_client.GetPhotoSizeList(photoId)
			//photo = photo2.MakeUserProfilePhoto(photoId, sizeList)
		}

		data := &userData{TLUser: &mtproto.TLUser{Data2: &mtproto.User_Data{
			Id:            do.Id,
			Self:          isSelf,
			Contact:       contact,
			MutualContact: mutualContact,
			AccessHash:    do.AccessHash,
			FirstName:     do.FirstName,
			LastName:      do.LastName,
			Username:      m.usernameCallback.GetAccountUsername(do.Id),
			Phone:         phone,
			Photo:         photo,
			Status:        status,
		}}}

		return data
	}
}

func (m *UserModel) GetUserByPhoneNumber(selfId int32, phoneNumber string) *userData {
	do := m.dao.UsersDAO.SelectByPhoneNumber(phoneNumber)
	if do == nil {
		return nil
	}
	do.Phone = phoneNumber
	return m.makeUserDataByDO(selfId, do)
}

func (m *UserModel) GetUserListByPhoneNumberList(selfId int32, phoneNumberList []string) []*userData {
	do := m.dao.UsersDAO.SelectByPhoneNumber(phoneNumberList[0])
	if do == nil {
		return nil
	}
	do.Phone = phoneNumberList[0]
	return []*userData{m.makeUserDataByDO(selfId, do)}
}

func (m *UserModel) GetUserByUsername(selfId int32, username string) *userData {
	do := m.dao.UsersDAO.SelectByUsername(username)
	if do == nil {
		return nil
	}
	return m.makeUserDataByDO(selfId, do)
}

func (m *UserModel) GetMyUserByPhoneNumber(phoneNumber string) *userData {
	do := m.dao.UsersDAO.SelectByPhoneNumber(phoneNumber)
	if do == nil {
		return nil
	}
	do.Phone = phoneNumber
	return m.makeUserDataByDO(do.Id, do)
}

func (m *UserModel) GetUserById(selfId int32, userId int32) *userData {
	do := m.dao.UsersDAO.SelectById(userId)
	return m.makeUserDataByDO(selfId, do)
}

func (m *UserModel) CreateNewUser(phoneNumber, countryCode, firstName, lastName string) *mtproto.TLUser {
	do := &dataobject.UsersDO{
		AccessHash:  rand.Int63(),
		Phone:       phoneNumber,
		FirstName:   firstName,
		LastName:    lastName,
		CountryCode: countryCode,
	}
	do.Id = int32(m.dao.UsersDAO.Insert(do))
	user := &mtproto.TLUser{Data2: &mtproto.User_Data{
		Id:            do.Id,
		Self:          true,
		Contact:       true,
		MutualContact: true,
		AccessHash:    do.AccessHash,
		FirstName:     do.FirstName,
		LastName:      do.LastName,
		Username:      do.Username,
		Phone:         phoneNumber,
		// TODO(@benqi): Load from db
		Photo:  mtproto.NewTLUserProfilePhotoEmpty().To_UserProfilePhoto(),
		Status: makeUserStatusOnline(),
	}}
	return user
}

func (m *UserModel) CreateNewUserPassword(userId int32) {
	// gen server_nonce
	do := &dataobject.UserPasswordsDO{
		UserId:     userId,
		ServerSalt: hex.EncodeToString(crypto.GenerateNonce(8)),
	}
	m.dao.UserPasswordsDAO.Insert(do)
}

func (m *UserModel) CheckAccessHashByUserId(userId int32, accessHash int64) bool {
	params := map[string]interface{}{
		"id":          userId,
		"access_hash": accessHash,
	}
	return m.dao.CommonDAO.CheckExists("users", params)
}

func (m *UserModel) GetCountryCodeByUser(userId int32) string {
	do := m.dao.UsersDAO.SelectCountryCode(userId)
	if do == nil {
		return ""
	} else {
		return do.CountryCode
	}
}
