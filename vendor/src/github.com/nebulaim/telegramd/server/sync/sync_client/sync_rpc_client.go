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

package sync_client

import (
	"context"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/grpc_util/service_discovery"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"time"
)

type syncClient struct {
	client mtproto.RPCSyncClient
}

var (
	syncInstance = &syncClient{}
)

func GetSyncClient() *syncClient {
	return syncInstance
}

func InstallSyncClient(discovery *service_discovery.ServiceDiscoveryClientConfig) {
	conn, err := grpc_util.NewRPCClientByServiceDiscovery(discovery)

	if err != nil {
		glog.Error(err)
		panic(err)
	}

	syncInstance.client = mtproto.NewRPCSyncClient(conn)
}

// sync.syncUpdates#3a077679 flags:# layer:int user_id:int auth_key_id:long server_id:flags.1?int not_me:flags.0?true updates:Updates = Bool;
func (c *syncClient) SyncUpdatesMe(userId int32, authKeyId int64, serverId int32, updates *mtproto.Updates) (bool, error) {
	m := &mtproto.TLSyncSyncUpdates{
		UserId:    userId,
		AuthKeyId: authKeyId,
		ServerId:  serverId,
		Updates:   updates,
	}

	r, err := c.client.SyncSyncUpdates(context.Background(), m)
	return mtproto.FromBool(r), err
}

// sync.syncUpdates#3a077679 flags:# layer:int user_id:int auth_key_id:long server_id:flags.1?int not_me:flags.0?true updates:Updates = Bool;
func (c *syncClient) SyncUpdatesNotMe(userId int32, authKeyId int64, updates *mtproto.Updates) (bool, error) {
	m := &mtproto.TLSyncSyncUpdates{
		UserId:    userId,
		AuthKeyId: authKeyId,
		Updates:   updates,
	}

	r, err := c.client.SyncSyncUpdates(context.Background(), m)
	return mtproto.FromBool(r), err
}

// sync.pushUpdates#5c612649 user_id:int updates:Updates = Bool;
func (c *syncClient) PushUpdates(userId int32, updates *mtproto.Updates) (bool, error) {
	m := &mtproto.TLSyncPushUpdates{
		UserId:  userId,
		Updates: updates,
	}

	r, err := c.client.SyncPushUpdates(context.Background(), m)
	return mtproto.FromBool(r), err
}

func (c *syncClient) SyncChannelUpdatesMe(channelId int32, participantId int32, authKeyId int64, serverId int32, updates *mtproto.Updates) (bool, error) {
	m := &mtproto.TLSyncSyncChannelUpdates{
		ChannelId: channelId,
		UserId:    participantId,
		AuthKeyId: authKeyId,
		ServerId:  serverId,
		Updates:   updates,
	}

	r, err := c.client.SyncSyncChannelUpdates(context.Background(), m)
	return mtproto.FromBool(r), err
}

func (c *syncClient) SyncChannelUpdatesNotMe(channelId int32, participantId int32, authKeyId int64, updates *mtproto.Updates) (bool, error) {
	m := &mtproto.TLSyncSyncChannelUpdates{
		ChannelId: channelId,
		UserId:    participantId,
		AuthKeyId: authKeyId,
		Updates:   updates,
	}

	r, err := c.client.SyncSyncChannelUpdates(context.Background(), m)
	return mtproto.FromBool(r), err
}

func (c *syncClient) PushChannelUpdates(channelId, userId int32, updates *mtproto.Updates) (bool, error) {
	m := &mtproto.TLSyncPushChannelUpdates{
		ChannelId: channelId,
		UserId:    userId,
		Updates:   updates,
	}

	r, err := c.client.SyncPushChannelUpdates(context.Background(), m)
	return mtproto.FromBool(r), err
}

