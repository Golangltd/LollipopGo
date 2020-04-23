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

// channelParticipantsRecent#de3f3c79 = ChannelParticipantsFilter;
// channelParticipantsAdmins#b4608969 = ChannelParticipantsFilter;
// channelParticipantsKicked#a3b54985 q:string = ChannelParticipantsFilter;
// channelParticipantsBots#b0d1865b = ChannelParticipantsFilter;
// channelParticipantsBanned#1427a5e1 q:string = ChannelParticipantsFilter;
// channelParticipantsSearch#656ac4b q:string = ChannelParticipantsFilter;
//
// channels.getParticipants#123e05e9 channel:InputChannel filter:ChannelParticipantsFilter offset:int limit:int hash:int = channels.ChannelParticipants;
func (s *ChannelsServiceImpl) ChannelsGetParticipants(ctx context.Context, request *mtproto.TLChannelsGetParticipants) (*mtproto.Channels_ChannelParticipants, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.getParticipants#123e05e9 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	if request.GetChannel().GetConstructor() == mtproto.TLConstructor_CRC32_inputChannelEmpty {
		glog.Errorf("channels.getParticipants#123e05e9 - channel is inputChannelEmpty")
		return nil, mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
	}

	// TODO(@benqi): check access_hash
	channelId := request.GetChannel().GetData2().GetChannelId()

	channelLogic, _ := s.ChannelModel.NewChannelLogicById(channelId)

	participants := make([]*mtproto.ChannelParticipant, 0)

	switch request.GetFilter().GetConstructor() {
	case mtproto.TLConstructor_CRC32_channelParticipantsRecent:
		participants = channelLogic.GetChannelParticipantListRecent(request.GetOffset(), request.GetLimit(), request.GetHash())
	case mtproto.TLConstructor_CRC32_channelParticipantsAdmins:
		participants = channelLogic.GetChannelParticipantListAdmins(request.GetOffset(), request.GetLimit(), request.GetHash())
	case mtproto.TLConstructor_CRC32_channelParticipantsKicked:
		participants = channelLogic.GetChannelParticipantListKicked(request.GetFilter().GetData2().GetQ(), request.GetOffset(), request.GetLimit(), request.GetHash())
	case mtproto.TLConstructor_CRC32_channelParticipantsBots:
		participants = channelLogic.GetChannelParticipantListRecent(request.GetOffset(), request.GetLimit(), request.GetHash())
	case mtproto.TLConstructor_CRC32_channelParticipantsBanned:
		participants = channelLogic.GetChannelParticipantListBanned(request.GetFilter().GetData2().GetQ(), request.GetOffset(), request.GetLimit(), request.GetHash())
	case mtproto.TLConstructor_CRC32_channelParticipantsSearch:
		participants = channelLogic.GetChannelParticipantListSearch(request.GetFilter().GetData2().GetQ(), request.GetOffset(), request.GetLimit(), request.GetHash())
	default:
		glog.Errorf("channels.getParticipants#123e05e9 - channel is inputChannelEmpty")
	}

	var userIdList []int32
	for _, participant := range participants {
		userIdList = append(userIdList, participant.GetData2().GetUserId())
	}
	users := s.UserModel.GetUsersBySelfAndIDList(md.UserId, userIdList)

	channelParticipants := &mtproto.TLChannelsChannelParticipants{Data2: &mtproto.Channels_ChannelParticipants_Data{
		Count:        int32(len(participants)),
		Participants: participants,
		Users:        users,
	}}

	glog.Infof("channels.getParticipants#123e05e9 - reply: {%s}", logger.JsonDebugData(channelParticipants))
	return channelParticipants.To_Channels_ChannelParticipants(), nil
}
