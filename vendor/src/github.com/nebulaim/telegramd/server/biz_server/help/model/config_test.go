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

package model

import (
	"encoding/json"
	"fmt"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"io/ioutil"
	"testing"
)

func TestGetHelpConfig(t *testing.T) {
	helpConfig := mtproto.NewTLConfig()
	// data2 := &ProfilePhotoIds{}
	configData, err := ioutil.ReadFile("./config_test.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal([]byte(configData), helpConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", logger.JsonDebugData(helpConfig))
}
