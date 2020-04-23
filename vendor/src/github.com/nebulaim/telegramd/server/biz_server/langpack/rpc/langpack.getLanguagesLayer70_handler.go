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

package rpc

import (
    "github.com/golang/glog"
    "github.com/nebulaim/telegramd/proto/mtproto"
    "golang.org/x/net/context"
    "github.com/nebulaim/telegramd/baselib/grpc_util"
    "github.com/nebulaim/telegramd/baselib/logger"
)

// langpack.getLanguages#800fd57d = Vector<LangPackLanguage>;
func (s *LangpackServiceImpl) LangpackGetLanguagesLayer70(ctx context.Context, request *mtproto.TLLangpackGetLanguagesLayer70) (*mtproto.Vector_LangPackLanguage, error) {
    md := grpc_util.RpcMetadataFromIncoming(ctx)
    glog.Infof("langpack.getLanguages#800fd57d - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

    // TODO(@benqi): Add other language
    language := &mtproto.TLLangPackLanguage{Data2: &mtproto.LangPackLanguage_Data{
        Name:       "English",
        NativeName: "English",
        LangCode:   "en",
    }}

    languages := &mtproto.Vector_LangPackLanguage{}
    languages.Datas = append(languages.Datas, language.To_LangPackLanguage())

    glog.Infof("langpack.getLanguages#800fd57d - reply: %s", logger.JsonDebugData(languages))
    return languages, nil
}
