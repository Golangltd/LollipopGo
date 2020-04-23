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

package phone_call

import (
	"encoding/hex"
	"fmt"
	base2 "github.com/nebulaim/telegramd/baselib/base"
	"github.com/nebulaim/telegramd/biz/core"
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"math/rand"
	"time"
)

// TODO(@benqi): Using redis storage phone_call_sessions

type phoneCallLogic struct {
	*PhoneCallSession
	dao *phoneCallsDAO
}

func (m *PhoneCallModel) NewPhoneCallLogic(adminId, participantId int32, ga []byte, protocol *mtproto.TLPhoneCallProtocol) *phoneCallLogic {
	phoneCallSession := &PhoneCallSession{
		Id:                    core.GetUUID(),
		AdminId:               adminId,
		AdminAccessHash:       rand.Int63(),
		ParticipantId:         participantId,
		ParticipantAccessHash: rand.Int63(),
		UdpP2P:                protocol.GetUdpP2P(),
		UdpReflector:          protocol.GetUdpReflector(),
		MinLayer:              protocol.GetMinLayer(),
		MaxLayer:              protocol.GetMaxLayer(),
		GA:                    ga,
		State:                 0,
		Date:                  time.Now().Unix(),
	}
	session := &phoneCallLogic{
		PhoneCallSession: phoneCallSession,
		dao:              m.dao,
	}

	do := &dataobject.PhoneCallSessionsDO{
		CallSessionId:         session.Id,
		AdminId:               session.AdminId,
		AdminAccessHash:       session.AdminAccessHash,
		ParticipantId:         session.ParticipantId,
		ParticipantAccessHash: session.ParticipantAccessHash,
		UdpP2p:                base2.BoolToInt8(session.UdpP2P),
		UdpReflector:          base2.BoolToInt8(session.UdpReflector),
		MinLayer:              session.MinLayer,
		MaxLayer:              session.MaxLayer,
		GA:                    hex.EncodeToString(session.GA),
		Date:                  int32(session.Date),
	}
	m.dao.PhoneCallSessionsDAO.Insert(do)
	return session
}

func (m *PhoneCallModel) MakePhoneCallLogcByLoad(id int64) (*phoneCallLogic, error) {
	do := m.dao.PhoneCallSessionsDAO.Select(id)
	if do == nil {
		err := fmt.Errorf("not found call session: %d", id)
		return nil, err
	}

	phoneCallSession := &PhoneCallSession{
		Id:                    do.CallSessionId,
		AdminId:               do.AdminId,
		AdminAccessHash:       do.AdminAccessHash,
		ParticipantId:         do.ParticipantId,
		ParticipantAccessHash: do.ParticipantAccessHash,
		UdpP2P:                do.UdpP2p == 1,
		UdpReflector:          do.UdpReflector == 1,
		MinLayer:              do.MinLayer,
		MaxLayer:              do.MaxLayer,
		// GA:                    do.GA,
		State: 0,
		Date:  int64(do.Date),
	}

	session := &phoneCallLogic{
		PhoneCallSession: phoneCallSession,
		dao:              m.dao,
	}

	session.GA, _ = hex.DecodeString(do.GA)
	return session, nil
}

func (p *phoneCallLogic) SetGB(gb []byte) {
	p.GB = gb
	p.dao.PhoneCallSessionsDAO.UpdateGB(hex.EncodeToString(gb), p.Id)
}

func (p *phoneCallLogic) SetAdminDebugData(dataJson string) {
	p.dao.PhoneCallSessionsDAO.UpdateAdminDebugData(dataJson, p.Id)
}

func (p *phoneCallLogic) SetParticipantDebugData(dataJson string) {
	p.dao.PhoneCallSessionsDAO.UpdateParticipantDebugData(dataJson, p.Id)
}

func (p *phoneCallLogic) toPhoneCallProtocol() *mtproto.PhoneCallProtocol {
	return &mtproto.PhoneCallProtocol{
		Constructor: mtproto.TLConstructor_CRC32_phoneCallProtocol,
		Data2: &mtproto.PhoneCallProtocol_Data{
			UdpP2P:       p.UdpP2P,
			UdpReflector: p.UdpReflector,
			MinLayer:     p.MinLayer,
			MaxLayer:     p.MaxLayer,
		},
	}
}

func (p *phoneCallLogic) ToPhoneCallProtocol() *mtproto.PhoneCallProtocol {
	return p.toPhoneCallProtocol()
}

