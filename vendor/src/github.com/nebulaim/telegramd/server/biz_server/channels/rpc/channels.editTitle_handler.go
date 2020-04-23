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
	"github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/biz/core"
	"github.com/nebulaim/telegramd/biz/core/message"
	"github.com/nebulaim/telegramd/biz/base"
)

// channels.editTitle#566decd0 channel:InputChannel title:string = Updates;
func (s *ChannelsServiceImpl) ChannelsEditTitle(ctx context.Context, request *mtproto.TLChannelsEditTitle) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.editTitle#566decd0 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): check channel_id and access_hash

	// channelId := request.GetChannel().GetData2().GetChannelId()
	peer := &base.PeerUtil{
		PeerType: base.PEER_CHANNEL,
		PeerId:   request.GetChannel().GetData2().GetChannelId(),
	}

	channelLogic, _ := s.ChannelModel.NewChannelLogicById(peer.PeerId)
	err := channelLogic.EditTitle(md.UserId, request.GetTitle())
	if err != nil {
		glog.Errorf("channels.editTitle#566decd0 - error: %s", err)
		return nil, err
	}

	channelEditMessage := channelLogic.MakeChannelEditTitleMessage(md.UserId, request.Title)
	randomId := core.GetUUID()

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
		randomId,
		channelEditMessage,
		resultCB,
		syncNotMeCB,
		pushCB)

	// reply := replyUpdates.To_Updates()
	glog.Infof("channels.editTitle#566decd0 - reply: {%s}", logger.JsonDebugData(replyUpdates))
	return replyUpdates, nil
}
