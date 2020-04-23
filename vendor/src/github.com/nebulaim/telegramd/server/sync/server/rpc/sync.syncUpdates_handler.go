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

package rpc

import (
    "github.com/golang/glog"
    "github.com/nebulaim/telegramd/proto/mtproto"
    "golang.org/x/net/context"
    "github.com/nebulaim/telegramd/baselib/logger"
)

// sync.syncUpdates#3a077679 flags:# layer:int user_id:int auth_key_id:long server_id:flags.1?int not_me:flags.0?true updates:Updates = Bool;
func (s *SyncServiceImpl) SyncSyncUpdates(ctx context.Context, request *mtproto.TLSyncSyncUpdates) (*mtproto.Bool, error) {
    glog.Infof("sync.syncUpdates#3a077679 - request: {%s}", logger.JsonDebugData(request))

    err := s.processUpdatesRequest(request.GetUserId(), request.GetUpdates())
    if err == nil {
		pushData := &mtproto.PushData{
			Constructor: mtproto.TLConstructor_CRC32_sync_pushUpdatesData,
			Data2:       &mtproto.PushData_Data{AuthKeyId: request.GetAuthKeyId(), Updates: request.GetUpdates()},
		}

		if request.GetServerId() == 0 {
			s.pushUpdatesToSession(syncTypeUserNotMe, request.GetUserId(), pushData, 0)
		} else {
			s.pushUpdatesToSession(syncTypeUserMe, request.GetUserId(), pushData, request.GetServerId())
		}
    } else {
        glog.Error(err)
        return mtproto.ToBool(false), nil
    }

	glog.Infof("sync.syncUpdates#3a077679 - reply: {true}",)
	return mtproto.ToBool(true), nil
}
