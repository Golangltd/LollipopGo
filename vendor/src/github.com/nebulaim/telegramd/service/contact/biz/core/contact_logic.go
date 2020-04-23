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

package core

import (
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/service/contact/biz/dal/dataobject"
	"github.com/nebulaim/telegramd/service/contact/proto"
	"time"
)

// exclude
type contactData *dataobject.UserContactsDO

type contactLogic struct {
	selfUserID int32
	dao        *contactsDAO
}

// include deleted
func (c *contactLogic) getContactList() []contactData {
	doList := c.dao.UserContactsDAO.SelectAllUserContacts(c.selfUserID)
	contactList := make([]contactData, 0, len(doList))
	for i := 0; i < len(doList); i++ {
		contactList = append(contactList, &doList[i])
	}
	return contactList

}

// exclude deleted
func (c *contactLogic) getNotDeletedContactList() []contactData {
	doList := c.dao.UserContactsDAO.SelectUserContacts(c.selfUserID)
	contactList := make([]contactData, 0, len(doList))
	for i := 0; i < len(doList); i++ {
		contactList = append(contactList, &doList[i])
	}
	return contactList
}

func (c *contactLogic) GetContactByID(id int32) contactData {
	return c.dao.UserContactsDAO.SelectUserContact(c.selfUserID, id)
}

func (c *contactLogic) GetContactListByIDList(idList []int32) []contactData {
	contactList := make([]contactData, 0, len(idList))

	// TODO(@benqi): add SelectUserContactList in user_contacts_dao.go
	for _, id := range idList {
		do := c.dao.UserContactsDAO.SelectUserContact(c.selfUserID, id)
		if do != nil {
			contactList = append(contactList, do)
		}
	}
	return contactList
}

//func findContaceByPhone(contacts []contactData, phone string) *dataobject.UserContactsDO {
//	for _, c := range contacts {
//		if c.ContactPhone == phone {
//			return c
//		}
//	}
//	return nil
//}
//
//// include deleted
//func (c contactLogic) GetAllContactList() []contactData {
//	doList := dao.GetUserContactsDAO(dao.DB_SLAVE).SelectAllUserContacts(c.selfUserID)
//	contactList := make([]contactData, 0, len(doList))
//	for index, _ := range doList {
//		contactList = append(contactList, &doList[index])
//	}
//	return contactList
//}
//
//// exclude deleted
//func (c contactLogic) GetContactList() []contactData {
//	doList := dao.GetUserContactsDAO(dao.DB_SLAVE).SelectUserContacts(c.selfUserID)
//	contactList := make([]contactData, 0, len(doList))
//	for index, _ := range doList {
//		contactList = append(contactList, &doList[index])
//	}
//	return contactList
//}
//

func (c *contactLogic) importContact(contact *contact.InputContactData) (imported *contact.ImportedContactData) {
	var (
		mutualUpdated = false
		importers     int32
	)

	// TODO(@benqi): phone is me???
	// ## A加B，检查AB和BA
	// - B不是A的联系人，定义为B未加过A或B已经删了A，此时A只能添加B为单向联系人
	// 	- AB不存在，A添加B
	// 	- AB存在，不管以前是否被删，更新first_name和last_name并设置is_deleted = 0
	// - B是A的联系人，即BA存在并且B未删A
	//  - AB不存在或AB已经被删，A添加B，设置AB的mutual，设置BA的mutual
	//	- AB存在，更新first_name和last_name并设置is_deleted = 0
	//		-
	// 		- 则AB和BA肯定都是mutual，仅仅需要

	selfDO := c.dao.UserContactsDAO.SelectUserContact(c.selfUserID, contact.UserId)
	// contact -> self
	contactDO := c.dao.UserContactsDAO.SelectUserContact(contact.UserId, c.selfUserID)
	now := int32(time.Now().Unix())
	if contactDO == nil || contactDO.IsDeleted == 1 {
		// self -> contact
		if selfDO == nil {
			// input不是我的联系人
			do := &dataobject.UserContactsDO{
				OwnerUserId:      c.selfUserID,
				ContactUserId:    contact.UserId,
				ContactPhone:     contact.Phone,
				ContactFirstName: contact.FirstName,
				ContactLastName:  contact.LastName,
				Mutual:           0,
				Date2:            now,
			}
			do.Id = c.dao.UserContactsDAO.Insert(do)
		} else {
			c.dao.UserContactsDAO.UpdateContactNameByID(contact.FirstName, contact.LastName, selfDO.Id)
		}
	} else {
		// 我不是input的联系人
		if selfDO == nil {
			// input不是我的联系人
			do := &dataobject.UserContactsDO{
				OwnerUserId:      c.selfUserID,
				ContactUserId:    contact.UserId,
				ContactPhone:     contact.Phone,
				ContactFirstName: contact.FirstName,
				ContactLastName:  contact.LastName,
				Mutual:           1,
				Date2:            now,
			}
			do.Id = c.dao.UserContactsDAO.Insert(do)
			c.dao.UserContactsDAO.UpdateMutualByID(1, contactDO.Id)
			mutualUpdated = true
		} else {
			// selfDeleted := selfDO.IsDeleted
			c.dao.UserContactsDAO.UpdateContactNameByID(contact.FirstName, contact.LastName, selfDO.Id)

			// update不影响selfDO里的值，直接拿来用
			if selfDO.IsDeleted == 1 {
				c.dao.UserContactsDAO.UpdateMutualByID(1, selfDO.Id)
				c.dao.UserContactsDAO.UpdateMutualByID(1, contactDO.Id)
				mutualUpdated = true
			}
		}
	}

	// increase importers
	importers = c.increasePopularContact(contact.Phone)

	return &contact.ImportedContactData{
		UserId:        contact.UserId,
		Importers:     importers,
		MutualUpdated: mutualUpdated,
	}
}

