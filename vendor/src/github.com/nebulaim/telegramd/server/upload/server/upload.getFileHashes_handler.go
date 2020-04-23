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
    "github.com/golang/glog"
    "github.com/nebulaim/telegramd/baselib/grpc_util"
    "github.com/nebulaim/telegramd/baselib/logger"
    "github.com/nebulaim/telegramd/proto/mtproto"
    "golang.org/x/net/context"
)

// upload.getFileHashes#c7025931 location:InputFileLocation offset:int = Vector<FileHash>;
func (s *UploadServiceImpl) UploadGetFileHashes(ctx context.Context, request *mtproto.TLUploadGetFileHashes) (*mtproto.Vector_FileHash, error) {
    md := grpc_util.RpcMetadataFromIncoming(ctx)
    glog.Infof("UploadGetFileHashes - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

    // TODO(@benqi): Impl UploadGetFileHashes logic

    return nil, fmt.Errorf("Not impl UploadGetFileHashes")
}
