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

package message

/*
import (
	"github.com/nebulaim/telegramd/biz/base"
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"time"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/biz/core"
)

type MessageBox struct {
	UserId          int32
	MessageId       int32
	DialogMessageId int64
	RandomId        int64
	Message         *mtproto.Message
	dao             *messagesDAO
}

type MessageBoxList []*MessageBox

// type OnOutboxCreated
type OnOutboxCreated func(int32)
type OnInboxSendOK func(int32, int32)

// 新增
//func (m *MessageModel) CreateMessageOutboxByNew(fromId int32, peer *base.PeerUtil, clientRandomId int64, message2 *mtproto.Message, cb OnOutboxCreated) (box *MessageBox) {
//	now := int32(time.Now().Unix())
//	messageDO := &dataobject.MessagesDO{
//		UserId:           fromId,
//		UserMessageBoxId: int32(core.NextMessageBoxId(fromId)),
//		DialogMessageId:  core.GetUUID(),
//		SenderUserId:     fromId,
//		MessageBoxType:   MESSAGE_BOX_TYPE_OUTGOING,
//		PeerType:         int8(peer.PeerType),
//		PeerId:           peer.PeerId,
//		RandomId:         clientRandomId,
//		Date2:            now,
//		Deleted:          0,
//	}
//
//	switch message2.GetConstructor() {
//	case mtproto.TLConstructor_CRC32_messageEmpty:
//		messageDO.MessageType = MESSAGE_TYPE_MESSAGE_EMPTY
//	case mtproto.TLConstructor_CRC32_message:
//		messageDO.MessageType = MESSAGE_TYPE_MESSAGE
//		message := message2.To_Message()
//
//		// mentioned = message.GetMentioned()
//		message.SetId(messageDO.UserMessageBoxId)
//	case mtproto.TLConstructor_CRC32_messageService:
//		messageDO.MessageType = MESSAGE_TYPE_MESSAGE_SERVICE
//		message := message2.To_MessageService()
//
//		// mentioned = message.GetMentioned()
//		message.SetId(messageDO.UserMessageBoxId)
//	}
//
//	messageData, _ := json.Marshal(message2)
//	messageDO.MessageData = string(messageData)
//
//	// TODO(@benqi): pocess clientRandomId dup
//	m.dao.MessagesDAO.Insert(messageDO)
//
//	box = &MessageBox{
//		UserId:          fromId,
//		MessageId:       messageDO.UserMessageBoxId,
//		DialogMessageId: messageDO.DialogMessageId,
//		RandomId:        clientRandomId,
//		Message:         message2,
//		dao:             m.dao,
//	}
//
//	if cb != nil {
//		cb(messageDO.UserMessageBoxId)
//	}
//	return
//}
//

func (m *MessageModel) CreateMessageOutboxByNew(fromId int32, peerType, peerId int32,clientRandomId, did, dialogMessageId int64, message2 *mtproto.Message, cb OnOutboxCreated) (box *MessageBox) {
	boxDO := &dataobject.MessageBoxesDO{
		UserId:           fromId,
		UserMessageBoxId: int32(core.NextMessageBoxId(fromId)),
		DialogId:         did,
		DialogMessageId:  dialogMessageId,
		SenderUserId:     fromId,
		MessageBoxType:   MESSAGE_BOX_TYPE_OUTGOING,
		PeerType:         int8(peerType),
		PeerId:           peerId,
		Date2:            int32(time.Now().Unix()),
		Deleted:          0,
	}

	m.dao.MessageBoxesDAO.Insert(boxDO)

	box = &MessageBox{
		UserId:          fromId,
		MessageId:       boxDO.UserMessageBoxId,
		DialogMessageId: dialogMessageId,
		RandomId:        clientRandomId,
		Message:         message2,
		dao:             m.dao,
	}

	if cb != nil {
		cb(boxDO.UserMessageBoxId)
	}
	return
}

////func (m *MessageModel) MakeMessageBoxByLoad(userId int32, peer *base.PeerUtil, messageId int32) (box *MessageBox) {
////	return nil
////}
//
//func (this *MessageBox) InsertMessageToInbox(fromId int32, peer *base.PeerUtil, cb OnInboxSendOK) (MessageBoxList, error) {
//	switch peer.PeerType {
//	case base.PEER_USER:
//		return this.insertUserMessageToInbox(fromId, peer, cb)
//	case base.PEER_CHAT:
//		return this.insertChatMessageToInbox(fromId, peer, cb)
//	// case base.PEER_CHANNEL:
//	// 	return this.insertChannelMessageToInbox(fromId, peer, cb)
//	default:
//		//	panic("invalid peer")
//		return nil, fmt.Errorf("invalid peer")
//	}
//}

func (this *MessageBox) getPeerMessageId(userId, messageId, peerId int32) int32 {
	do := this.dao.MessagesDAO.SelectPeerMessageId(peerId, userId, messageId)
	if do == nil {
		return 0
	} else {
		return do.UserMessageBoxId
	}
}

//func (this *MessageBox) makeInboxMessageDO(fromId int32, peer *base.PeerUtil, inboxUserId int32) *MessageBox {
//	now := int32(time.Now().Unix())
//	messageDO := &dataobject.MessagesDO{
//		UserId:           inboxUserId,
//		UserMessageBoxId: int32(core.NextMessageBoxId(inboxUserId)),
//		DialogMessageId:  this.DialogMessageId,
//		SenderUserId:     this.UserId,
//		MessageBoxType:   MESSAGE_BOX_TYPE_INCOMING,
//		PeerType:         int8(peer.PeerType),
//		PeerId:           peer.PeerId,
//		RandomId:         this.RandomId,
//		Date2:            now,
//		Deleted:          0,
//	}
//
//	inboxMessage := proto.Clone(this.Message).(*mtproto.Message)
//	// var mentioned = false
//
//	switch this.Message.GetConstructor() {
//	case mtproto.TLConstructor_CRC32_messageEmpty:
//		messageDO.MessageType = MESSAGE_TYPE_MESSAGE_EMPTY
//	case mtproto.TLConstructor_CRC32_message:
//		messageDO.MessageType = MESSAGE_TYPE_MESSAGE
//
//		m2 := inboxMessage.To_Message()
//		m2.SetOut(false)
//		if m2.GetReplyToMsgId() != 0 {
//			replyMsgId := this.getPeerMessageId(fromId, m2.GetReplyToMsgId(), inboxUserId)
//			m2.SetReplyToMsgId(replyMsgId)
//		}
//		m2.SetId(messageDO.UserMessageBoxId)
//		// mentioned = m2.GetMentioned()
//	case mtproto.TLConstructor_CRC32_messageService:
//		messageDO.MessageType = MESSAGE_TYPE_MESSAGE_SERVICE
//
//		m2 := inboxMessage.To_MessageService()
//		m2.SetOut(false)
//		m2.SetId(messageDO.UserMessageBoxId)
//	}
//
//	messageData, _ := json.Marshal(inboxMessage)
//	messageDO.MessageData = string(messageData)
//
//	// TODO(@benqi): rpocess clientRandomId dup
//	this.dao.MessagesDAO.Insert(messageDO)
//
//	return &MessageBox{
//		UserId:          inboxUserId,
//		MessageId:       messageDO.UserMessageBoxId,
//		DialogMessageId: messageDO.DialogMessageId,
//		RandomId:        this.RandomId,
//		Message:         inboxMessage,
//	}
//}

func (this *MessageBox) makeInboxMessageDO(fromId int32, peerType int, peerId int32, inboxUserId int32) *MessageBox {
	now := int32(time.Now().Unix())
	messageBoxDO := &dataobject.MessageBoxesDO{
		UserId:           inboxUserId,
		UserMessageBoxId: int32(core.NextMessageBoxId(inboxUserId)),
		DialogMessageId:  this.DialogMessageId,
		SenderUserId:     this.UserId,
		MessageBoxType:   MESSAGE_BOX_TYPE_INCOMING,
		PeerType:         int8(peerType),
		PeerId:           peerId,
		RandomId:         this.RandomId,
		Date2:            now,
		Deleted:          0,
	}

	inboxMessage := proto.Clone(this.Message).(*mtproto.Message)
	//// var mentioned = false
	//
	//switch this.Message.GetConstructor() {
	//case mtproto.TLConstructor_CRC32_messageEmpty:
	//	messageDO.MessageType = MESSAGE_TYPE_MESSAGE_EMPTY
	//case mtproto.TLConstructor_CRC32_message:
	//	messageDO.MessageType = MESSAGE_TYPE_MESSAGE
	//
	//	m2 := inboxMessage.To_Message()
	//	m2.SetOut(false)
	//	if m2.GetReplyToMsgId() != 0 {
	//		replyMsgId := this.getPeerMessageId(fromId, m2.GetReplyToMsgId(), inboxUserId)
	//		m2.SetReplyToMsgId(replyMsgId)
	//	}
	//	m2.SetId(messageDO.UserMessageBoxId)
	//	// mentioned = m2.GetMentioned()
	//case mtproto.TLConstructor_CRC32_messageService:
	//	messageDO.MessageType = MESSAGE_TYPE_MESSAGE_SERVICE
	//
	//	m2 := inboxMessage.To_MessageService()
	//	m2.SetOut(false)
	//	m2.SetId(messageDO.UserMessageBoxId)
	//}
	//
	//messageData, _ := json.Marshal(inboxMessage)
	//messageDO.MessageData = string(messageData)

	// TODO(@benqi): rpocess clientRandomId dup
	this.dao.MessageBoxesDAO.Insert(messageBoxDO)

	return &MessageBox{
		UserId:          inboxUserId,
		MessageId:       messageBoxDO.UserMessageBoxId,
		DialogMessageId: messageBoxDO.DialogMessageId,
		RandomId:        this.RandomId,
		Message:         inboxMessage,
	}
}

//// 发送到收件箱
//func (this *MessageBox) insertUserMessageToInbox(fromId, peerId int32, cb OnInboxSendOK) (MessageBoxList, error) {
//	inbox := this.makeInboxMessageDO2(fromId, peer, peer.PeerId)
//	if cb != nil {
//		cb(inbox.UserId, inbox.MessageId)
//	}
//	return []*MessageBox{inbox}, nil
//}

// 发送到收件箱
func (this *MessageBox) InsertUserMessageToInbox(fromId, peerId int32, cb OnInboxSendOK) (*MessageBox, error) {
	inbox := this.makeInboxMessageDO(fromId, int(base.PEER_USER), peerId, peerId)
	if cb != nil {
		cb(inbox.UserId, inbox.MessageId)
	}
	return inbox, nil
}

//// 发送chat message到收件箱
//func (this *MessageBox) insertChatMessageToInbox(fromId int32, peer *base.PeerUtil, cb OnInboxSendOK) (MessageBoxList, error) {
//	doList := this.dao.ChatParticipantsDAO.SelectByChatId(peer.PeerId)
//
//	var inoutBoxList MessageBoxList = make([]*MessageBox, 0, len(doList))
//	for _, do := range doList {
//		if do.UserId == this.UserId {
//			continue
//		}
//		inbox := this.makeInboxMessageDO(fromId, peer, do.UserId)
//		glog.Info("insertChatMessageToInbox - ", inbox)
//		if cb != nil {
//			cb(inbox.UserId, inbox.MessageId)
//		}
//		inoutBoxList = append(inoutBoxList, inbox)
//	}
//
//	return inoutBoxList, nil
//}

// 发送chat message到收件箱
func (this *MessageBox) InsertChatMessageToInbox(fromId, peerId int32, cb OnInboxSendOK) (MessageBoxList, error) {
	doList := this.dao.ChatParticipantsDAO.SelectByChatId(peerId)

	var inoutBoxList MessageBoxList = make([]*MessageBox, 0, len(doList))
	for _, do := range doList {
		if do.UserId == this.UserId {
			continue
		}
		inbox := this.makeInboxMessageDO(fromId, int(base.PEER_CHAT), peerId, do.UserId)
		glog.Info("insertChatMessageToInbox - ", inbox)
		if cb != nil {
			cb(inbox.UserId, inbox.MessageId)
		}
		inoutBoxList = append(inoutBoxList, inbox)
	}

	return inoutBoxList, nil
}

// 发送channel message到收件箱
//func (this *MessageBox) insertChannelMessageToInbox(fromId int32, peer *base.PeerUtil, cb OnInboxSendOK) (MessageBoxList, error) {
//	switch this.Message.GetConstructor() {
//	case mtproto.TLConstructor_CRC32_message:
//	case mtproto.TLConstructor_CRC32_messageService:
//	default:
//		panic("invalid messageEmpty type")
//		// return
//	}
//	return []*MessageBox{}, nil
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (this *MessageBoxList) ToMessageList() []*mtproto.Message {
	messageList := make([]*mtproto.Message, 0, len(*this))
	for _, box := range messageList {
		messageList = append(messageList, box)
	}
	return messageList
}
*/
