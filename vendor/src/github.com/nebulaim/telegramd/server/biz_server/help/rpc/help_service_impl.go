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
	"github.com/nebulaim/telegramd/biz/core"
	// model2 "github.com/nebulaim/telegramd/server/biz_server/help/model"
	"encoding/json"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"io/ioutil"
)

const (
	CONFIG_FILE = "./config.json"

	// date = 1509066502,    2017/10/27 09:08:22
	// expires = 1509070295, 2017/10/27 10:11:35
	EXPIRES_TIMEOUT = 3600 // 超时时间设置为3600秒

	// support user: @benqi
	SUPPORT_USER_ID = 2
)

var config mtproto.TLConfig

func init() {
	configData, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		panic(err)
		return
	}

	err = json.Unmarshal([]byte(configData), &config)
	if err != nil {
		panic(err)
		return
	}
}

type HelpServiceImpl struct {
}

func NewHelpServiceImpl(models []core.CoreModel) *HelpServiceImpl {
	impl := &HelpServiceImpl{}

	for _, m := range models {
		switch m.(type) {
		}
	}

	return impl
}
