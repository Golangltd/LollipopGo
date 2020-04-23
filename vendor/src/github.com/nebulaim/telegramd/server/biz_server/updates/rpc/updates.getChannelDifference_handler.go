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

package rpc

import (
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"golang.org/x/net/context"
)

// updates.getChannelDifference#3173d78 flags:# force:flags.0?true channel:InputChannel filter:ChannelMessagesFilter pts:int limit:int = updates.ChannelDifference;
func (s *UpdatesServiceImpl) UpdatesGetChannelDifference(ctx context.Context, request *mtproto.TLUpdatesGetChannelDifference) (*mtproto.Updates_ChannelDifference, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("updates.getChannelDifference#3173d78 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
	// lastPts = request.GetPts()
	//otherUpdates []*mtproto.Update
	// messages []*mtproto.Message
	//userList []*mtproto.User
	//chatList []*mtproto.Chat

	)

	// var difference *mtproto.Updates_ChannelDifference

	channelId := request.GetChannel().GetData2().ChannelId
	channelLogic, _ := s.ChannelModel.NewChannelLogicById(channelId)
	participant := channelLogic.GetChannelParticipant(md.UserId)
	switch participant.GetConstructor() {
	case mtproto.TLConstructor_CRC32_channelParticipantsBanned:
		// TODO(@benqi):
		//banned := channel.MakeChannelBannedRights(participant.GetData2().GetBannedRights().To_ChannelBannedRights())
		//if banned.IsForbidden() {
		//	return nil, mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_CHANNEL_PRIVATE)
		//}
		return nil, mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_CHANNEL_PRIVATE)

	}

	// updateList, _ := sync_client.GetSyncClient().SyncGetDifference()
	// .GetChannelUpdateListByGtPts(request.GetChannel().GetData2().GetChannelId(), lastPts)
	//
	//for _, update := range updateList {
	//	switch update.GetConstructor() {
	//	case mtproto.TLConstructor_CRC32_updateNewChannelMessage:
	//		newMessage := update.To_UpdateNewChannelMessage()
	//		messages = append(messages, newMessage.GetMessage())
	//		// otherUpdates = append(otherUpdates, update)
	//
	//	case mtproto.TLConstructor_CRC32_updateDeleteChannelMessages:
	//		// readHistoryOutbox := update.To_UpdateReadHistoryOutbox()
	//		// readHistoryOutbox.SetPtsCount(0)
	//		// otherUpdates = append(otherUpdates, readHistoryOutbox.To_Update())
	//	case mtproto.TLConstructor_CRC32_updateEditChannelMessage:
	//		// readHistoryInbox := update.To_UpdateReadHistoryInbox()
	//		// readHistoryInbox.SetPtsCount(0)
	//		// otherUpdates = append(otherUpdates, readHistoryInbox.To_Update())
	//	case mtproto.TLConstructor_CRC32_updateChannelWebPage:
	//	default:
	//		continue
	//	}
	//	if update.Data2.GetPts() > lastPts {
	//		lastPts = update.Data2.GetPts()
	//	}
	//}

	//otherUpdates, boxIDList, lastPts := model.GetUpdatesModel().GetUpdatesByGtPts(md.UserId, request.GetPts())
	//messages := model.GetMessageModel().GetMessagesByPeerAndMessageIdList2(md.UserId, boxIDList)
	// userIdList, chatIdList, _ := message.PickAllIDListByMessages(messages)
	// userList = user.GetUsersBySelfAndIDList(md.UserId, userIdList)
	// chatList = chat.GetChatListBySelfAndIDList(md.UserId, chatIdList)
	//
	//state := &mtproto.TLUpdatesState{Data2: &mtproto.Updates_State_Data{
	//	Pts:         lastPts,
	//	Date:        int32(time.Now().Unix()),
	//	UnreadCount: 0,
	//	Seq:         int32(model.GetSequenceModel().CurrentSeqId(base2.Int32ToString(md.UserId))),
	//	Seq:         0,
	//}}

	// updates.channelDifference#2064674e flags:# final:flags.0?true pts:int timeout:flags.1?int new_messages:Vector<Message> other_updates:Vector<Update> chats:Vector<Chat> users:Vector<User> = updates.ChannelDifference;
	var difference *mtproto.Updates_ChannelDifference

	//if len(updateList) == 0 {
	// pts, _ := sync_client.GetSyncClient().GetCurrentChannelPts(request.GetChannel().GetData2().GetChannelId())
	difference = &mtproto.Updates_ChannelDifference{
		Constructor: mtproto.TLConstructor_CRC32_updates_channelDifferenceEmpty,
		Data2: &mtproto.Updates_ChannelDifference_Data{
			Final:   true,
			Pts:     0,
			Timeout: 30,
		},
	}
	//} else {
	//	difference = &mtproto.Updates_ChannelDifference{
	//		Constructor: mtproto.TLConstructor_CRC32_updates_channelDifferenceEmpty,
	//		Data2:  &mtproto.Updates_ChannelDifference_Data{
	//			Final:   true,
	//			Pts:     2,
	//			Timeout: 3,
	//		},
	//	}
	//	//difference := &mtproto.TLUpdatesChannelDifference{Data2: &mtproto.Updates_ChannelDifference_Data{
	//	//	Pts: lastPts,
	//	//	Timeout: 3,
	//	//	NewMessages:  messages,
	//	//	OtherUpdates: otherUpdates,
	//	//	Users:        userList,
	//	//	Chats:        chatList,
	//	//	// State:        state.To_Updates_State(),
	//	//}}
	//	//
	//	// TODO(@benqi): remove to received ack handler.
	//	// update2.UpdateAuthStateSeq(md.AuthId, lastPts, 0)
	//}

	glog.Infof("updates.getChannelDifference#3173d78 - reply: %s", logger.JsonDebugData(difference))
	return difference, nil

}
