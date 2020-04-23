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
	"github.com/nebulaim/telegramd/server/sync/sync_client"
	"time"
	"github.com/nebulaim/telegramd/biz/core/message"
	"fmt"
	"github.com/nebulaim/telegramd/biz/core/update"
	"net/url"
)


// TODO(@benqi): mention...
//
func makeMessageBySendMessage(fromId, peerType, peerId int32, request *mtproto.TLMessagesSendMessage) (message *mtproto.TLMessage) {
	message = &mtproto.TLMessage{Data2: &mtproto.Message_Data{
		Out:          true,
		Silent:       request.GetSilent(),
		FromId:       fromId,
		ToId:         base.ToPeerByTypeAndID(int8(peerType), peerId),
		ReplyToMsgId: request.GetReplyToMsgId(),
		Message:      request.GetMessage(),
		ReplyMarkup:  request.GetReplyMarkup(),
		Entities:     request.GetEntities(),
		Date:         int32(time.Now().Unix()),
	}}

	// TODO(@benqi): check channel or super chat
	if peerType == base.PEER_CHANNEL {
		message.SetPost(true)
		message.SetViews(1)
	}

	//// TODO(@benqi):
	// isWebPageMessage = false

	message.Data2.Media = &mtproto.MessageMedia{
		Constructor: mtproto.TLConstructor_CRC32_messageMediaEmpty,
		Data2:       &mtproto.MessageMedia_Data{},
	}
	entities := request.GetEntities()

	u, err := url.Parse(request.Message)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		entityUrl := &mtproto.TLMessageEntityUrl{Data2: &mtproto.MessageEntity_Data{
			Offset: 0,
			Length: int32(len(request.Message)),
		}}
		entities = append(entities, entityUrl.To_MessageEntity())
	}
	message.Data2.Entities = entities

	return
}

func (s *MessagesServiceImpl) DoClearDraft(userId int32, authKeyId int64, peer *base.PeerUtil) {
	hasClearDraft := s.DialogModel.ClearDraftMessage(userId, peer.PeerType, peer.PeerId)

	// ClearDraft
	if hasClearDraft {
		updateDraftMessage := &mtproto.TLUpdateDraftMessage{Data2: &mtproto.Update_Data{
			Peer_39: peer.ToPeer(),
			Draft:   mtproto.NewTLDraftMessageEmpty().To_DraftMessage(),
		}}

		updates := &mtproto.TLUpdates{Data2: &mtproto.Updates_Data{
			Updates: []*mtproto.Update{updateDraftMessage.To_Update()},
			Users:   []*mtproto.User{},
			Chats:   []*mtproto.Chat{},
			Date:    int32(time.Now().Unix()),
			Seq:     0,
		}}

		sync_client.GetSyncClient().SyncUpdatesNotMe(userId, authKeyId, updates.To_Updates())
	}
}