func (c *contactLogic) importContacts(contacts []*contact.InputContactData) []*contact.ImportedContactData {
	// TODO(@benqi): 优化, 减少数据库操作：
	// ## importContacts优化方法
	// - user_contacts表添加relation_id字段
	// - 由selfUserId和contact.userId生成relation_id，由relation_id_list查出所有contacts
	// - 处理完后批量更新到数据库
	// - 处理importers
	//  - 方法1，存储在mysql里：
	// 		- 如果AB或BA存在，则importer肯定存在，只需要批量自增importers
	//  	- AB或BA都不存在，批量查出importers，计算最终的importers存储到db里（插入或自增）
	//  - 方法2，存储在redis里：
	//		- 批量执行redis的increase
	//
	importeds := make([]*contact.ImportedContactData, 0, len(contacts))
	for _, contact := range contacts {
		importeds = append(importeds, c.importContact(contact))
	}

	return importeds
}

func (c *contactLogic) deleteContact(deleteId int32) *contact.DeleteResult {
	// A 删除 B
	// 如果AB is mutual，则BA设置为非mutual

	//var needUpdate = false
	//
	//c.dao.UserContactsDAO.DeleteContacts(c.selfUserID, []int32{deleteId})
	//if deleteId != c.selfUserID && mutual {
	//	c.dao.UserContactsDAO.UpdateMutual(0, deleteId, c.selfUserID)
	//	needUpdate = true
	//}

	return nil
}

func (c *contactLogic) deleteContacts(deleteId []int32) []*contact.DeleteResult {
	// A 删除 B
	// 如果AB is mutual，则BA设置为非mutual

	//var needUpdate = false
	//
	//c.dao.UserContactsDAO.DeleteContacts(c.selfUserID, []int32{deleteId})
	//if deleteId != c.selfUserID && mutual {
	//	c.dao.UserContactsDAO.UpdateMutual(0, deleteId, c.selfUserID)
	//	needUpdate = true
	//}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////
func (c *contactLogic) blockUser(blockId int32) bool {
	c.dao.UserContactsDAO.UpdateBlock(1, c.selfUserID, blockId)
	return true
}

func (c *contactLogic) unBlockUser(blockedId int32) bool {
	c.dao.UserContactsDAO.UpdateBlock(0, c.selfUserID, blockedId)
	return true
}

func (c *contactLogic) getBlockedList(offset, limit int32) []*mtproto.ContactBlocked {
	// TODO(@benqi): enable offset
	doList := c.dao.UserContactsDAO.SelectBlockedList(c.selfUserID, limit)
	bockedList := make([]*mtproto.ContactBlocked, 0, len(doList))
	for i := 0; i < len(doList); i++ {
		blocked := &mtproto.ContactBlocked{
			Constructor: mtproto.TLConstructor_CRC32_contactBlocked,
			Data2: &mtproto.ContactBlocked_Data{
				UserId: doList[i].ContactUserId,
				Date:   doList[i].Date2,
			},
		}
		bockedList = append(bockedList, blocked)
	}

	return bockedList
}

func (c *contactLogic) getPopularContacts(phones []string) []int32 {
	doList := c.dao.PopularContactsDAO.SelectImportersList(phones)
	_ = doList
	return nil
}

func (c *contactLogic) increasePopularContact(phone string) int32 {
	// TODO(@benqi): storage to redis
	do := c.dao.PopularContactsDAO.SelectImporters(phone)
	if do == nil {
		// importers = 1
		do = &dataobject.PopularContactsDO{
			Phone:     phone,
			Importers: 1,
		}
		c.dao.PopularContactsDAO.Insert(do)
	} else {
		do.Importers += 1
		c.dao.PopularContactsDAO.IncreaseImporters(phone)
	}
	return do.Importers
}

func (c *contactLogic) checkContactAndMutual(contactId int32) (bool, bool) {
	do := c.dao.UserContactsDAO.SelectUserContact(c.selfUserID, contactId)
	if do == nil {
		return false, false
	} else {
		if do.IsDeleted == 1 {
			return false, false
		} else {
			return true, do.Mutual == 1
		}
	}
}
