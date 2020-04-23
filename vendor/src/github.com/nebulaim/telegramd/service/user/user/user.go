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
	"fmt"
	"github.com/nebulaim/telegramd/proto/mtproto"
)

type UserFacade interface {
	Initialize(config string) error
	GetUser(id int32) (*mtproto.User, error)
	GetUserList(idList []int32) ([]*mtproto.User, error)
	GetUserByPhoneNumber(phone string) (*mtproto.User, error)

	CheckPhoneNumberExist(phoneNumber string) bool
	CheckBannedByPhoneNumber(phoneNumber string) bool
}

type Instance func() UserFacade

var instances = make(map[string]Instance)

func Register(name string, inst Instance) {
	if inst == nil {
		panic("register instance is nil")
	}
	if _, ok := instances[name]; ok {
		panic("register called twice for instance " + name)
	}
	instances[name] = inst
}

func NewUserFacade(name, config string) (inst UserFacade, err error) {
	instanceFunc, ok := instances[name]
	if !ok {
		err = fmt.Errorf("unknown adapter name %q (forgot to import?)", name)
		return
	}
	inst = instanceFunc()
	err = inst.Initialize(config)
	if err != nil {
		inst = nil
	}
	return
}