// 流程：
//  1. 入库: 1）存消息数据，2）分别存到发件箱和收件箱里
//  2. 离线推送
//  3. 在线推送
// messages.sendMessage#fa88427a flags:# no_webpage:flags.1?true silent:flags.5?true background:flags.6?true clear_draft:flags.7?true peer:InputPeer reply_to_msg_id:flags.0?int message:string random_id:long reply_markup:flags.2?ReplyMarkup entities:flags.3?Vector<MessageEntity> = Updates;
func (s *MessagesServiceImpl) MessagesSendMessage(ctx context.Context, request *mtproto.TLMessagesSendMessage) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.sendMessage#fa88427a - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// peer
	var (
		peer               *base.PeerUtil
		err                error
	)

	if request.GetPeer().GetConstructor() == mtproto.TLConstructor_CRC32_inputPeerEmpty {
		err = mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error("messages.sendMessage#fa88427a - invalid peer", err)
		return nil, err
	}

	// TODO(@benqi): check user or channels's access_hash
	if request.GetPeer().GetConstructor() == mtproto.TLConstructor_CRC32_inputPeerSelf {
		peer = &base.PeerUtil{
			PeerType: base.PEER_USER,
			PeerId:   md.UserId,
		}
	} else {
		peer = base.FromInputPeer(request.GetPeer())
	}

	// 1. draft
	if request.GetClearDraft() {
		s.DoClearDraft(md.UserId, md.AuthId, peer)
	}

	outboxMessage := makeMessageBySendMessage(md.UserId, peer.PeerType, peer.PeerId, request)

	if peer.PeerType != base.PEER_CHANNEL {
		resultCB := func(pts, ptsCount int32, outBox *message.MessageBox2) (*mtproto.Updates, error) {
			sentMessage := message.MessageToUpdateShortSentMessage(outBox.ToMessage(md.UserId))
			sentMessage.SetPts(pts)
			sentMessage.SetPtsCount(ptsCount)
			return sentMessage.To_Updates(), nil
		}

		syncNotMeCB := func(pts, ptsCount int32, outBox *message.MessageBox2) (int64, *mtproto.Updates, error) {
			var syncUpdates *mtproto.Updates
			switch peer.PeerType {
			case base.PEER_USER:
				syncShortMessage := message.MessageToUpdateShortMessage(outBox.ToMessage(md.UserId))
				syncShortMessage.SetPts(pts)
				syncShortMessage.SetPtsCount(ptsCount)
				syncUpdates = syncShortMessage.To_Updates()
			case base.PEER_CHAT:
				syncShortMessage := message.MessageToUpdateShortChatMessage(outBox.ToMessage(md.UserId))
				syncShortMessage.SetPts(pts)
				syncShortMessage.SetPtsCount(ptsCount)
				syncUpdates = syncShortMessage.To_Updates()
			case base.PEER_CHANNEL:
				err = fmt.Errorf("peer_channel not impl")
			default:
				err = fmt.Errorf("invalid peer_type")
			}
			return md.AuthId, syncUpdates, nil
		}

		pushCB := func(pts, ptsCount int32, inBox *message.MessageBox2) (*mtproto.Updates, error) {
			var (
				updates *mtproto.Updates
				err error
			)

			glog.Info(inBox)
			switch peer.PeerType {
			case base.PEER_USER:
				shortMessage := message.MessageToUpdateShortMessage(inBox.ToMessage(inBox.OwnerId))
				shortMessage.SetPts(pts)
				shortMessage.SetPtsCount(ptsCount)
				updates = shortMessage.To_Updates()
			case base.PEER_CHAT:
				shortMessage := message.MessageToUpdateShortChatMessage(inBox.ToMessage(inBox.OwnerId))
				shortMessage.SetPts(pts)
				shortMessage.SetPtsCount(ptsCount)
				updates = shortMessage.To_Updates()
			case base.PEER_CHANNEL:
				err = fmt.Errorf("peer_channel not impl")
			default:
				err = fmt.Errorf("invalid peer_type")
			}
			return updates, err
		}

		replyUpdates, err := s.MessageModel.SendMessage(
			md.UserId,
			peer,
			request.GetRandomId(),
			outboxMessage.To_Message(),
			resultCB,
			syncNotMeCB,
			pushCB)

		glog.Infof("messages.sendMessage#fa88427a - reply: %s", logger.JsonDebugData(replyUpdates))
		return replyUpdates, err
	} else {
		channelLogic, _ := s.ChannelModel.NewChannelLogicById(peer.PeerId)
		resultCB := func(pts, ptsCount int32, channelBox *message.MessageBox2) *mtproto.Updates {
			replyUpdates := updates.NewUpdatesLogic(md.UserId)
			channelLogic.SetTopMessage(channelBox.MessageId)

			replyUpdates.AddUpdateMessageId(channelBox.MessageId, channelBox.RandomId)
			updateReadChannelInbox := &mtproto.TLUpdateReadChannelInbox{ Data2: &mtproto.Update_Data{
				ChannelId: channelBox.OwnerId,
				MaxId:     channelBox.MessageId,
			}}
			replyUpdates.AddUpdate(updateReadChannelInbox.To_Update())
			replyUpdates.AddUpdateNewChannelMessage(pts, ptsCount, channelBox.ToMessage(md.UserId))
			replyUpdates.AddChat(channelLogic.ToChannel(md.UserId))

			return replyUpdates.ToUpdates()
		}

		syncNotMeCB := func(pts, ptsCount int32, channelBox *message.MessageBox2) ([]int32, int64, *mtproto.Updates, error) {
			syncUpdates := updates.NewUpdatesLogic(md.UserId)

			updateReadChannelInbox := &mtproto.TLUpdateReadChannelInbox{ Data2: &mtproto.Update_Data{
				ChannelId: channelBox.OwnerId,
				MaxId:     channelBox.MessageId,
			}}
			syncUpdates.AddUpdate(updateReadChannelInbox.To_Update())
			syncUpdates.AddUpdateNewChannelMessage(pts, ptsCount, channelBox.ToMessage(md.UserId))
			syncUpdates.AddChat(channelLogic.ToChannel(md.UserId))

			idList := channelLogic.GetChannelParticipantIdList(md.UserId)
			return idList, md.AuthId, syncUpdates.ToUpdates(), nil
		}

		pushCB := func(userId, pts, ptsCount int32, channelBox *message.MessageBox2) (*mtproto.Updates, error) {
			pushUpdates := updates.NewUpdatesLogic(userId)

			pushUpdates.AddUpdateNewChannelMessage(pts, ptsCount, channelBox.ToMessage(userId))
			pushUpdates.AddChat(channelLogic.ToChannel(userId))

			return pushUpdates.ToUpdates(), nil
		}

		replyUpdates, err := s.MessageModel.SendChannelMessage(
			md.UserId,
			peer,
			request.GetRandomId(),
			outboxMessage.To_Message(),
			resultCB,
			syncNotMeCB,
			pushCB)

		glog.Infof("messages.sendMessage#fa88427a - reply: %s", logger.JsonDebugData(replyUpdates))
		return replyUpdates, err
	}
}
