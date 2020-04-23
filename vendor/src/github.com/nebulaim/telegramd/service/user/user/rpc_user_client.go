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

package user_client

import (
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/grpc_util/service_discovery"
	"github.com/nebulaim/telegramd/proto/mtproto"
)

type rpcUserClient struct {
	// client mtproto.RPCNbfsClient
}

func rpcUserClientInstance() UserFacade {
	return &rpcUserClient{}
}

func NewRpcUserClient(discovery *service_discovery.ServiceDiscoveryClientConfig) *rpcUserClient {
	_, err := grpc_util.NewRPCClientByServiceDiscovery(discovery)

	if err != nil {
		glog.Error(err)
		panic(err)
	}

	return &rpcUserClient{}
}

//func InstallRpcUserClient(discovery *service_discovery.ServiceDiscoveryClientConfig) {
//	//conn, err := grpc_util.NewRPCClientByServiceDiscovery(discovery)
//	//
//	//if err != nil {
//	//	glog.Error(err)
//	//	panic(err)
//	//}
//	//
//	//nbfsInstance.client = mtproto.NewRPCNbfsClient(conn)
//}

func (c *rpcUserClient) Initialize(config string) error {
	var err error
	return err
}

func (c *rpcUserClient) GetUser(id int32) (*mtproto.User, error) {
	return nil, nil
}

func (c *rpcUserClient) GetUserList(idList []int32) ([]*mtproto.User, error) {
	return nil, nil
}

func (c *rpcUserClient) GetUserByPhoneNumber(phone string) (*mtproto.User, error) {
	return nil, nil
}

//func (c *rpcUserClient) GetSelfUserByPhoneNumber(phone string) (*mtproto.User, error) {
//	return nil, nil
//}
//
//GetUser(id int32) (*mtproto.User, error)
//GetUserList(idList []int32) ([]*mtproto.User, error)
//GetUserByPhoneNumber(phone string) (*mtproto.User, error)
