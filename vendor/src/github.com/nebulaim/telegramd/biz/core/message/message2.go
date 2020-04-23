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

package message

import (
	"github.com/nebulaim/telegramd/biz/base"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/golang/glog"
	base2 "github.com/nebulaim/telegramd/baselib/base"
)

const (
	MESSAGE_TYPE_UNKNOWN         = 0
	MESSAGE_TYPE_MESSAGE_EMPTY   = 1
	MESSAGE_TYPE_MESSAGE         = 2
	MESSAGE_TYPE_MESSAGE_SERVICE = 3
)
const (
	MESSAGE_BOX_TYPE_INCOMING = 0
	MESSAGE_BOX_TYPE_OUTGOING = 1
	MESSAGE_BOX_TYPE_CHANNEL  = 2
)

const (
	PTS_UNKNOWN             = 0
	PTS_MESSAGE_OUTBOX      = 1
	PTS_MESSAGE_INBOX       = 2
	PTS_READ_HISTORY_OUTBOX = 3
	PTS_READ_HISTORY_INBOX  = 4
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// Loadhistory
func (m *MessageModel) LoadBackwardHistoryMessages(userId int32, peerType, peerId int32, offset int32, limit int32) (messages []*mtproto.Message) {
	// TODO(@benqi): chat and channel

	did := makeDialogId(userId, peerType, peerId)

	switch peerType {
	case base.PEER_USER, base.PEER_CHAT:
		boxDOList := m.dao.MessageBoxesDAO.SelectBackwardByOffsetLimit(userId, did, offset, limit)
		if len(boxDOList) == 0 {
			messages = []*mtproto.Message{}
			return
		}

		dialogMessageIdList := make([]int32, 0, len(boxDOList))
		for i := 0; i < len(boxDOList); i++ {
			dialogMessageIdList = append(dialogMessageIdList, boxDOList[i].DialogMessageId)
		}
		mDataDOList := m.dao.MessageDatasDAO.SelectMessageList(did, dialogMessageIdList)
		if len(mDataDOList) == 0 {
			messages = []*mtproto.Message{}

			// TODO(@benqi): logo
			return
		}

		for i := 0; i < len(boxDOList); i++ {
			for j := 0; j < len(mDataDOList); j++ {
				if boxDOList[i].DialogMessageId == mDataDOList[j].DialogMessageId {
					box := m.makeMessageBoxByDO(&boxDOList[i], &mDataDOList[j])
					messages = append(messages, box.ToMessage(userId))
					break
				}
			}
		}

	case base.PEER_CHANNEL:
		boxDOList := m.dao.ChannelMessagesDAO.SelectBackwardByOffsetLimit(peerId, offset, limit)
		for i := 0; i < len(boxDOList); i++ {
			box := m.makeChannelMessageBoxByDO(&boxDOList[i])
			messages = append(messages, box.ToMessage(userId))
		}
	default:
		// TODO(@benqi): log
	}
	return
}

func (m *MessageModel) LoadForwardHistoryMessages(userId int32, peerType, peerId int32, offset int32, limit int32) (messages []*mtproto.Message) {
	did := makeDialogId(userId, peerType, peerId)

	switch peerType {
	case base.PEER_USER, base.PEER_CHAT:
		boxDOList := m.dao.MessageBoxesDAO.SelectForwardByPeerOffsetLimit(userId, did, offset, limit)
		if len(boxDOList) == 0 {
			messages = []*mtproto.Message{}
			return
		}

		dialogMessageIdList := make([]int32, 0, len(boxDOList))
		for i := 0; i < len(boxDOList); i++ {
			dialogMessageIdList = append(dialogMessageIdList, boxDOList[i].DialogMessageId)
		}
		mDataDOList := m.dao.MessageDatasDAO.SelectMessageList(did, dialogMessageIdList)
		if len(mDataDOList) == 0 {
			messages = []*mtproto.Message{}

			// TODO(@benqi): log
			return
		}

		for i := 0; i < len(boxDOList); i++ {
			for j := 0; j < len(mDataDOList); j++ {
				if boxDOList[i].DialogMessageId == mDataDOList[j].DialogMessageId {
					box := m.makeMessageBoxByDO(&boxDOList[i], &mDataDOList[j])
					messages = append(messages, box.ToMessage(userId))
					break
				}
			}
		}

	case base.PEER_CHANNEL:
		boxDOList := m.dao.ChannelMessagesDAO.SelectForwardByOffsetLimit(peerId, offset, limit)
		for i := 0; i < len(boxDOList); i++ {
			box := m.makeChannelMessageBoxByDO(&boxDOList[i])
			messages = append(messages, box.ToMessage(userId))
		}
	default:
		// TODO(@benqi): log
	}
	return
}

func (m *MessageModel) GetUserMessagesByMessageIdList(userId int32, idList []int32) (messages []*mtproto.Message) {
	if len(idList) == 0 {
		messages = []*mtproto.Message{}
	} else {
		boxDOList := m.dao.MessageBoxesDAO.SelectByMessageIdList(userId, idList)
		glog.Info(boxDOList)
		if len(boxDOList) == 0 {
			messages = []*mtproto.Message{}
			return
		}

		messageDataIdList := make([]int64, 0, len(boxDOList))
		for i := 0; i < len(boxDOList); i++ {
			messageDataIdList = append(messageDataIdList, boxDOList[i].MessageDataId)
		}
		mDataDOList := m.dao.MessageDatasDAO.SelectMessageListByDataIdList(messageDataIdList)
		glog.Info(mDataDOList)
		if len(mDataDOList) == 0 {
			messages = []*mtproto.Message{}
			// TODO(@benqi): log
			return
		}

		for i := 0; i < len(boxDOList); i++ {
			for j := 0; j < len(mDataDOList); j++ {
				if boxDOList[i].DialogMessageId == mDataDOList[j].DialogMessageId {
					box := m.makeMessageBoxByDO(&boxDOList[i], &mDataDOList[j])
					messages = append(messages, box.ToMessage(userId))
					break
				}
			}
		}
	}
	return
}

func (m *MessageModel) GetPeerMessageListByMessageDataId(userId int32, messageDataId int64) (boxList []*MessageBox2) {
	doList := m.dao.MessageBoxesDAO.SelectPeerMessageList(userId, messageDataId)
	for _, do := range doList {
		// TODO(@benqi): check data
		box := &MessageBox2{
			OwnerId:        do.UserId,
			MessageId:      do.UserMessageBoxId,
			MessageBoxType: do.MessageBoxType,
			MediaUnread:    base2.Int8ToBool(do.MediaUnread),
		}
		boxList = append(boxList, box)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func (m *MessageModel) GetPeerMessageId(userId, messageId, peerId int32) int32 {
	do := m.dao.MessageBoxesDAO.SelectPeerMessageId(peerId, userId, messageId)
	if do == nil {
		return 0
	} else {
		return do.UserMessageBoxId
	}
}


func (m *MessageModel) DeleteByMessageIdList(userId int32, idList []int32) {
	if len(idList) > 0 {
		m.dao.MessageBoxesDAO.DeleteMessagesByMessageIdList(userId, idList)
	}
}

func (m *MessageModel) GetPeerDialogMessageIdList(userId int32, idList []int32) map[int32][]int32 {
	doList := m.dao.MessageBoxesDAO.SelectPeerDialogMessageIdList(userId, idList)
	peerMessageIdListMap := make(map[int32][]int32)

	for _, do := range doList {
		if messageIdList, ok := peerMessageIdListMap[do.UserId]; !ok {
			peerMessageIdListMap[do.UserId] = []int32{do.UserMessageBoxId}
		} else {
			peerMessageIdListMap[do.UserId] = append(messageIdList, do.UserMessageBoxId)
		}
	}

	return peerMessageIdListMap
}

func (m *MessageModel) GetMessageIdListByDialog(userId int32, peer *base.PeerUtil) []int32 {
	did := makeDialogId(userId, peer.PeerType, peer.PeerId)
	doList := m.dao.MessageBoxesDAO.SelectDialogMessageIdList(userId, did)
	idList := make([]int32, 0, len(doList))
	for i := 0; i < len(doList); i++ {
		idList = append(idList, doList[i].UserMessageBoxId)
	}
	return idList
}


func (m *MessageModel) GetClearHistoryMessages(userId int32, peer *base.PeerUtil) (lastMessageBox *MessageBox2, idList []int32) {
	idList = []int32{}
	did := makeDialogId(userId, peer.PeerType, peer.PeerId)
	doList := m.dao.MessageBoxesDAO.SelectDialogMessageIdList(userId, did)
	for i := 0; i < len(doList); i++ {
		if i == 0 {
			var err error
			lastMessageBox, err  = m.GetMessageBox2(int32(base.PEER_USER), userId, doList[0].UserMessageBoxId)
			if err != nil {
				return
			}
		} else {
			idList = append(idList, doList[i].UserMessageBoxId)
		}
	}
	return
}

func (m *MessageModel) GetChannelMessagesViews(channelId int32, idList []int32, increment bool) ([]int32) {
	viewsDOList := m.dao.ChannelMessagesDAO.SelectMessagesViews(channelId, idList)
	viewsList := make([]int32, 0, len(idList))

	for _, id := range idList {
		views := int32(1)
		for i := 0; i < len(viewsDOList); i++ {
			if viewsDOList[i].ChannelMessageId == id {
				if increment {
					views = viewsDOList[i].Views + 1
				} else {
					views = viewsDOList[i].Views
				}
				break
			}
		}
		viewsList = append(viewsList, views)
	}

	return viewsList
}

func (m *MessageModel) IncrementChannelMessagesViews(channelId int32, idList []int32) {
	m.dao.ChannelMessagesDAO.UpdateMessagesViews(channelId, idList)
}

/*
func (m *MessageModel) GetMessageByPeerAndMessageId(userId int32, messageId int32) (message *mtproto.Message) {
	do := m.dao.MessagesDAO.SelectByMessageId(userId, messageId)
	if do != nil {
		message, _ = messageDOToMessage(do)
	}
	return
}

func (m *MessageModel) GetMessageBoxListByMessageIdList(userId int32, idList []int32) []*MessageBox {
	doList := m.dao.MessagesDAO.SelectByMessageIdList(userId, idList)
	boxList := make([]*MessageBox, 0, len(doList))
	for _, do := range doList {
		message, _ := messageDOToMessage(&do)
		box := &MessageBox{
			UserId:    do.UserId,
			MessageId: do.UserMessageBoxId,
			Message:   message,
		}
		boxList = append(boxList, box)
	}
	return boxList
}

//func (m *MessageModel) GetPeerDialogMessageListByMessageId(userId int32, messageId int32) (messages *InboxMessages) {
//	doList := m.dao.MessagesDAO.SelectPeerDialogMessageListByMessageId(userId, messageId)
//	messages = &InboxMessages{
//		UserIds:  make([]int32, 0, len(doList)),
//		Messages: make([]*mtproto.Message, 0, len(doList)),
//	}
//	for _, do := range doList {
//		// TODO(@benqi): check data
//		m, _ := messageDOToMessage(&do)
//		messages.Messages = append(messages.Messages, m)
//		messages.UserIds = append(messages.UserIds, do.UserId)
//	}
//	return
//}

func (m *MessageModel) GetMessagesByPeerAndMessageIdList2(userId int32, idList []int32) (messages []*mtproto.Message) {
	if len(idList) == 0 {
		messages = []*mtproto.Message{}
	} else {
		doList := m.dao.MessagesDAO.SelectByMessageIdList(userId, idList)
		messages = make([]*mtproto.Message, 0, len(doList))
		for i := 0; i < len(doList); i++ {
			// TODO(@benqi): check data
			m, _ := messageDOToMessage(&doList[i])
			if m != nil {
				messages = append(messages, m)
			}
		}
	}
	return
}

/////////////////////////////////////
func (m *MessageModel) GetMessageIdListByDialog(userId int32, peer *base.PeerUtil) []int32 {
	doList := m.dao.MessagesDAO.SelectDialogMessageIdList(userId, peer.PeerId, int8(peer.PeerType))
	idList := make([]int32, 0, len(doList))
	for i := 0; i < len(doList); i++ {
		idList = append(idList, doList[i].UserMessageBoxId)
	}
	return idList
}

/////////////////////////////////////
func (m *MessageModel) GetPeerMessageId(userId, messageId, peerId int32) int32 {
	do := m.dao.MessagesDAO.SelectPeerMessageId(peerId, userId, messageId)
	if do == nil {
		return 0
	} else {
		return do.UserMessageBoxId
	}
}

func (m *MessageModel) SaveMessage(message *mtproto.Message, userId, messageId int32) error {
	var err error
	messageData, err := json.Marshal(message)
	m.dao.MessagesDAO.UpdateMessagesData(string(messageData), userId, messageId)
	return err
}

func (m *MessageModel) DeleteByMessageIdList(userId int32, idList []int32) {
	m.dao.MessagesDAO.DeleteMessagesByMessageIdList(userId, idList)
}

func (m *MessageModel) GetPeerDialogMessageIdList(userId int32, idList []int32) map[int32][]int32{
	doList := m.dao.MessagesDAO.SelectPeerDialogMessageIdList(userId, idList)
	deleteIdListMap := make(map[int32][]int32)
	for _, do := range doList {
		if messageIdList, ok := deleteIdListMap[do.UserId]; !ok {
			deleteIdListMap[do.UserId] = []int32{do.UserMessageBoxId}
		} else {
			deleteIdListMap[do.UserId] = append(messageIdList, do.UserMessageBoxId)
		}
	}

	return deleteIdListMap
}


*/
