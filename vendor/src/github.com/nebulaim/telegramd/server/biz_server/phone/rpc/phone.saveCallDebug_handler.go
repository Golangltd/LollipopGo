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
	"fmt"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"golang.org/x/net/context"
)

// phone.saveCallDebug#277add7e peer:InputPhoneCall debug:DataJSON = Bool;
func (s *PhoneServiceImpl) PhoneSaveCallDebug(ctx context.Context, request *mtproto.TLPhoneSaveCallDebug) (*mtproto.Bool, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("phone.saveCallDebug#277add7e - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	//// TODO(@benqi): check peer
	peer := request.GetPeer().To_InputPhoneCall()

	callSession, err := s.PhoneCallModel.MakePhoneCallLogcByLoad(peer.GetId())
	if err != nil {
		glog.Errorf("invalid peer: {%v}, err: %v", peer, err)
		return nil, err
	}

	if md.UserId == callSession.AdminId {
		if peer.GetAccessHash() != callSession.AdminAccessHash {
			err = fmt.Errorf("invalid peer: {%v}", peer)
			glog.Errorf("invalid peer: {%v}", peer)
			return nil, err
		}

		callSession.SetAdminDebugData(request.GetDebug().GetData2().GetData())
	} else {
		if peer.GetAccessHash() != callSession.ParticipantAccessHash {
			err = fmt.Errorf("invalid peer: {%v}", peer)
			glog.Errorf("invalid peer: {%v}", peer)
			return nil, err
		}

		callSession.SetParticipantDebugData(request.GetDebug().GetData2().GetData())
	}

	glog.Infof("phone.saveCallDebug#277add7e - reply: {true}")
	return mtproto.ToBool(true), nil
}
