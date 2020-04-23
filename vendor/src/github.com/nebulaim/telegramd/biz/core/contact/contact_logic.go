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

package contact

import (
	"time"

	"github.com/nebulaim/telegramd/biz/dal/dataobject"
	"github.com/nebulaim/telegramd/proto/mtproto"
	// "github.com/mattermost/mattermost-server/model"
)

//type contactUser struct {
//	userId int32
//	phone string
//	firstName string
//	lastName string
//}

// exclude

type contactData *dataobject.UserContactsDO
type contactLogic struct {
	selfUserId int32
	dao        *contactsDAO
}

func (m *ContactModel) MakeContactLogic(userId int32) *contactLogic {
	return &contactLogic{
		selfUserId: userId,
		dao:        m.dao,
	}
}

func findContaceByPhone(contacts []contactData, phone string) *dataobject.UserContactsDO {
	for _, c := range contacts {
		if c.ContactPhone == phone {
			return c
		}
	}
	return nil
}

// include deleted
func (c contactLogic) GetAllContactList() []contactData {
	doList := c.dao.UserContactsDAO.SelectAllUserContacts(c.selfUserId)
	contactList := make([]contactData, 0, len(doList))
	for index, _ := range doList {
		contactList = append(contactList, &doList[index])
	}
	return contactList
}

// exclude deleted
func (c contactLogic) GetContactList() []contactData {
	doList := c.dao.UserContactsDAO.SelectUserContacts(c.selfUserId)
	contactList := make([]contactData, 0, len(doList))
	for index, _ := range doList {
		contactList = append(contactList, &doList[index])
	}
	return contactList
}

func (c contactLogic) ImportContact(userId int32, phone, firstName, lastName string) bool {
	var needUpdate bool = false

	// TODO(@benqi): phone is me???

	// 我->input
	byMy := c.dao.UserContactsDAO.SelectUserContact(c.selfUserId, userId)
	// input->我
	byInput := c.dao.UserContactsDAO.SelectUserContact(userId, c.selfUserId)

	now := int32(time.Now().Unix())
	if byInput == nil {
		// 我不是input的联系人
		if byMy == nil {
			// input不是我的联系人
			do := &dataobject.UserContactsDO{
				OwnerUserId:      c.selfUserId,
				ContactUserId:    userId,
				ContactPhone:     phone,
				ContactFirstName: firstName,
				ContactLastName:  lastName,
				Mutual:           0,
				Date2:            now,
			}
			do.Id = c.dao.UserContactsDAO.Insert(do)
		} else {
			c.dao.UserContactsDAO.UpdateContactNameById(firstName, lastName, byMy.Id)
		}
	} else {
		// 我不是input的联系人
		if byMy == nil {
			// input不是我的联系人
			do := &dataobject.UserContactsDO{
				OwnerUserId:      c.selfUserId,
				ContactUserId:    userId,
				ContactPhone:     phone,
				ContactFirstName: firstName,
				ContactLastName:  lastName,
				Mutual:           1,
				Date2:            now,
			}
			do.Id = c.dao.UserContactsDAO.Insert(do)
			c.dao.UserContactsDAO.UpdateMutual(1, userId, c.selfUserId)
			needUpdate = true
		} else {
			c.dao.UserContactsDAO.UpdateContactNameById(firstName, lastName, byMy.Id)
			if byMy.IsDeleted == 1 {
				c.dao.UserContactsDAO.UpdateMutual(1, userId, c.selfUserId)
				c.dao.UserContactsDAO.UpdateMutual(1, c.selfUserId, userId)
				needUpdate = true
			}
		}
	}

	return needUpdate
}

func (c contactLogic) DeleteContact(deleteId int32, mutual bool) bool {
	// A 删除 B
	// 如果AB is mutual，则BA设置为非mutual

	var needUpdate = false

	c.dao.UserContactsDAO.DeleteContacts(c.selfUserId, []int32{deleteId})

	if deleteId != c.selfUserId && mutual {
		c.dao.UserContactsDAO.UpdateMutual(0, deleteId, c.selfUserId)
		needUpdate = true
	}

	return needUpdate
}

