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

package core

import (
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/base"
	"github.com/nebulaim/telegramd/service/idgen/client"
)

const (
	// TODO(@benqi): 使用更紧凑的前缀
	messageDataNgenId       = "message_data_ngen_"
	boxUpdatesNgenId        = "message_box_ngen_"
	channelBoxUpdatesNgenId = "channel_message_box_ngen_"
	seqUpdatesNgenId        = "seq_updates_ngen_"
	ptsUpdatesNgenId        = "pts_updates_ngen_"
	qtsUpdatesNgenId        = "qts_updates_ngen_"
	channelPtsUpdatesNgenId = "channel_pts_updates_ngen_"
)

var seqIDGen idgen.SeqIDGen

func initSeqIDGen(redisName string) {
	var err error
	seqIDGen, err = idgen.NewSeqIDGen("redis", redisName)
	if err != nil {
		glog.Fatal("seqidGen init error: ", err)
	}
}

func NextMessageBoxId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetNextSeqID(boxUpdatesNgenId + base.Int32ToString(key))
	return
}

func CurrentMessageBoxId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetCurrentSeqID(boxUpdatesNgenId + base.Int32ToString(key))
	return
}

func NextMessageDataId(key int64) (seq int64) {
	seq, _ = seqIDGen.GetNextSeqID(messageDataNgenId + base.Int64ToString(key))
	return
}

func NextChannelMessageBoxId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetNextSeqID(channelBoxUpdatesNgenId + base.Int32ToString(key))
	return
}

func CurrentChannelMessageBoxId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetCurrentSeqID(channelBoxUpdatesNgenId + base.Int32ToString(key))
	return
}

func NextSeqId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetNextSeqID(seqUpdatesNgenId + base.Int32ToString(key))
	return
}

func CurrentSeqId(key int32) (seq int64) {
	var err error
	seq, _ = seqIDGen.GetCurrentSeqID(seqUpdatesNgenId + base.Int32ToString(key))

	if err != nil {
		seq = -1
	}
	return
}

func NextPtsId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetNextSeqID(ptsUpdatesNgenId + base.Int32ToString(key))
	return
}

func NextNPtsId(key int32, n int) (seq int64) {
	seq, _ = seqIDGen.GetNextNSeqID(ptsUpdatesNgenId + base.Int32ToString(key), n)
	return
}

func CurrentPtsId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetCurrentSeqID(ptsUpdatesNgenId + base.Int32ToString(key))
	return
}

func NextQtsId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetNextSeqID(qtsUpdatesNgenId + base.Int32ToString(key))
	return
}

func CurrentQtsId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetNextSeqID(qtsUpdatesNgenId + base.Int32ToString(key))
	return
}

func NextChannelPtsId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetNextSeqID(channelPtsUpdatesNgenId + base.Int32ToString(key))
	return
}

func NextChannelNPtsId(key int32, n int) (seq int64) {
	seq, _ = seqIDGen.GetNextNSeqID(channelPtsUpdatesNgenId + base.Int32ToString(key), n)
	return
}

func CurrentChannelPtsId(key int32) (seq int64) {
	seq, _ = seqIDGen.GetNextSeqID(channelPtsUpdatesNgenId + base.Int32ToString(key))
	return
}
