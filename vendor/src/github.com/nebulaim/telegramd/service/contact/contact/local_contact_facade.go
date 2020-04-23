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
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/service/contact/biz/core"
)

type localContactFacade struct {
	*core.ContactModel
}

func localContactFacadeInstance() ContactFacade {
	return &localContactFacade{}
}

func NewLocalUserApi(dbName string) (*localContactFacade, error) {
	var err error

	facade := &localContactFacade{}
	facade.ContactModel, err = core.InitContactModel(dbName)

	return facade, err
}

func (c *localContactFacade) Initialize(config string) error {
	glog.Info("localUserApi - Initialize config: ", config)

	var err error

	dbName := config
	c.ContactModel, err = core.InitContactModel(dbName)

	return err
}