//// imported int64, 低32位为InputContact的index， 高32位为userId
//func (c contactLogic) AddContactList(contactList []*mtproto.InputContact) (importedList []int64, retryList []int64) {
//	contacts := c.GetAllContactList()
//	// c.dao.UserContactsDAO.SelectUserContacts(c.selfUserId)
//	for i, v := range contactList {
//		inputContact := v.To_InputPhoneContact()
//		found := findContaceByPhone(contacts, inputContact.GetPhone())
//
//		// TODO(@benqi): ?? popularContact#5ce14175 client_id:long importers:int = PopularContact;
//		if found == nil {
//			// Not found, insert.
//			// Check user exist by phone number
//			// TODO(@benqi): mutual
//			do := &dataobject.UserContactsDO{
//				OwnerUserId:      c.selfUserId,
//				ContactPhone:     inputContact.GetPhone(),
//				ContactFirstName: inputContact.GetFirstName(),
//				ContactLastName:  inputContact.GetLastName(),
//			}
//
//			do.Id = int32(c.dao.UserContactsDAO.Insert(do))
//
//			// 低32位为InputContact的index， 高32位为userId
//			importedList = append(importedList, int64(i) | int64(do.Id) << 32)
//		} else {
//			// delete
//			if found.IsDeleted == 1 {
//				// 如果已经删除，则将delete设置为0
//				// update delete = 0
//				c.dao.UserContactsDAO.UpdateContactNameById(inputContact.GetFirstName(), inputContact.GetLastName(), found.Id)
//				importedList = append(importedList, int64(i) | int64(found.Id) << 32)
//			} else {
//				if found.ContactFirstName != inputContact.GetFirstName() || found.ContactLastName != inputContact.GetLastName() {
//					// 修改联系人名字
//					c.dao.UserContactsDAO.UpdateContactNameById(inputContact.GetFirstName(), inputContact.GetLastName(), found.Id)
//					importedList = append(importedList, int64(i) | int64(found.Id) << 32)
//				} else {
//					retryList = append(retryList, inputContact.GetClientId())
//				}
//			}
//		}
//	}
//
//	return
//}

/////////////////////////////////////////////////////////////////////////////////////////
func (c contactLogic) BlockUser(blockId int32) bool {
	c.dao.UserContactsDAO.UpdateBlock(1, c.selfUserId, blockId)
	return true
}

func (c contactLogic) UnBlockUser(blockedId int32) bool {
	c.dao.UserContactsDAO.UpdateBlock(0, c.selfUserId, blockedId)
	return true
}

func (c contactLogic) GetBlockedList(offset, limit int32) []*mtproto.ContactBlocked {
	// TODO(@benqi): enable offset
	doList := c.dao.UserContactsDAO.SelectBlockedList(c.selfUserId, limit)
	bockedList := make([]*mtproto.ContactBlocked, 0, len(doList))
	for _, do := range doList {
		blocked := &mtproto.ContactBlocked{
			Constructor: mtproto.TLConstructor_CRC32_contactBlocked,
			Data2: &mtproto.ContactBlocked_Data{
				UserId: do.ContactUserId,
				Date:   do.Date2,
			},
		}
		bockedList = append(bockedList, blocked)
	}
	return bockedList
}

func (c contactLogic) SearchContacts(q string, limit int32) []int32 {
	contactList := c.GetContactList()
	idList := make([]int32, 0, len(contactList)+1)
	idList = append(idList, c.selfUserId)
	for _, c2 := range contactList {
		idList = append(idList, c2.ContactUserId)
	}

	// TODO(@benqi): 区分大小写

	// 构造模糊查询字符串
	q = "%" + q + "%"
	doList := c.dao.UsersDAO.SearchByQueryNotIdList(q, idList, limit)
	founds := make([]int32, 0, len(doList))
	for _, do := range doList {
		founds = append(founds, do.Id)
	}
	return founds
}

