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

// help.getAppUpdate#c812ac7e device_model:string system_version:string app_version:string lang_code:string = help.AppUpdate;
func (s *HelpServiceImpl) HelpGetAppUpdateLayer46(ctx context.Context, request *mtproto.TLHelpGetAppUpdateLayer46) (*mtproto.Help_AppUpdate, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("help.getAppUpdate#c812ac7e - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// TODO(@benqi): Impl HelpGetAppUpdate logic
	reply := &mtproto.TLHelpNoAppUpdate{Data2: &mtproto.Help_AppUpdate_Data{}}

	glog.Infof("help.getAppUpdate#c812ac7e - reply: %s\n", logger.JsonDebugData(reply))
	return reply.To_Help_AppUpdate(), nil
}
