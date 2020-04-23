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
    "fmt"
    "github.com/nebulaim/telegramd/baselib/grpc_util"
    "github.com/nebulaim/telegramd/baselib/logger"
)

// updates.getChannelDifference#bb32d7c0 channel:InputChannel filter:ChannelMessagesFilter pts:int limit:int = updates.ChannelDifference;
func (s *UpdatesServiceImpl) UpdatesGetChannelDifferenceLayer46(ctx context.Context, request *mtproto.TLUpdatesGetChannelDifferenceLayer46) (*mtproto.Updates_ChannelDifference, error) {
    md := grpc_util.RpcMetadataFromIncoming(ctx)
    glog.Infof("UpdatesGetChannelDifferenceLayer46 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

    // TODO(@benqi): Impl UpdatesGetChannelDifferenceLayer46 logic

    return nil, fmt.Errorf("Not impl UpdatesGetChannelDifferenceLayer46")
}