// phoneCallRequested#83761ce4 id:long access_hash:long date:int admin_id:int participant_id:int g_a_hash:bytes protocol:PhoneCallProtocol = PhoneCall;
func (p *phoneCallLogic) ToPhoneCallRequested() *mtproto.TLPhoneCallRequested {
	return &mtproto.TLPhoneCallRequested{Data2: &mtproto.PhoneCall_Data{
		Id:            p.Id,
		AccessHash:    p.ParticipantAccessHash,
		Date:          int32(p.Date),
		AdminId:       p.AdminId,
		ParticipantId: p.ParticipantId,
		GAHash:        p.GA,
		Protocol:      p.toPhoneCallProtocol(),
	}}
}

// phoneCallWaiting#1b8f4ad1 flags:# id:long access_hash:long date:int admin_id:int participant_id:int protocol:PhoneCallProtocol receive_date:flags.0?int = PhoneCall;
func (p *phoneCallLogic) ToPhoneCallWaiting(selfId int32, receiveDate int32) *mtproto.TLPhoneCallWaiting {
	var (
		accessHash int64
	)

	if selfId == p.AdminId {
		accessHash = p.AdminAccessHash
	} else {
		accessHash = p.ParticipantAccessHash
	}

	return &mtproto.TLPhoneCallWaiting{Data2: &mtproto.PhoneCall_Data{
		Id:            p.Id,
		AccessHash:    accessHash,
		Date:          int32(p.Date),
		AdminId:       p.AdminId,
		ParticipantId: p.ParticipantId,
		GAHash:        p.GA,
		Protocol:      p.toPhoneCallProtocol(),
		ReceiveDate:   receiveDate,
	}}
}

// phoneCallAccepted#6d003d3f id:long access_hash:long date:int admin_id:int participant_id:int g_b:bytes protocol:PhoneCallProtocol = PhoneCall;
func (p *phoneCallLogic) ToPhoneCallAccepted() *mtproto.TLPhoneCallAccepted {
	return &mtproto.TLPhoneCallAccepted{Data2: &mtproto.PhoneCall_Data{
		Id:            p.Id,
		AccessHash:    p.AdminAccessHash,
		Date:          int32(p.Date),
		AdminId:       p.AdminId,
		ParticipantId: p.ParticipantId,
		GB:            p.GB,
		Protocol:      p.toPhoneCallProtocol(),
	}}
}

// phoneConnection#9d4c17c0 id:long ip:string ipv6:string port:int peer_tag:bytes = PhoneConnection;
func makeConnection(relayIp string) *mtproto.PhoneConnection {
	return &mtproto.PhoneConnection{
		Constructor: mtproto.TLConstructor_CRC32_phoneConnection,
		Data2: &mtproto.PhoneConnection_Data{
			Id: 50003,
			// Ip:      "192.168.4.32",
			Ip:      relayIp,
			Ipv6:    "",
			Port:    50001,
			PeerTag: []byte("24ffcbeb7980d28b"),
		},
	}
}

// phoneCall#ffe6ab67 id:long access_hash:long date:int admin_id:int participant_id:int g_a_or_b:bytes key_fingerprint:long protocol:PhoneCallProtocol connection:PhoneConnection alternative_connections:Vector<PhoneConnection> start_date:int = PhoneCall;
func (p *phoneCallLogic) ToPhoneCall(selfId int32, keyFingerprint int64, relayIp string) *mtproto.TLPhoneCall {
	var (
		accessHash int64
		gaOrGb     []byte
	)

	if selfId == p.AdminId {
		accessHash = p.AdminAccessHash
		gaOrGb = p.GB
	} else {
		accessHash = p.ParticipantAccessHash
		gaOrGb = p.GA
	}

	return &mtproto.TLPhoneCall{Data2: &mtproto.PhoneCall_Data{
		Id:                     p.Id,
		AccessHash:             accessHash,
		Date:                   int32(p.Date),
		AdminId:                p.AdminId,
		ParticipantId:          p.ParticipantId,
		GAOrB:                  gaOrGb,
		KeyFingerprint:         keyFingerprint,
		Protocol:               p.toPhoneCallProtocol(),
		Connection:             makeConnection(relayIp),
		AlternativeConnections: []*mtproto.PhoneConnection{}, // TODO(@benqi):
		StartDate:              0,
	}}
}