//const (
//	_unregisted = iota
//	_noneContact
//	_contact
//	_mutualContact
//)

type contactItem struct {
	c               *mtproto.InputContact_Data
	unregisted      bool
	userId          int32
	contactId       int32
	importContactId int32
}

func (c *contactLogic) importContacts(contacts []*mtproto.InputContact_Data) []*mtproto.ImportedContact_Data {
	//// TODO(@benqi): 优化, 减少数据库操作：
	//// ## importContacts优化方法
	//// - user_contacts表添加relation_id字段
	//// - 由selfUserId和contact.userId生成relation_id，由relation_id_list查出所有contacts
	//// - 处理完后批量更新到数据库
	//// - 处理importers
	////  - 方法1，存储在mysql里：
	//// 		- 如果AB或BA存在，则importer肯定存在，只需要批量自增importers
	////  	- AB或BA都不存在，批量查出importers，计算最终的importers存储到db里（插入或自增）
	////  - 方法2，存储在redis里：
	////		- 批量执行redis的increase
	////
	//importeds := make([]*contact.ImportedContactData, 0, len(contacts))
	//for _, contact := range contacts {
	//	importeds = append(importeds, c.importContact(contact))
	//}
	//
	//return importeds

	importContacts := make(map[string]*contactItem)
	// 1. 找出phoneList里的所有已注册用户
	phoneList := make([]string, 0, len(contacts))
	for _, c2 := range contacts {
		phoneList = append(phoneList, c2.Phone)
		importContacts[c2.Phone] = &contactItem{unregisted: true, c: c2}
	}

	// 已注册
	registeredContacts := c.dao.UsersDAO.SelectUsersByPhoneList(phoneList)
	for i := 0; i < len(registeredContacts); i++ {
		if c2, ok := importContacts[registeredContacts[i].Phone]; ok {
			c2.unregisted = false
			c2.userId = registeredContacts[i].Id
		}
	}

	myContacts := c.dao.UserContactsDAO.SelectAllUserContacts(c.selfUserId)
	for i := 0; i < len(myContacts); i++ {
		if c2, ok := importContacts[myContacts[i].ContactPhone]; ok {
			c2.contactId = myContacts[i].ContactUserId
		}
	}

	importedMyContacts := c.dao.ImportedContactsDAO.SelectImportedContacts(c.selfUserId)
	for i := 0; i < len(importedMyContacts); i++ {
		for _, c2 := range importContacts {
			if c2.userId == importedMyContacts[i].ImportedUserId {
				c2.importContactId = c2.userId
				break
			}
		}
	}

	for _, c2 := range importContacts {
		if c2.unregisted {
			do := &dataobject.UnregisteredContactsDO{
				ImporterUserId:   c.selfUserId,
				ContactPhone:     c2.c.Phone,
				ContactFirstName: c2.c.FirstName,
				ContactLastName:  c2.c.LastName,
			}
			c.dao.UnregisteredContactsDAO.Insert(do)
		} else {
			if c2.contactId != 0 || c2.contactId == 0 && c2.importContactId == 0 {
			} else {
				//if c2.importContactId == 0 {
				//
				//} else {
				//
				//}
			}
		}
		//if c2.userId == importedMyContacts[i].ImportedUserId == {
		//	c2.mutualContact = true
		//	break
		//}
	}

	//// select * from users where phone in (:phoneList)
	//// users := GetUserListByPhoneNumberList[]*mtproto.User{}// model.GetGroupNameFromUserIds()
	//for _, c2 := range contacts {
	//	if checkExists(c2.Phone) {
	//
	//	}
	//	// c.dao.SelectByPhoneNumber()
	//}

	return nil
}
