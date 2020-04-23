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
	"github.com/nebulaim/telegramd/biz/base"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"golang.org/x/net/context"
	"github.com/nebulaim/telegramd/biz/core"
	"github.com/nebulaim/telegramd/biz/core/message"
	"github.com/nebulaim/telegramd/biz/core/update"
)

// messages.createChat#9cb126e users:Vector<InputUser> title:string = Updates;
func (s *MessagesServiceImpl) MessagesCreateChat(ctx context.Context, request *mtproto.TLMessagesCreateChat) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.createChat#9cb126e - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// logic.NewChatLogicByCreateChat()
	//// TODO(@benqi): Impl MessagesCreateChat logic
	//randomId := md.ClientMsgId

	chatUserIdList := make([]int32, 0, len(request.GetUsers()))
	// chatUserIdList = append(chatUserIdList, md.UserId)
	for _, u := range request.GetUsers() {
		switch u.GetConstructor() {
		case mtproto.TLConstructor_CRC32_inputUser:
			chatUserIdList = append(chatUserIdList, u.GetData2().GetUserId())
		default:
			// TODO(@benqi): chatUser不能是inputUser和inputUserSelf
			err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
			glog.Error("messages.createChat#9cb126e - error: ", err, "; InputPeer invalid")
			return nil, err
		}
	}

	chat := s.ChatModel.NewChatLogicByCreateChat(md.UserId, chatUserIdList, request.GetTitle())

	peer := &base.PeerUtil{
		PeerType: base.PEER_CHAT,
		PeerId:   chat.GetChatId(),
	}

	createChatMessage := chat.MakeCreateChatMessage(md.UserId)
	randomId := core.GetUUID()

	resultCB := func(pts, ptsCount int32, outBox *message.MessageBox2) (*mtproto.Updates, error) {
		syncUpdates := updates.NewUpdatesLogic(md.UserId)

		updateChatParticipants := &mtproto.TLUpdateChatParticipants{Data2: &mtproto.Update_Data{
			Participants: chat.GetChatParticipants().To_ChatParticipants(),
		}}
		syncUpdates.AddUpdate(updateChatParticipants.To_Update())
		syncUpdates.AddUpdateNewMessage(pts, ptsCount, outBox.ToMessage(md.UserId))
		syncUpdates.AddUsers(s.UserModel.GetUsersBySelfAndIDList(md.UserId, chat.GetChatParticipantIdList()))
		syncUpdates.AddChat(chat.ToChat(md.UserId))
		syncUpdates.AddUpdateMessageId(outBox.MessageId, outBox.RandomId)

		return syncUpdates.ToUpdates(), nil
	}

	syncNotMeCB := func(pts, ptsCount int32, outBox *message.MessageBox2) (int64, *mtproto.Updates, error) {
		syncUpdates := updates.NewUpdatesLogic(md.UserId)

		updateChatParticipants := &mtproto.TLUpdateChatParticipants{Data2: &mtproto.Update_Data{
			Participants: chat.GetChatParticipants().To_ChatParticipants(),
		}}
		syncUpdates.AddUpdate(updateChatParticipants.To_Update())
		syncUpdates.AddUpdateNewMessage(pts, ptsCount, outBox.ToMessage(md.UserId))
		syncUpdates.AddUsers(s.UserModel.GetUsersBySelfAndIDList(md.UserId, chat.GetChatParticipantIdList()))
		syncUpdates.AddChat(chat.ToChat(md.UserId))

		return md.AuthId, syncUpdates.ToUpdates(), nil
	}

	pushCB := func(pts, ptsCount int32, inBox *message.MessageBox2) (*mtproto.Updates, error) {
		pushUpdates := updates.NewUpdatesLogic(md.UserId)
		updateChatParticipants := &mtproto.TLUpdateChatParticipants{Data2: &mtproto.Update_Data{
			Participants: chat.GetChatParticipants().To_ChatParticipants(),
		}}
		pushUpdates.AddUpdate(updateChatParticipants.To_Update())
		pushUpdates.AddUpdateNewMessage(pts, ptsCount, inBox.ToMessage(inBox.OwnerId))
		pushUpdates.AddUsers(s.UserModel.GetUsersBySelfAndIDList(inBox.OwnerId, chat.GetChatParticipantIdList()))
		pushUpdates.AddChat(chat.ToChat(inBox.OwnerId))

		return pushUpdates.ToUpdates(), nil
	}

	replyUpdates, _ := s.MessageModel.SendMessage(
		md.UserId,
		peer,
		randomId,
		createChatMessage,
		resultCB,
		syncNotMeCB,
		pushCB)

	glog.Infof("messages.createChat#9cb126e - reply: {%v}", replyUpdates)
	return replyUpdates, nil
}
