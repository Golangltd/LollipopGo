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
	"github.com/BurntSushi/toml"
	"github.com/nebulaim/telegramd/biz/core"
	model2 "github.com/nebulaim/telegramd/server/biz_server/langpack/model"
)

const (
	LANG_PACK_EN_FILE = "./lang_pack_en.toml"
)

var langs model2.LangPacks

func init() {
	if _, err := toml.DecodeFile(LANG_PACK_EN_FILE, &langs); err != nil {
		panic(err)
	}
	// fmt.Print(langs)
}

type LangpackServiceImpl struct {
}

func NewLangpackServiceImpl(models []core.CoreModel) *LangpackServiceImpl {
	impl := &LangpackServiceImpl{}

	for _, m := range models {
		switch m.(type) {
		}
	}

	return impl
}
