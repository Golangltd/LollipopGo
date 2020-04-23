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

package chat

import (
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/proto/mtproto"
	// photo2 "github.com/nebulaim/telegramd/biz/core/photo"
	// "time"
	"github.com/nebulaim/telegramd/biz/base"
	// "github.com/nebulaim/telegramd/biz/core/account"
	// "github.com/nebulaim/telegramd/server/nbfs/nbfs_client"
	"github.com/nebulaim/telegramd/biz/core"
	"github.com/nebulaim/telegramd/biz/dal/dao"
	"github.com/nebulaim/telegramd/biz/dal/dao/mysql_dao"
)

type chatsDAO struct {
	*mysql_dao.CommonDAO
	*mysql_dao.UsersDAO
	*mysql_dao.ChatsDAO
	*mysql_dao.ChatParticipantsDAO
}

type ChatModel struct {
	dao           *chatsDAO
	photoCallback core.PhotoCallback
	notifySetting core.NotifySettingCallback
}

func (m *ChatModel) RegisterCallback(cb interface{}) {
	switch cb.(type) {
	case core.PhotoCallback:
		glog.Info("chatModel - register core.PhotoCallback")
		m.photoCallback = cb.(core.PhotoCallback)
	case core.NotifySettingCallback:
		glog.Info("chatModel - register core.NotifySettingCallback")
		m.notifySetting = cb.(core.NotifySettingCallback)
	}
}

func (m *ChatModel) InstallModel() {
	m.dao.CommonDAO = dao.GetCommonDAO(dao.DB_MASTER)
	m.dao.UsersDAO = dao.GetUsersDAO(dao.DB_MASTER)
	m.dao.ChatsDAO = dao.GetChatsDAO(dao.DB_MASTER)
	m.dao.ChatParticipantsDAO = dao.GetChatParticipantsDAO(dao.DB_MASTER)
}

func (m *ChatModel) GetChatListBySelfAndIDList(selfUserId int32, idList []int32) (chats []*mtproto.Chat) {
	if len(idList) == 0 {
		return []*mtproto.Chat{}
	}

	chats = make([]*mtproto.Chat, 0, len(idList))

	// TODO(@benqi): 性能优化，从DB里一次性取出所有的chatList
	for _, id := range idList {
		chatData, err := m.NewChatLogicById(id)
		if err != nil {
			glog.Error("getChatListBySelfIDList - not find chat_id: ", id)
			chatEmpty := &mtproto.TLChatEmpty{Data2: &mtproto.Chat_Data{
				Id: id,
			}}
			chats = append(chats, chatEmpty.To_Chat())
		} else {
			chats = append(chats, chatData.ToChat(selfUserId))
		}
	}

	return
}

func (m *ChatModel) GetChatBySelfID(selfUserId, chatId int32) (chat *mtproto.Chat) {
	chatData, err := m.NewChatLogicById(chatId)
	if err != nil {
		glog.Error("getChatBySelfID - not find chat_id: ", chatId)
		chatEmpty := &mtproto.TLChatEmpty{Data2: &mtproto.Chat_Data{
			Id: chatId,
		}}
		chat = chatEmpty.To_Chat()
	} else {
		chat = chatData.ToChat(selfUserId)
	}

	return
}

func (m *ChatModel) GetChatFullBySelfId(selfUserId int32, chatData *chatLogicData) *mtproto.TLChatFull {
	// sizes, _ := nbfs_client.GetPhotoSizeList(chatData.chat.PhotoId)
	// photo2 := photo2.MakeUserProfilePhoto(photoId, sizes)
	var photo *mtproto.Photo

	if chatData.GetPhotoId() == 0 {
		photoEmpty := &mtproto.TLPhotoEmpty{Data2: &mtproto.Photo_Data{
			Id: 0,
		}}
		photo = photoEmpty.To_Photo()
	} else {
		//chatPhoto := &mtproto.TLPhoto{ Data2: &mtproto.Photo_Data{
		//	Id:          chatData.chat.PhotoId,
		//	HasStickers: false,
		//	AccessHash:  chatData.chat.PhotoId, // photo2.GetFileAccessHash(file.GetData2().GetId(), file.GetData2().GetParts()),
		//	Date:        int32(time.Now().Unix()),
		//	Sizes:       sizes,
		//}}
		//photo = chatPhoto.To_Photo()
		photo = m.photoCallback.GetPhoto(chatData.chat.PhotoId)
	}

	peer := &base.PeerUtil{
		PeerType: base.PEER_CHAT,
		PeerId:   chatData.GetChatId(),
	}
	notifySettings := m.notifySetting.GetNotifySettings(selfUserId, peer)

	chatFull := &mtproto.TLChatFull{Data2: &mtproto.ChatFull_Data{
		Id:             chatData.GetChatId(),
		Participants:   chatData.GetChatParticipants().To_ChatParticipants(),
		ChatPhoto:      photo,
		NotifySettings: notifySettings,
		ExportedInvite: mtproto.NewTLChatInviteEmpty().To_ExportedChatInvite(), // TODO(@benqi):
		BotInfo:        []*mtproto.BotInfo{},
	}}

	return chatFull
}

func init() {
	core.RegisterCoreModel(&ChatModel{dao: &chatsDAO{}})
}