//
//// sync.pushChannelUpdates#bfd3d677 channel_id:int exclude_user_id:int updates:Updates = Bool;
//func (c *syncClient) PushChannelUpdates(channelId int32, channelParticipantIds []int32, updates *mtproto.Updates) (bool, error) {
//	m := &mtproto.TLSyncPushChannelUpdates{
//		ChannelId:             channelId,
//		ChannelParticipantIds: channelParticipantIds,
//		Updates:               updates,
//	}
//
//	r, err := c.client.SyncPushChannelUpdates(context.Background(), m)
//	return mtproto.FromBool(r), err
//}

// sync.pushRpcResult#1bf9b15e auth_key_id:long req_msg_id:long result:bytes = Bool;
func (c *syncClient) SyncPushRpcResult(authKeyId int64, serverId int32, clientReqMsgId int64, result []byte) (bool, error) {
	m := &mtproto.TLSyncPushRpcResult{
		AuthKeyId: authKeyId,
		ServerId:  serverId,
		ReqMsgId:  clientReqMsgId,
		Result:    result,
	}

	r, err := c.client.SyncPushRpcResult(context.Background(), m)
	return mtproto.FromBool(r), err
}

// sync.getState auth_key_id:long user_id:int = updates.State;
func (c *syncClient) SyncGetState(authKeyId int64, userId int32) (*mtproto.Updates_State, error) {
	req := &mtproto.TLSyncGetState{
		AuthKeyId: authKeyId,
		UserId:    userId,
	}

	state, err := c.client.SyncGetState(context.Background(), req)
	return state, err
}

// sync.getDifference flags:# auth_key_id:long user_id:int pts:int pts_total_limit:flags.0?int date:int qts:int = updates.Difference;
func (c *syncClient) SyncGetDifference(authKeyId int64, userId, pts int32) (*mtproto.Updates_Difference, error) {
	req := &mtproto.TLSyncGetDifference{
		AuthKeyId: authKeyId,
		UserId:    userId,
		Pts:       pts,
		Date:      int32(time.Now().Unix()),
		Qts:       0,
	}

	difference, err := c.client.SyncGetDifference(context.Background(), req)
	return difference, err
}


func (c *syncClient) SyncGetChannelDifference(authKeyId int64, userId, pts int32) (*mtproto.Updates_ChannelDifference, error) {
	req := &mtproto.TLSyncGetChannelDifference{
		AuthKeyId: authKeyId,
		UserId:    userId,
		Pts:       pts,
		// Date:      int32(time.Now().Unix()),
		// Qts:       0,
	}

	difference, err := c.client.SyncGetChannelDifference(context.Background(), req)
	return difference, err
}


