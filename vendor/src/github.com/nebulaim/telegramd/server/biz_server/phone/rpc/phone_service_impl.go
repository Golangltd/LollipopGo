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
	"github.com/nebulaim/telegramd/biz/core/message"
	"github.com/nebulaim/telegramd/biz/core/phone_call"
	"github.com/nebulaim/telegramd/biz/core/user"
	"github.com/nebulaim/telegramd/proto/mtproto"
)

// Before a voice call is ready, some preliminary actions have to be performed.
// The calling party needs to contact the party to be called and check whether it is ready to accept the call.
// Besides that, the parties have to negotiate the protocols to be used,
// learn the IP addresses of each other or of the Telegram relay servers to be used (so-called reflectors),
// and generate a one-time encryption key for this voice call with the aid of Diffie—Hellman key exchange.
// All of this is accomplished in parallel with the aid of several Telegram API methods and related notifications.
//

var (
	fingerprint uint64 = 12240908862933197005
)

const (
	PHONE_STATE_UNKNOWN = iota
	PHONE_STATE_REQUEST_CALL
)

type phoneCallState int

type phoneCallSession struct {
	id                    int64
	adminId               int32
	adminAccessHash       int64
	participantId         int32
	participantAccessHash int64
	date                  int32
	state                 int // phoneCallstate
	protocol              *mtproto.TLPhoneCallProtocol
	g_b                   []byte // acceptCall
	g_a                   []byte // confirm
}

// TODO(@benqi): 存储到redis里
var phoneCallSessionManager = make(map[int64]*phoneCallSession)

type PhoneServiceImpl struct {
	*user.UserModel
	*phone_call.PhoneCallModel
	*message.MessageModel
	RelayIp string
}

func NewPhoneServiceImpl(models []core.CoreModel, relayIp string) *PhoneServiceImpl {
	impl := &PhoneServiceImpl{RelayIp: relayIp}

	for _, m := range models {
		switch m.(type) {
		case *phone_call.PhoneCallModel:
			impl.PhoneCallModel = m.(*phone_call.PhoneCallModel)
		case *user.UserModel:
			impl.UserModel = m.(*user.UserModel)
		case *message.MessageModel:
			impl.MessageModel = m.(*message.MessageModel)
		}
	}

	return impl
}
