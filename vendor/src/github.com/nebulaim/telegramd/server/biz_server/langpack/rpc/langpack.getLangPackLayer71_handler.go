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
    "github.com/BurntSushi/toml"
)

// langpack.getLangPack#9ab5c58e lang_code:string = LangPackDifference;
func (s *LangpackServiceImpl) LangpackGetLangPackLayer71(ctx context.Context, request *mtproto.TLLangpackGetLangPackLayer71) (*mtproto.LangPackDifference, error) {
    md := grpc_util.RpcMetadataFromIncoming(ctx)
    glog.Infof("LangpackGetLangPack - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

    if _, err := toml.DecodeFile(LANG_PACK_EN_FILE, &langs); err != nil {
        glog.Errorf("LangpackGetLangPack - decode file %s error: %v", LANG_PACK_EN_FILE, err)
        return nil, err
    }

    diff := mtproto.NewTLLangPackDifference()
    diff.SetLangCode(request.LangCode)
    diff.SetVersion(langs.Version)
    diff.SetFromVersion(0)

    diffStrings := make([]*mtproto.LangPackString, 0)
    for _, strings := range langs.Strings {
        diffStrings = append(diffStrings, &mtproto.LangPackString{
            Constructor: mtproto.TLConstructor_CRC32_langPackString,
            Data2:       strings,
        })
    }

    for _, stringPluralizeds := range langs.StringPluralizeds {
        diffStrings = append(diffStrings, &mtproto.LangPackString{
            Constructor: mtproto.TLConstructor_CRC32_langPackStringPluralized,
            Data2:       stringPluralizeds,
        })
    }

    // reply := mtproto.MakeLangPackDifference(diff)
    glog.Infof("LangpackGetLangPack - reply: %s", logger.JsonDebugData(diff))
    return diff.To_LangPackDifference(), nil
}