/*
func (c *syncClient) SyncOneUpdateData2(serverId int32, authKeyId, sessionId int64, pushUserId int32, clientMsgId int64, update *mtproto.Update) (reply *mtproto.ClientUpdatesState, err error) {
	updates := &mtproto.TLUpdates{Data2: &mtproto.Updates_Data{
		Updates: []*mtproto.Update{update},
	}}

	m := &mtproto.UpdatesRequest{
		PushType:    mtproto.SyncType_SYNC_TYPE_RPC_RESULT,
		ServerId:    serverId,
		AuthKeyId:   authKeyId,
		SessionId:   sessionId,
		PushUserId:  pushUserId,
		ClientMsgId: clientMsgId,
		Updates:     updates.To_Updates(),
		RpcResult: &mtproto.RpcResultData{
			AffectedMessages: mtproto.NewTLMessagesAffectedMessages(),
		},
	}
	reply, err = c.client.SyncUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) SyncOneUpdateData3(serverId int32, authKeyId, sessionId int64, pushUserId int32, clientMsgId int64, update *mtproto.Update) (reply *mtproto.ClientUpdatesState, err error) {
	updates := &mtproto.TLUpdates{Data2: &mtproto.Updates_Data{
		Updates: []*mtproto.Update{update},
	}}

	m := &mtproto.UpdatesRequest{
		PushType:    mtproto.SyncType_SYNC_TYPE_RPC_RESULT,
		ServerId:    serverId,
		AuthKeyId:   authKeyId,
		SessionId:   sessionId,
		PushUserId:  pushUserId,
		ClientMsgId: clientMsgId,
		Updates:     updates.To_Updates(),
		RpcResult: &mtproto.RpcResultData{
			AffectedHistory: mtproto.NewTLMessagesAffectedHistory(),
		},
	}
	reply, err = c.client.SyncUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) SyncOneUpdateData(authKeyId, sessionId int64, pushUserId int32, update *mtproto.Update) (reply *mtproto.ClientUpdatesState, err error) {
	updates := &mtproto.TLUpdates{Data2: &mtproto.Updates_Data{
		Updates: []*mtproto.Update{update},
	}}

	m := &mtproto.UpdatesRequest{
		PushType:   mtproto.SyncType_SYNC_TYPE_USER_NOTME,
		AuthKeyId:  authKeyId,
		SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates.To_Updates(),
	}
	reply, err = c.client.SyncUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) PushToUserNotMeOneUpdateData(authKeyId, sessionId int64, pushUserId int32, update *mtproto.Update) (reply *mtproto.VoidRsp, err error) {
	updates := &mtproto.TLUpdates{Data2: &mtproto.Updates_Data{
		Updates: []*mtproto.Update{update},
	}}

	m := &mtproto.UpdatesRequest{
		PushType:   mtproto.SyncType_SYNC_TYPE_USER_NOTME,
		AuthKeyId:  authKeyId,
		SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates.To_Updates(),
	}
	reply, err = c.client.PushUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) PushToUserMeOneUpdateData(authKeyId, sessionId int64, pushUserId int32, update *mtproto.Update) (reply *mtproto.VoidRsp, err error) {
	updates := &mtproto.TLUpdates{Data2: &mtproto.Updates_Data{
		Updates: []*mtproto.Update{update},
	}}

	m := &mtproto.UpdatesRequest{
		PushType:   mtproto.SyncType_SYNC_TYPE_USER_ME,
		AuthKeyId:  authKeyId,
		SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates.To_Updates(),
	}
	reply, err = c.client.PushUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) PushToUserOneUpdateData(pushUserId int32, update *mtproto.Update) (reply *mtproto.VoidRsp, err error) {
	updates := &mtproto.TLUpdates{Data2: &mtproto.Updates_Data{
		Updates: []*mtproto.Update{update},
	}}

	m := &mtproto.UpdatesRequest{
		PushType: mtproto.SyncType_SYNC_TYPE_USER,
		// AuthKeyId:  authKeyId,
		// SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates.To_Updates(),
	}
	reply, err = c.client.PushUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) PushToUserNotMeUpdateShortData(authKeyId, sessionId int64, pushUserId int32, update *mtproto.Update) (reply *mtproto.VoidRsp, err error) {
	updates := &mtproto.TLUpdateShort{Data2: &mtproto.Updates_Data{
		Update: update,
	}}

	m := &mtproto.UpdatesRequest{
		PushType:   mtproto.SyncType_SYNC_TYPE_USER_NOTME,
		AuthKeyId:  authKeyId,
		SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates.To_Updates(),
	}
	reply, err = c.client.PushUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) PushToUserMeUpdateShortData(authKeyId, sessionId int64, pushUserId int32, update *mtproto.Update) (reply *mtproto.VoidRsp, err error) {
	updates := &mtproto.TLUpdateShort{Data2: &mtproto.Updates_Data{
		Update: update,
	}}

	m := &mtproto.UpdatesRequest{
		PushType:   mtproto.SyncType_SYNC_TYPE_USER_ME,
		AuthKeyId:  authKeyId,
		SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates.To_Updates(),
	}
	reply, err = c.client.PushUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) PushToUserUpdateShortData(pushUserId int32, update *mtproto.Update) (reply *mtproto.VoidRsp, err error) {
	updates := &mtproto.TLUpdateShort{Data2: &mtproto.Updates_Data{
		Update: update,
	}}

	m := &mtproto.UpdatesRequest{
		PushType: mtproto.SyncType_SYNC_TYPE_USER,
		// AuthKeyId:  authKeyId,
		// SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates.To_Updates(),
	}
	reply, err = c.client.PushUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) SyncUpdatesData(authKeyId, sessionId int64, pushUserId int32, updates *mtproto.Updates) (reply *mtproto.ClientUpdatesState, err error) {
	m := &mtproto.UpdatesRequest{
		PushType:   mtproto.SyncType_SYNC_TYPE_USER_NOTME,
		AuthKeyId:  authKeyId,
		SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates,
	}
	reply, err = c.client.SyncUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) PushToUserNotMeUpdatesData(authKeyId, sessionId int64, pushUserId int32, updates *mtproto.Updates) (reply *mtproto.VoidRsp, err error) {
	m := &mtproto.UpdatesRequest{
		PushType:   mtproto.SyncType_SYNC_TYPE_USER_NOTME,
		AuthKeyId:  authKeyId,
		SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates,
	}
	reply, err = c.client.PushUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) PushToUserMeUpdatesData(authKeyId, sessionId int64, pushUserId int32, updates *mtproto.Updates) (reply *mtproto.VoidRsp, err error) {
	m := &mtproto.UpdatesRequest{
		PushType:   mtproto.SyncType_SYNC_TYPE_USER_ME,
		AuthKeyId:  authKeyId,
		SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates,
	}
	reply, err = c.client.PushUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) PushToUserUpdatesData(pushUserId int32, updates *mtproto.Updates) (reply *mtproto.VoidRsp, err error) {
	m := &mtproto.UpdatesRequest{
		PushType: mtproto.SyncType_SYNC_TYPE_USER,
		// AuthKeyId:  authKeyId,
		// SessionId:  sessionId,
		PushUserId: pushUserId,
		Updates:    updates,
	}
	reply, err = c.client.PushUpdatesData(context.Background(), m)
	return
}

func (c *syncClient) GetCurrentChannelPts(channelId int32) (pts int32, err error) {
	req := &mtproto.ChannelPtsRequest{
		ChannelId: channelId,
	}
	var ptsId *mtproto.SeqId
	ptsId, err = c.client.GetCurrentChannelPts(context.Background(), req)
	if err == nil {
		pts = ptsId.Pts
	}
	return
}

func (c *syncClient) GetUpdateListByGtPts(userId, pts int32) (updateList []*mtproto.Update, err error) {
	req := &mtproto.UserGtPtsUpdatesRequest{
		UserId: userId,
		Pts:    pts,
	}

	var updates *mtproto.Updates
	updates, err = c.client.GetUserGtPtsUpdatesData(context.Background(), req)
	if err == nil {
		updateList = updates.GetData2().GetUpdates()
	}
	return
}

func (c *syncClient) GetChannelUpdateListByGtPts(channelId, pts int32) (updateList []*mtproto.Update, err error) {
	req := &mtproto.ChannelGtPtsUpdatesRequest{
		ChannelId: channelId,
		Pts:       pts,
	}

	var updates *mtproto.Updates
	updates, err = c.client.GetChannelGtPtsUpdatesData(context.Background(), req)
	if err == nil {
		updateList = updates.GetData2().GetUpdates()
	}
	return
}

func (c *syncClient) GetServerUpdatesState(authKeyId int64, userId int32) (state *mtproto.TLUpdatesState, err error) {
	req := &mtproto.UpdatesStateRequest{
		AuthKeyId: authKeyId,
		UserId:    userId,
	}

	var state2 *mtproto.Updates_State
	state2, err = c.client.GetServerUpdatesState(context.Background(), req)
	if err == nil {
		state = state2.To_UpdatesState()
	}
	return
}

func (c *syncClient) UpdateAuthStateSeq(authKeyId int64, pts, qts int32) (err error) {
	req := &mtproto.UpdatesStateRequest{
		AuthKeyId: authKeyId,
		Pts:       pts,
		Qts:       qts,
	}
	_, err = c.client.UpdateUpdatesState(context.Background(), req)
	return
}
*/
