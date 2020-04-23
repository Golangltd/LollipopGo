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

package server

import (
	"github.com/nebulaim/telegramd/baselib/base"
	"github.com/nebulaim/telegramd/service/nbfs/cachefs"
	"github.com/nebulaim/telegramd/service/nbfs/nbfs"
)

type UploadServiceImpl struct {
	nbfs_client.NbfsFacade
}

func NewUploadServiceImpl(serverId int32, dataPath string) *UploadServiceImpl {
	if dataPath == "" {
		dataPath = "/opt/nbfs"
	}

	cachefs.InitCacheFS(dataPath)

	s := &UploadServiceImpl{}
	s.NbfsFacade, _ = nbfs_client.NewNbfsFacade("local", base.Int32ToString(serverId))
	return s
}
