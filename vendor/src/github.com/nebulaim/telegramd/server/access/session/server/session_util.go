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

package server

import (
	"github.com/nebulaim/telegramd/baselib/app"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/proto/zproto"
)

//func sendDataByConnection(conn *net2.TcpConnection, sessionID uint64, md *zproto.ZProtoMetadata, buf []byte) error {
//	smsg := &zproto.ZProtoSessionData{
//		SessionId: sessionID,
//		MtpRawData: buf,
//	}
//	//zmsg := &mtproto.ZProtoMessage{
//	//	SessionId: sessionID,
//	//	Metadata:  md,
//	//	SeqNum:    2,
//	//	Message: &mtproto.ZProtoRawPayload{
//	//		Payload: smsg.Encode(),
//	//	},
//	//}
//	return zproto.SendMessageByConn(conn, md, smsg)
//	// conn.Send(zmsg)
//}

func sendSessionDataByConnID(connID uint64, md *zproto.ZProtoMetadata, sessData *zproto.ZProtoSessionData) error {
	return app.GAppInstance.(*SessionServer).server.SendMessageByConnID(connID, md, sessData)
}

func getBizRPCClient() (*grpc_util.RPCClient, error) {
	return app.GAppInstance.(*SessionServer).bizRpcClient, nil
}

func getNbfsRPCClient() (*grpc_util.RPCClient, error) {
	return app.GAppInstance.(*SessionServer).nbfsRpcClient, nil
}

func getSyncRPCClient() (mtproto.RPCSyncClient, error) {
	return app.GAppInstance.(*SessionServer).syncRpcClient, nil
}

func getAuthSessionRPCClient() (mtproto.RPCSessionClient, error) {
	return app.GAppInstance.(*SessionServer).authSessionRpcClient, nil
}

func deleteClientSessionManager(authKeyID int64) {
	app.GAppInstance.(*SessionServer).sessionManager.onCloseSessionClientManager(authKeyID)
}

func getServerID() int32 {
	return Conf.ServerId
}

func getUUID() int64 {
	uuid, _ := app.GAppInstance.(*SessionServer).idgen.GetUUID()
	return uuid
}

func setOnline(userId int32, authKeyId int64, serverId, layer int32) {
	app.GAppInstance.(*SessionServer).status.SetSessionOnline(userId, authKeyId, serverId, layer)
}
