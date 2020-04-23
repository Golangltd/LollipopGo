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

package contact

import (
	"fmt"
	"github.com/nebulaim/telegramd/service/contact/proto"
)

type ContactFacade interface {
	Initialize(config string) error
	ImportContacts(selfUserId int32, contacts []*contact.InputContactData) []*contact.ImportedContactData
	DeleteContact(selfUserId, contactUserId int32) *contact.DeleteResult
	DeleteContacts(selfUserId int32, contactUserIdList []int32) []*contact.DeleteResult
	BlockUser(selfUserId, id int32) bool
	UnBlockUser(selfUserId, id int32) bool
	CheckContactAndMutual(selfUserId, id int32) (bool, bool)
}

type Instance func() ContactFacade

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

func NewContactFacade(name, config string) (inst ContactFacade, err error) {
	instanceFunc, ok := instances[name]
	if !ok {
		err = fmt.Errorf("unknown instance name %q (forgot to import?)", name)
		return
	}
	inst = instanceFunc()
	err = inst.Initialize(config)
	if err != nil {
		inst = nil
	}
	return
}
