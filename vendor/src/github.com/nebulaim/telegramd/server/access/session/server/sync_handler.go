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
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/net2"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/proto/zproto"
)

func init() {
	proto.RegisterType((*mtproto.TLSyncConnectToSessionServer)(nil), "mtproto.TLSyncConnectToSessionServer")
	proto.RegisterType((*mtproto.TLSyncSessionServerConnected)(nil), "mtproto.TLSyncSessionServerConnected")
	proto.RegisterType((*mtproto.TLSyncPushUpdatesData)(nil), "mtproto.TLSyncPushUpdatesData")
	proto.RegisterType((*mtproto.TLSyncPushRpcResultData)(nil), "mtproto.TLSyncPushRpcResultData")
	proto.RegisterType((*mtproto.PushData)(nil), "mtproto.PushData")
	proto.RegisterType((*mtproto.Bool)(nil), "mtproto.Bool")
}

type syncHandler struct {
	smgr *sessionManager
}

func newSyncHandler(smgr *sessionManager) *syncHandler {
	s := &syncHandler{
		smgr: smgr,
	}
	return s
}

func protoToSyncData(m proto.Message) (*zproto.ZProtoSyncData, error) {
	x := mtproto.NewEncodeBuf(128)
	// x.UInt(mtproto.SYNC_DATA)
	n := proto.MessageName(m)
	x.Int(int32(len(n)))
	x.Bytes([]byte(n))
	b, err := proto.Marshal(m)
	x.Bytes(b)
	return &zproto.ZProtoSyncData{SyncRawData: x.GetBuf()}, err
}

func (s *syncHandler) onSyncData(conn *net2.TcpConnection, syncData *zproto.ZProtoSyncData) (*zproto.ZProtoSyncData, error) {
	dbuf := mtproto.NewDecodeBuf(syncData.SyncRawData)
	len2 := int(dbuf.Int())
	messageName := string(dbuf.Bytes(len2))
	message, err := grpc_util.NewMessageByName(messageName)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	err = proto.Unmarshal(syncData.SyncRawData[4+len2:], message)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	switch message.(type) {
	case *mtproto.TLSyncConnectToSessionServer:
		glog.Infof("onSyncData - request(ConnectToSessionServerReq): {%v}", message)
		return protoToSyncData(&mtproto.TLSyncSessionServerConnected{Data2: &mtproto.ServerConnected_Data{
			SessionServerId: getServerID(),
			ServerName:      "session",
		}})
	case *mtproto.PushData:
		glog.Infof("onSyncData - request(PushUpdatesData): {%v}", message)

		// TODO(@benqi): dispatch to session_client
		pushData, _ := message.(*mtproto.PushData)
		// dbuf := mtproto.NewDecodeBuf(pushData.GetUpdatesData())
		mdata := &messageData{
			confirmFlag:  true,
			compressFlag: false,
			// obj:          pushData.GetUpdates(),
		}
		md := &zproto.ZProtoMetadata{}

		switch pushData.GetConstructor() {
		case mtproto.TLConstructor_CRC32_sync_pushUpdatesData:
			mdata.obj = pushData.Data2.Updates
			if mdata.obj == nil {
				glog.Errorf("onSyncData - recv invalid pushData: {%v}", message)
				return protoToSyncData(mtproto.ToBool(false))
			}
			// push
			// s.smgr.pushToSessionData(pushData.GetAuthKeyId(), pushData.GetSessionId(), md, mdata)
			s.smgr.onSyncData(pushData.Data2.AuthKeyId, md, mdata)
		case mtproto.TLConstructor_CRC32_sync_pushRpcResultData:
			dbuf := mtproto.NewDecodeBuf(pushData.Data2.Result)
			result := dbuf.Object()
			glog.Info(result)
			if result == nil {
				glog.Errorf("onSyncData - recv invalid pushData: {%v}", message)
				return protoToSyncData(mtproto.ToBool(false))
			}
			mdata.obj = &mtproto.TLRpcResult{
				ReqMsgId: pushData.Data2.ClientReqMsgId,
				Result:   result,
			}
			// push
			// s.smgr.pushToSessionData(pushData.GetAuthKeyId(), pushData.GetSessionId(), md, mdata)
			s.smgr.onSyncData(pushData.Data2.AuthKeyId, md, mdata)
		default:
			glog.Errorf("invalid pushData - ", pushData)
			return protoToSyncData(mtproto.ToBool(false))
		}

		return protoToSyncData(mtproto.ToBool(true))
	default:
		err := fmt.Errorf("invalid register proto type: {%v}", message)
		glog.Error(err)
		return nil, err
	}
}
