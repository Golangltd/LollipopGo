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
	"github.com/nebulaim/telegramd/biz/core"
	update2 "github.com/nebulaim/telegramd/biz/core/update"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"golang.org/x/net/context"
	"github.com/nebulaim/telegramd/biz/core/message"
	"github.com/nebulaim/telegramd/server/sync/sync_client"
)

/*
 request:
	body: { channels_createChannel
	  flags: 1 [INT],
	  broadcast: YES [ BY BIT 0 IN FIELD flags ],
	  megagroup: [ SKIPPED BY BIT 1 IN FIELD flags ],
	  title: "channel-001" [STRING],
	  about: "channel-001" [STRING],
	},

 response:
    result: { updates
      updates: [ vector<0x0>
        { updateMessageID
          id: 1 [INT],
          random_id: 4140694743218070306 [LONG],
        },
        { updateChannel
          channel_id: 1261873857 [INT],
        },
        { updateReadChannelInbox
          channel_id: 1261873857 [INT],
          max_id: 1 [INT],
        },
        { updateNewChannelMessage
          message: { messageService
            flags: 16386 [INT],
            out: YES [ BY BIT 1 IN FIELD flags ],
            mentioned: [ SKIPPED BY BIT 4 IN FIELD flags ],
            media_unread: [ SKIPPED BY BIT 5 IN FIELD flags ],
            silent: [ SKIPPED BY BIT 13 IN FIELD flags ],
            post: YES [ BY BIT 14 IN FIELD flags ],
            id: 1 [INT],
            from_id: [ SKIPPED BY BIT 8 IN FIELD flags ],
            to_id: { peerChannel
              channel_id: 1261873857 [INT],
            },
            reply_to_msg_id: [ SKIPPED BY BIT 3 IN FIELD flags ],
            date: 1529328456 [INT],
            action: { messageActionChannelCreate
              title: "channel-001" [STRING],
            },
          },
          pts: 2 [INT],
          pts_count: 1 [INT],
        },
      ],
      users: [ vector<0x0> ],
      chats: [ vector<0x0>
        { channel
          flags: 8225 [INT],
          creator: YES [ BY BIT 0 IN FIELD flags ],
          left: [ SKIPPED BY BIT 2 IN FIELD flags ],
          editor: [ SKIPPED BY BIT 3 IN FIELD flags ],
          broadcast: YES [ BY BIT 5 IN FIELD flags ],
          verified: [ SKIPPED BY BIT 7 IN FIELD flags ],
          megagroup: [ SKIPPED BY BIT 8 IN FIELD flags ],
          restricted: [ SKIPPED BY BIT 9 IN FIELD flags ],
          democracy: [ SKIPPED BY BIT 10 IN FIELD flags ],
          signatures: [ SKIPPED BY BIT 11 IN FIELD flags ],
          min: [ SKIPPED BY BIT 12 IN FIELD flags ],
          id: 1261873857 [INT],
          access_hash: 18367393077902002260 [LONG],
          title: "channel-001" [STRING],
          username: [ SKIPPED BY BIT 6 IN FIELD flags ],
          photo: { chatPhotoEmpty },
          date: 1529328455 [INT],
          version: 0 [INT],
          restriction_reason: [ SKIPPED BY BIT 9 IN FIELD flags ],
          admin_rights: [ SKIPPED BY BIT 14 IN FIELD flags ],
          banned_rights: [ SKIPPED BY BIT 15 IN FIELD flags ],
        },
      ],
      date: 1529328455 [INT],
      seq: 0 [INT],
    },
*/

// channels.createChannel#f4893d7f flags:# broadcast:flags.0?true megagroup:flags.1?true title:string about:string = Updates;
func (s *ChannelsServiceImpl) ChannelsCreateChannel(ctx context.Context, request *mtproto.TLChannelsCreateChannel) (*mtproto.Updates, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("channels.createChannel#f4893d7f - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))


	// 1. 创建channel
	channel := s.ChannelModel.NewChannelLogicByCreateChannel(md.UserId, request.GetTitle(), request.GetAbout())

	peer := &base.PeerUtil{
		PeerType: base.PEER_CHANNEL,
		PeerId:   channel.GetChannelId(),
	}
	createChannelMessage := channel.MakeCreateChannelMessage(md.UserId)
	randomId := core.GetUUID()

	// 2. 创建channel createChannel message
	var boxList []*message.MessageBox2
	s.MessageModel.SendInternalMessage(md.UserId, peer, randomId, false, createChannelMessage, func(i int32, box2 *message.MessageBox2) {
		boxList = append(boxList, box2)
	})

	syncUpdates := update2.NewUpdatesLogic(md.UserId)
	pts := int32(core.NextChannelPtsId(peer.PeerId))
	ptsCount := int32(1)

	syncUpdates.AddUpdateNewMessage(pts, ptsCount, boxList[0].ToMessage(md.UserId))
	syncUpdates.AddChat(channel.ToChannel(md.UserId))

	sync_client.GetSyncClient().SyncUpdatesNotMe(md.UserId, md.AuthId, syncUpdates.ToUpdates())

	resultUpdates := syncUpdates
	resultUpdates.AddUpdateMessageId(boxList[0].MessageId, randomId)
	updateChannel := &mtproto.TLUpdateChannel{Data2: &mtproto.Update_Data{
		ChannelId: channel.GetChannelId(),
	}}
	resultUpdates.AddUpdate(updateChannel.To_Update())
	updateReadChannelInbox := &mtproto.TLUpdateReadChannelInbox{ Data2: &mtproto.Update_Data{
		ChannelId: channel.GetChannelId(),
		MaxId:     boxList[0].MessageId,
	}}
	resultUpdates.AddUpdate(updateReadChannelInbox.To_Update())

	glog.Infof("channels.createChannel#f4893d7f - reply: {%v}", resultUpdates)
	return resultUpdates.ToUpdates(), nil
}
