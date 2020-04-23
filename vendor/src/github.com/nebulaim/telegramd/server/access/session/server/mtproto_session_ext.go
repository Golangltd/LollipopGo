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

package server

import (
	"fmt"
	"github.com/nebulaim/telegramd/proto/mtproto"
)

// invokeAfterMsg#cb9f372d {X:Type} msg_id:long query:!X = X;
type TLInvokeAfterMsgExt struct {
	MsgId int64
	Query mtproto.TLObject
}

func NewInvokeAfterMsgExt(invokeAfterMsg *mtproto.TLInvokeAfterMsg) *TLInvokeAfterMsgExt {
	dbuf := mtproto.NewDecodeBuf(invokeAfterMsg.Query)
	query := dbuf.Object()

	return &TLInvokeAfterMsgExt{
		MsgId: invokeAfterMsg.MsgId,
		Query: query,
	}
}

func (m *TLInvokeAfterMsgExt) Encode() []byte {
	return nil
}

func (m *TLInvokeAfterMsgExt) EncodeToLayer(int) []byte {
	return m.Encode()
}

func (m *TLInvokeAfterMsgExt) Decode(dbuf *mtproto.DecodeBuf) error {
	return nil
}

func (m *TLInvokeAfterMsgExt) String() string {
	return fmt.Sprintf("{msg_id: %d, query: {%v}}", m.MsgId, m.Query)
}

// invokeAfterMsgs#3dc4b4f0 {X:Type} msg_ids:Vector<long> query:!X = X;
type TLInvokeAfterMsgsExt struct {
	MsgIds []int64
	Query  mtproto.TLObject
}

func NewInvokeAfterMsgsExt(invokeAfterMsgs *mtproto.TLInvokeAfterMsgs) *TLInvokeAfterMsgsExt {
	dbuf := mtproto.NewDecodeBuf(invokeAfterMsgs.Query)
	query := dbuf.Object()

	return &TLInvokeAfterMsgsExt{
		MsgIds: invokeAfterMsgs.MsgIds,
		Query:  query,
	}
}

func (m *TLInvokeAfterMsgsExt) Encode() []byte {
	return nil
}

func (m *TLInvokeAfterMsgsExt) EncodeToLayer(int) []byte {
	return m.Encode()
}

func (m *TLInvokeAfterMsgsExt) Decode(dbuf *mtproto.DecodeBuf) error {
	return nil
}

func (m *TLInvokeAfterMsgsExt) String() string {
	return fmt.Sprintf("{msg_ids: {%v}, query: {%v}}", m.MsgIds, m.Query)
}

// initConnection#c7481da6 {X:Type} api_id:int device_model:string system_version:string app_version:string system_lang_code:string lang_pack:string lang_code:string query:!X = X;
type TLInitConnectionExt struct {
	ApiId          int32
	DeviceMode     string
	SystemVersion  string
	AppVersion     string
	SystemLangCode string
	LangPack       string
	LangCode       string
	Query          mtproto.TLObject
}

func NewInitConnectionExt(initConnection *mtproto.TLInitConnection) *TLInitConnectionExt {
	dbuf := mtproto.NewDecodeBuf(initConnection.Query)
	query := dbuf.Object()

	return &TLInitConnectionExt{
		ApiId:          initConnection.ApiId,
		DeviceMode:     initConnection.DeviceModel,
		SystemVersion:  initConnection.SystemVersion,
		AppVersion:     initConnection.AppVersion,
		SystemLangCode: initConnection.SystemLangCode,
		LangCode:       initConnection.LangCode,
		LangPack:       initConnection.LangPack,
		Query:          query,
	}
}

func (m *TLInitConnectionExt) Encode() []byte {
	return nil
}

func (m *TLInitConnectionExt) Decode(dbuf *mtproto.DecodeBuf) error {
	return nil
}

func (m *TLInitConnectionExt) EncodeToLayer(int) []byte {
	return m.Encode()
}

func (m *TLInitConnectionExt) String() string {
	return fmt.Sprintf("{api_id: %d, device_mode: %s, system_version: %s, app_version: %s, system_lang_code: %s, lang_pack: %s, lang_code: %s, query: {%v}}",
		m.ApiId, m.DeviceMode, m.SystemVersion, m.AppVersion, m.SystemLangCode, m.LangCode, m.LangPack, m.Query)
}

// invokeWithLayer#da9b0d0d {X:Type} layer:int query:!X = X;
type TLInvokeWithLayerExt struct {
	Layer int32
	Query mtproto.TLObject
}

func NewInvokeWithLayerExt(invokeWithLayer *mtproto.TLInvokeWithLayer) *TLInvokeWithLayerExt {
	dbuf := mtproto.NewDecodeBuf(invokeWithLayer.Query)
	query := dbuf.Object()

	return &TLInvokeWithLayerExt{
		Query: query,
	}
}

func (m *TLInvokeWithLayerExt) Encode() []byte {
	return nil
}

func (m *TLInvokeWithLayerExt) EncodeToLayer(int) []byte {
	return m.Encode()
}

func (m *TLInvokeWithLayerExt) Decode(dbuf *mtproto.DecodeBuf) error {
	return nil
}

func (m *TLInvokeWithLayerExt) String() string {
	return ""
}

// invokeWithoutUpdates#bf9459b7 {X:Type} query:!X = X;
type TLInvokeWithoutUpdatesExt struct {
	Query mtproto.TLObject
}

func NewInvokeWithoutUpdatesExt(invokeWithoutUpdates *mtproto.TLInvokeWithoutUpdates) *TLInvokeWithoutUpdatesExt {
	dbuf := mtproto.NewDecodeBuf(invokeWithoutUpdates.Query)
	query := dbuf.Object()

	return &TLInvokeWithoutUpdatesExt{
		Query: query,
	}
}

func (m *TLInvokeWithoutUpdatesExt) Encode() []byte {
	return nil
}

func (m *TLInvokeWithoutUpdatesExt) EncodeToLayer(int) []byte {
	return m.Encode()
}

func (m *TLInvokeWithoutUpdatesExt) Decode(dbuf *mtproto.DecodeBuf) error {
	return nil
}

func (m *TLInvokeWithoutUpdatesExt) String() string {
	return fmt.Sprintf("{query: %v}", m.Query)
}
