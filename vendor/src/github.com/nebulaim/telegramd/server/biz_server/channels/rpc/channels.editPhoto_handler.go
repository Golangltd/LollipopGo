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
	"github.com/nebulaim/telegramd/service/document/client"
	"time"
	"github.com/nebulaim/telegramd/biz/base"
	"github.com/nebulaim/telegramd/biz/core"
	"github.com/nebulaim/telegramd/biz/core/message"
	"github.com/nebulaim/telegramd/biz/core/channel"
)

// channels.editPhoto#f12e57c9 channel:InputChannel photo:InputChatPhoto = Updates;
func (s *ChannelsServiceImpl) ChannelsEditPhoto(ctx context.Context, request *mtproto.TLChannelsEditPhoto) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.editPhoto#f12e57c9 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): check channel_id and access_hash

	peer := &base.PeerUtil{
		PeerType: base.PEER_CHANNEL,
		PeerId:   request.GetChannel().GetData2().GetChannelId(),
	}

	var (
		photoId int64 = 0
		action  *mtproto.MessageAction
	)

	channelLogic, err := s.ChannelModel.NewChannelLogicById(peer.PeerId)
	if err != nil {
		glog.Error("messages.editChatTitle#dc452855 - error: ", err)
		return nil, err
	}

	chatPhoto := request.GetPhoto()
	switch chatPhoto.GetConstructor() {
	case mtproto.TLConstructor_CRC32_inputChatPhotoEmpty:
		photoId = 0
		action = mtproto.NewTLMessageActionChatDeletePhoto().To_MessageAction()
	case mtproto.TLConstructor_CRC32_inputChatUploadedPhoto:
		file := chatPhoto.GetData2().GetFile()
		// photoId = helper.NextSnowflakeId()
		result, err := document_client.UploadPhotoFile(md.AuthId, file) // photoId, file.GetData2().GetId(), file.GetData2().GetParts(), file.GetData2().GetName(), file.GetData2().GetMd5Checksum())
		if err != nil {
			glog.Errorf("UploadPhoto error: %v", err)
			return nil, err
		}
		photoId = result.PhotoId
		// user.SetUserPhotoID(md.UserId, uuid)
		// fileData := mediaData.GetFile().GetData2()
		photo := &mtproto.TLPhoto{Data2: &mtproto.Photo_Data{
			Id:          photoId,
			HasStickers: false,
			AccessHash:  result.AccessHash, // photo2.GetFileAccessHash(file.GetData2().GetId(), file.GetData2().GetParts()),
			Date:        int32(time.Now().Unix()),
			Sizes:       result.SizeList,
		}}
		action2 := &mtproto.TLMessageActionChatEditPhoto{Data2: &mtproto.MessageAction_Data{
			Photo: photo.To_Photo(),
		}}
		action = action2.To_MessageAction()
	case mtproto.TLConstructor_CRC32_inputChatPhoto:
		// photo := chatPhoto.GetData2().GetId()
	}

	channelLogic.EditPhoto(md.UserId, photoId)
	editChannelPhotoMessage := channel.MakeChannelMessageService(md.UserId, channelLogic.GetChannelId(), action)
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
		editChannelPhotoMessage,
		resultCB,
		syncNotMeCB,
		pushCB)

	// reply := replyUpdates.To_Updates()
	glog.Infof("channels.editPhoto#f12e57c9 - reply: {%s}", logger.JsonDebugData(replyUpdates))
	return replyUpdates, nil
}
