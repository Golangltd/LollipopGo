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
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/nebulaim/telegramd/baselib/bytes2"
	"github.com/nebulaim/telegramd/baselib/crypto"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/baselib/net2"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/proto/zproto"
	"math/big"
	"time"
	"context"
)

const (
	SHA_DIGEST_LENGTH = 20
)

var (
	headerRpcMetadata = "auth_key_metadata"

	// TODO(@benqi): 预先计算出fingerprint
	// 这里直接使用了0xc3b42b026ce86b21
	fingerprint uint64 = 12240908862933197005

	// TODO(@benqi): 使用算法生成PQ
	// 这里直接指定了PQ值: {0x17, 0xED, 0x48, 0x94, 0x1A, 0x08, 0xF9, 0x81}
	pq = string([]byte{0x17, 0xED, 0x48, 0x94, 0x1A, 0x08, 0xF9, 0x81})

	// TODO(@benqi): 直接指定了p和q
	p = []byte{0x49, 0x4C, 0x55, 0x3B}
	q = []byte{0x53, 0x91, 0x10, 0x73}

	// TODO(@benqi): 直接指定了dh2048_p和dh2048_g!!!
	// andriod client 指定的good prime
	//
	// static const char *goodPrime = "
	//
	// c71caeb9c6b1c9048e6c522f
	// 70f13f73980d40238e3e21c1
	// 4934d037563d930f48198a0a
	// a7c14058229493d22530f4db
	// fa336f6e0ac925139543aed4
	// 4cce7c3720fd51f69458705a
	// c68cd4fe6b6b13abdc974651
	// 2969328454f18faf8c595f64
	// 2477fe96bb2a941d5bcd1d4a
	// c8cc49880708fa9b378e3c4f
	// 3a9060bee67cf9a4a4a69581
	// 1051907e162753b56b0f6b41
	// 0dba74d8a84b2a14b3144e0e
	// f1284754fd17ed950d5965b4
	// b9dd46582db1178d169c6bc4
	// 65b0d6ff9ca3928fef5b9ae4
	// e418fc15e83ebea0f87fa9ff
	// 5eed70050ded2849f47bf959
	// d956850ce929851f0d8115f6
	// 35b105ee2e4e15d04b2454bf
	// 6f4fadf034b10403119cd8e3
	// b92fcc5b";
	//
	dh2048_p = []byte{
		0xc7, 0x1c, 0xae, 0xb9, 0xc6, 0xb1, 0xc9, 0x04, 0x8e, 0x6c, 0x52, 0x2f,
		0x70, 0xf1, 0x3f, 0x73, 0x98, 0x0d, 0x40, 0x23, 0x8e, 0x3e, 0x21, 0xc1,
		0x49, 0x34, 0xd0, 0x37, 0x56, 0x3d, 0x93, 0x0f, 0x48, 0x19, 0x8a, 0x0a,
		0xa7, 0xc1, 0x40, 0x58, 0x22, 0x94, 0x93, 0xd2, 0x25, 0x30, 0xf4, 0xdb,
		0xfa, 0x33, 0x6f, 0x6e, 0x0a, 0xc9, 0x25, 0x13, 0x95, 0x43, 0xae, 0xd4,
		0x4c, 0xce, 0x7c, 0x37, 0x20, 0xfd, 0x51, 0xf6, 0x94, 0x58, 0x70, 0x5a,
		0xc6, 0x8c, 0xd4, 0xfe, 0x6b, 0x6b, 0x13, 0xab, 0xdc, 0x97, 0x46, 0x51,
		0x29, 0x69, 0x32, 0x84, 0x54, 0xf1, 0x8f, 0xaf, 0x8c, 0x59, 0x5f, 0x64,
		0x24, 0x77, 0xfe, 0x96, 0xbb, 0x2a, 0x94, 0x1d, 0x5b, 0xcd, 0x1d, 0x4a,
		0xc8, 0xcc, 0x49, 0x88, 0x07, 0x08, 0xfa, 0x9b, 0x37, 0x8e, 0x3c, 0x4f,
		0x3a, 0x90, 0x60, 0xbe, 0xe6, 0x7c, 0xf9, 0xa4, 0xa4, 0xa6, 0x95, 0x81,
		0x10, 0x51, 0x90, 0x7e, 0x16, 0x27, 0x53, 0xb5, 0x6b, 0x0f, 0x6b, 0x41,
		0x0d, 0xba, 0x74, 0xd8, 0xa8, 0x4b, 0x2a, 0x14, 0xb3, 0x14, 0x4e, 0x0e,
		0xf1, 0x28, 0x47, 0x54, 0xfd, 0x17, 0xed, 0x95, 0x0d, 0x59, 0x65, 0xb4,
		0xb9, 0xdd, 0x46, 0x58, 0x2d, 0xb1, 0x17, 0x8d, 0x16, 0x9c, 0x6b, 0xc4,
		0x65, 0xb0, 0xd6, 0xff, 0x9c, 0xa3, 0x92, 0x8f, 0xef, 0x5b, 0x9a, 0xe4,
		0xe4, 0x18, 0xfc, 0x15, 0xe8, 0x3e, 0xbe, 0xa0, 0xf8, 0x7f, 0xa9, 0xff,
		0x5e, 0xed, 0x70, 0x05, 0x0d, 0xed, 0x28, 0x49, 0xf4, 0x7b, 0xf9, 0x59,
		0xd9, 0x56, 0x85, 0x0c, 0xe9, 0x29, 0x85, 0x1f, 0x0d, 0x81, 0x15, 0xf6,
		0x35, 0xb1, 0x05, 0xee, 0x2e, 0x4e, 0x15, 0xd0, 0x4b, 0x24, 0x54, 0xbf,
		0x6f, 0x4f, 0xad, 0xf0, 0x34, 0xb1, 0x04, 0x03, 0x11, 0x9c, 0xd8, 0xe3,
		0xb9, 0x2f, 0xcc, 0x5b,
	}

	dh2048_g = []byte{0x02}
)

type handshake struct {
	rsa                  *crypto.RSACryptor
	dh2048p              []byte
	dh2048g              []byte
	bigIntDH2048G        *big.Int
	bigIntDH2048P        *big.Int
	authSessionRpcClient mtproto.RPCSessionClient
}

func newHandshake(c mtproto.RPCSessionClient) *handshake {
	s := &handshake{
		rsa:                  crypto.NewRSACryptor(),
		dh2048p:              dh2048_p,
		dh2048g:              dh2048_g,
		bigIntDH2048P:        new(big.Int).SetBytes(dh2048_p),
		bigIntDH2048G:        new(big.Int).SetBytes(dh2048_g),
		authSessionRpcClient: c,
	}
	return s
}

func (s *handshake) onHandshake(conn *net2.TcpConnection, hmsg *zproto.ZProtoHandshakeMessage) (*zproto.ZProtoHandshakeMessage, error) {
	var (
		state = hmsg.State
		err   error
		res   mtproto.TLObject
	)

	mtpMessage := &mtproto.UnencryptedMessage{}
	err = mtpMessage.Decode(hmsg.MTPRawData[8:])
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	switch mtpMessage.Object.(type) {
	case *mtproto.TLReqPq:
		res, err = s.onReqPq(state, mtpMessage.Object.(*mtproto.TLReqPq))
		state.State = zproto.STATE_pq_res
	case *mtproto.TLReqPqMulti:
		res, err = s.onReqPqMulti(state, mtpMessage.Object.(*mtproto.TLReqPqMulti))
		state.State = zproto.STATE_pq_res
	case *mtproto.TLReq_DHParams:
		res, err = s.onReq_DHParams(hmsg.State, mtpMessage.Object.(*mtproto.TLReq_DHParams))
		state.State = zproto.STATE_DH_params_res
	case *mtproto.TLSetClient_DHParams:
		res, err = s.onSetClient_DHParams(hmsg.State, mtpMessage.Object.(*mtproto.TLSetClient_DHParams))
		state.State = zproto.STATE_dh_gen_res
	case *mtproto.TLMsgsAck:
		// func (s *handshake) onMsgsAck(state *mtproto.HandshakeState, request *mtproto.TLMsgsAck) error {
		err = s.onMsgsAck(hmsg.State, mtpMessage.Object.(*mtproto.TLMsgsAck))
		return nil, err
	default:
		err = fmt.Errorf("invalid handshake type: %v", hmsg)
		return nil, err
	}

	if err != nil {
		state.ResState = zproto.RES_STATE_ERROR
	} else {
		state.ResState = zproto.RES_STATE_OK
	}

	mtpMessage2 := &mtproto.UnencryptedMessage{
		MessageId: mtproto.GenerateMessageId(),
		Object:    res,
	}

	zmsg := &zproto.ZProtoHandshakeMessage{
		SessionId:  hmsg.SessionId,
		State:      state,
		MTPRawData: mtpMessage2.Encode(),
	}
	return zmsg, nil
}

func toAuthKeyMetadata(ctx []byte) (md *mtproto.AuthKeyMetadata, err error) {
	md = &mtproto.AuthKeyMetadata{}
	err = proto.Unmarshal(ctx, md)
	return
}

// req_pq#60469778 nonce:int128 = ResPQ;
func (s *handshake) onReqPq(state *zproto.HandshakeState, request *mtproto.TLReqPq) (*mtproto.TLResPQ, error) {
	authKeyMD, err := toAuthKeyMetadata(state.Ctx)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	glog.Infof("onReqPq - state: %d, res_state: %d, metadata: {%v}, request: %s",
		state.State, state.ResState, authKeyMD, logger.JsonDebugData(request))

	// check State and ResState

	// 检查数据是否合法
	if request.GetNonce() == nil || len(request.GetNonce()) != 16 {
		err := fmt.Errorf("onReqPq - invalid nonce: %v", request)
		glog.Error(err)
		return nil, err
	}

	resPQ := &mtproto.TLResPQ{Data2: &mtproto.ResPQ_Data{
		Nonce:                       request.Nonce,
		ServerNonce:                 crypto.GenerateNonce(16),
		Pq:                          pq,
		ServerPublicKeyFingerprints: []int64{int64(fingerprint)},
	}}
	//
	//resPQ := mtproto.NewTLResPQ()
	//resPQ.SetNonce(make([]byte, 16))
	//copy(resPQ.Data2.Nonce, request.GetNonce())
	//resPQ.SetServerNonce(crypto.GenerateNonce(16))
	//resPQ.SetPq(pq)
	//resPQ.SetServerPublicKeyFingerprints([]int64{int64(fingerprint)})

	// 缓存客户端Nonce
	authKeyMD.Nonce = request.GetNonce()
	authKeyMD.ServerNonce = resPQ.GetServerNonce()

	state.Ctx, _ = proto.Marshal(authKeyMD)

	glog.Infof("onReqPq - metadata: {%v}, reply: %s", authKeyMD, logger.JsonDebugData(resPQ))
	return resPQ, nil
}

// req_pq#60469778 nonce:int128 = ResPQ;
func (s *handshake) onReqPqMulti(state *zproto.HandshakeState, request *mtproto.TLReqPqMulti) (*mtproto.TLResPQ, error) {
	authKeyMD, err := toAuthKeyMetadata(state.Ctx)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	glog.Infof("onReqPq - state: %d, res_state: %d, metadata: {%v}, request: %s",
		state.State, state.ResState, authKeyMD, logger.JsonDebugData(request))

	// check State and ResState

	// 检查数据是否合法
	if request.GetNonce() == nil || len(request.GetNonce()) != 16 {
		err := fmt.Errorf("onReqPq - invalid nonce: %v", request)
		glog.Error(err)
		return nil, err
	}

	resPQ := &mtproto.TLResPQ{Data2: &mtproto.ResPQ_Data{
		Nonce:                       request.Nonce,
		ServerNonce:                 crypto.GenerateNonce(16),
		Pq:                          pq,
		ServerPublicKeyFingerprints: []int64{int64(fingerprint)},
	}}
	//
	//resPQ := mtproto.NewTLResPQ()
	//resPQ.SetNonce(make([]byte, 16))
	//copy(resPQ.Data2.Nonce, request.GetNonce())
	//resPQ.SetServerNonce(crypto.GenerateNonce(16))
	//resPQ.SetPq(pq)
	//resPQ.SetServerPublicKeyFingerprints([]int64{int64(fingerprint)})

	// 缓存客户端Nonce
	authKeyMD.Nonce = request.GetNonce()
	authKeyMD.ServerNonce = resPQ.GetServerNonce()

	state.Ctx, _ = proto.Marshal(authKeyMD)

	glog.Infof("onReqPq - metadata: {%v}, reply: %s", authKeyMD, logger.JsonDebugData(resPQ))
	return resPQ, nil
}

// req_DH_params#d712e4be nonce:int128 server_nonce:int128 p:string q:string public_key_fingerprint:long encrypted_data:string = Server_DH_Params;
func (s *handshake) onReq_DHParams(state *zproto.HandshakeState, request *mtproto.TLReq_DHParams) (*mtproto.Server_DH_Params, error) {
	authKeyMD, err := toAuthKeyMetadata(state.Ctx)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	glog.Infof("onReq_DHParams - state: %d, res_state: %d, metadata: {%v}, request: %s",
		state.State, state.ResState, authKeyMD, logger.JsonDebugData(request))

	// 客户端传输数据解析
	// check Nonce
	if !bytes.Equal(request.Nonce, authKeyMD.Nonce) {
		err = fmt.Errorf("onReq_DHParams - Invalid Nonce, req: %s, back: %s",
			bytes2.HexDump(request.Nonce),
			bytes2.HexDump(authKeyMD.Nonce))
		glog.Error(err)
		return nil, err
	}

	// check ServerNonce
	if !bytes.Equal(request.ServerNonce, authKeyMD.ServerNonce) {
		err = fmt.Errorf("onReq_DHParams - Wrong ServerNonce, req: %s, back: %s",
			bytes2.HexDump(request.ServerNonce),
			bytes2.HexDump(authKeyMD.ServerNonce))
		glog.Error(err)
		return nil, err
	}

	// check P
	if !bytes.Equal([]byte(request.P), p) {
		err = fmt.Errorf("onReq_DHParams - Invalid p valuee")
		glog.Error(err)
		return nil, err
	}

	// check Q
	if !bytes.Equal([]byte(request.Q), q) {
		err = fmt.Errorf("onReq_DHParams - Invalid q value")
		glog.Error(err)
		return nil, err
	}

	if request.PublicKeyFingerprint != int64(fingerprint) {
		err = fmt.Errorf("onReq_DHParams - Invalid PublicKeyFingerprint value")
		glog.Error(err)
		return nil, err
	}

	// new_nonce := another (good) random number generated by the client;
	// after this query, it is known to both client and server;
	//
	// data := a serialization of
	//
	// p_q_inner_data#83c95aec pq:string p:string q:string nonce:int128 server_nonce:int128 new_nonce:int256 = P_Q_inner_data
	// or of
	// p_q_inner_data_temp#3c6a84d4 pq:string p:string q:string nonce:int128 server_nonce:int128 new_nonce:int256 expires_in:int = P_Q_inner_data;
	//
	// data_with_hash := SHA1(data) + data + (any random bytes);
	// 	 such that the length equal 255 bytes;
	// encrypted_data := RSA (data_with_hash, server_public_key);
	// 	 a 255-byte long number (big endian) is raised to the requisite power over the requisite modulus,
	// 	 and the result is stored as a 256-byte number.
	//

	// encryptedData := []byte(request.EncryptedData)
	//
	// glog.Info("EncryptedData: len = ", len(encryptedData), ", data: ", hex.EncodeToString(encryptedData))
	//
	// 1. 解密
	encryptedPQInnerData := s.rsa.Decrypt([]byte(request.EncryptedData))

	// TODO(@benqi): sha1_check
	// sha1Check := sha1.Sum(encryptedPQInnerData[20:])
	// glog.Info(hex.EncodeToString(sha1Check[:]))
	// glog.Info(hex.EncodeToString(encryptedPQInnerData[:20]))
	//if !bytes.Equal(sha1Check[:], encryptedPQInnerData[0:20]) {
	//	glog.Error("process Req_DHParams - sha1Check error")
	//	return nil, fmt.Errorf("process Req_DHParams - sha1Check error")
	//}

	// pqInnerData := new(mtproto.P_QInnerData)
	// pqInnerData.Decode()
	//
	// 2. 反序列化出pqInnerData
	pqInnerData := mtproto.NewTLPQInnerData()
	dbuf := mtproto.NewDecodeBuf(encryptedPQInnerData[SHA_DIGEST_LENGTH+4:])
	err = pqInnerData.Decode(dbuf)
	if err != nil {
		glog.Errorf("process Req_DHParams - TLPQInnerData decode error: %v", err)
		return nil, fmt.Errorf("process Req_DHParams - TLPQInnerData decode error: %v", err)
	}

	// 2. 再检查一遍p_q_inner_data里的pq, p, q, nonce, server_nonce合法性
	// 客户端传输数据解析
	// PQ
	if !bytes.Equal([]byte(pqInnerData.GetPq()), []byte(pq)) {
		glog.Error("process Req_DHParams - Invalid p_q_inner_data.pq value")
		return nil, fmt.Errorf("process Req_DHParams - Invalid p_q_inner_data.pq value")
	}

	// P
	if !bytes.Equal([]byte(pqInnerData.GetP()), p) {
		glog.Error("process Req_DHParams - Invalid p_q_inner_data.p value")
		return nil, fmt.Errorf("process Req_DHParams - Invalid p_q_inner_data.p value")
	}

	// Q
	if !bytes.Equal([]byte(pqInnerData.GetQ()), q) {
		glog.Error("process Req_DHParams - Invalid p_q_inner_data.q value")
		return nil, fmt.Errorf("process Req_DHParams - Invalid p_q_inner_data.q value")
	}

	// Nonce
	if !bytes.Equal(pqInnerData.GetNonce(), authKeyMD.Nonce) {
		glog.Error("process Req_DHParams - Invalid Nonce")
		return nil, fmt.Errorf("process Req_DHParams - InvalidNonce")
	}

	// ServerNonce
	if !bytes.Equal(request.GetServerNonce(), authKeyMD.ServerNonce) {
		glog.Error("process Req_DHParams - Wrong ServerNonce")
		return nil, fmt.Errorf("process Req_DHParams - Wrong ServerNonce")
	}

	// glog.Info("processReq_DHParams - pqInnerData Decode sucess: ", pqInnerData.String())

	// 检查NewNonce的长度(int256)
	// 缓存NewNonce
	authKeyMD.NewNonce = pqInnerData.GetNewNonce()
	authKeyMD.A = crypto.GenerateNonce(256)
	authKeyMD.P = s.dh2048p

	bigIntA := new(big.Int).SetBytes(authKeyMD.A)
	// bigIntP := new(big.Int).SetBytes(authKeyMD.P)

	//c.A = new(big.Int).SetBytes()
	//c.P = new(big.Int).SetBytes(dh2048_p)

	// 服务端计算GA = g^a mod p
	g_a := new(big.Int)
	g_a.Exp(s.bigIntDH2048G, bigIntA, s.bigIntDH2048P)

	// ServerNonce
	server_DHInnerData := &mtproto.TLServer_DHInnerData{Data2: &mtproto.Server_DHInnerData_Data{
		Nonce:       authKeyMD.Nonce,
		ServerNonce: authKeyMD.ServerNonce,
		G:           int32(s.dh2048g[0]),
		GA:          string(g_a.Bytes()),
		DhPrime:     string(s.dh2048p),
		ServerTime:  int32(time.Now().Unix()),
	}}

	server_DHInnerData_buf := server_DHInnerData.Encode()
	// server_DHInnerData_buf_sha1 := sha1.Sum(server_DHInnerData_buf)

	// 创建aes和iv key
	tmp_aes_key_and_iv := make([]byte, 64)
	sha1_a := sha1.Sum(append(authKeyMD.NewNonce, authKeyMD.ServerNonce...))
	sha1_b := sha1.Sum(append(authKeyMD.ServerNonce, authKeyMD.NewNonce...))
	sha1_c := sha1.Sum(append(authKeyMD.NewNonce, authKeyMD.NewNonce...))
	copy(tmp_aes_key_and_iv, sha1_a[:])
	copy(tmp_aes_key_and_iv[20:], sha1_b[:])
	copy(tmp_aes_key_and_iv[40:], sha1_c[:])
	copy(tmp_aes_key_and_iv[60:], authKeyMD.NewNonce[:4])

	tmpLen := 20 + len(server_DHInnerData_buf)
	if tmpLen%16 > 0 {
		tmpLen = (tmpLen/16 + 1) * 16
	} else {
		tmpLen = 20 + len(server_DHInnerData_buf)
	}

	tmp_encrypted_answer := make([]byte, tmpLen)
	sha1_tmp := sha1.Sum(server_DHInnerData_buf)
	copy(tmp_encrypted_answer, sha1_tmp[:])
	copy(tmp_encrypted_answer[20:], server_DHInnerData_buf)

	e := crypto.NewAES256IGECryptor(tmp_aes_key_and_iv[:32], tmp_aes_key_and_iv[32:64])
	tmp_encrypted_answer, _ = e.Encrypt(tmp_encrypted_answer)

	server_DHParamsOk := &mtproto.TLServer_DHParamsOk{Data2: &mtproto.Server_DH_Params_Data{
		Nonce:           authKeyMD.Nonce,
		ServerNonce:     authKeyMD.ServerNonce,
		EncryptedAnswer: string(tmp_encrypted_answer),
	}}

	glog.Infof("onReq_DHParams - metadata: {%v}, reply: %s", authKeyMD, logger.JsonDebugData(server_DHParamsOk))

	// s.authKeyMetadataToTrailer(ctx, authKeyMD)
	state.Ctx, _ = proto.Marshal(authKeyMD)
	return server_DHParamsOk.To_Server_DH_Params(), nil
}

// set_client_DH_params#f5045f1f nonce:int128 server_nonce:int128 encrypted_data:string = Set_client_DH_params_answer;
func (s *handshake) onSetClient_DHParams(state *zproto.HandshakeState, request *mtproto.TLSetClient_DHParams) (*mtproto.SetClient_DHParamsAnswer, error) {
	authKeyMD, err := toAuthKeyMetadata(state.Ctx)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	glog.Infof("onSetClient_DHParams - state: %d, res_state: %d, metadata: {%v}, request: %s",
		state.State, state.ResState, authKeyMD, logger.JsonDebugData(request))

	// TODO(@benqi): Impl SetClient_DHParams logic
	// 客户端传输数据解析
	// Nonce
	if !bytes.Equal(request.Nonce, authKeyMD.Nonce) {
		err := fmt.Errorf("process SetClient_DHParams - Wrong Nonce")
		glog.Error(err)
		return nil, err
	}

	// ServerNonce
	if !bytes.Equal(request.ServerNonce, authKeyMD.ServerNonce) {
		err := fmt.Errorf("process SetClient_DHParams - Wrong ServerNonce")
		glog.Error(err)
		return nil, err
	}

	bEncryptedData := []byte(request.EncryptedData)

	// 创建aes和iv key
	tmp_aes_key_and_iv := make([]byte, 64)
	sha1_a := sha1.Sum(append(authKeyMD.NewNonce, authKeyMD.ServerNonce...))
	sha1_b := sha1.Sum(append(authKeyMD.ServerNonce, authKeyMD.NewNonce...))
	sha1_c := sha1.Sum(append(authKeyMD.NewNonce, authKeyMD.NewNonce...))
	copy(tmp_aes_key_and_iv, sha1_a[:])
	copy(tmp_aes_key_and_iv[20:], sha1_b[:])
	copy(tmp_aes_key_and_iv[40:], sha1_c[:])
	copy(tmp_aes_key_and_iv[60:], authKeyMD.NewNonce[:4])

	d := crypto.NewAES256IGECryptor(tmp_aes_key_and_iv[:32], tmp_aes_key_and_iv[32:64])
	decryptedData, err := d.Decrypt(bEncryptedData)
	if err != nil {
		err := fmt.Errorf("process SetClient_DHParams - AES256IGECryptor descrypt error")
		glog.Error(err)
		return nil, err
	}

	// TODO(@benqi): 检查签名是否合法
	dbuf := mtproto.NewDecodeBuf(decryptedData[24:])
	client_DHInnerData := mtproto.NewTLClient_DHInnerData()
	// &TLClient_DHInnerData{}
	err = client_DHInnerData.Decode(dbuf)
	if err != nil {
		glog.Errorf("processSetClient_DHParams - TLClient_DHInnerData decode error: %s", err)
		return nil, err
	}

	glog.Info("processSetClient_DHParams - client_DHInnerData: ", client_DHInnerData.String())

	//
	if !bytes.Equal(client_DHInnerData.GetNonce(), authKeyMD.Nonce) {
		err := fmt.Errorf("process SetClient_DHParams - Wrong client_DHInnerData's Nonce")
		glog.Error(err)
		return nil, err
	}

	// ServerNonce
	if !bytes.Equal(client_DHInnerData.GetServerNonce(), authKeyMD.ServerNonce) {
		err := fmt.Errorf("process SetClient_DHParams - Wrong client_DHInnerData's ServerNonce")
		glog.Error(err)
		return nil, err
	}

	bigIntA := new(big.Int).SetBytes(authKeyMD.A)
	// bigIntP := new(big.Int).SetBytes(authKeyMD.P)

	// hash_key
	authKeyNum := new(big.Int)
	authKeyNum.Exp(new(big.Int).SetBytes([]byte(client_DHInnerData.GetGB())), bigIntA, s.bigIntDH2048P)

	authKey := make([]byte, 256)

	// TODO(@benqi): dhGenRetry and dhGenFail
	copy(authKey[256-len(authKeyNum.Bytes()):], authKeyNum.Bytes())

	authKeyAuxHash := make([]byte, len(authKeyMD.NewNonce))
	copy(authKeyAuxHash, authKeyMD.NewNonce)
	authKeyAuxHash = append(authKeyAuxHash, byte(0x01))
	sha1_d := sha1.Sum(authKey)
	authKeyAuxHash = append(authKeyAuxHash, sha1_d[:]...)
	sha1_e := sha1.Sum(authKeyAuxHash[:len(authKeyAuxHash)-12])
	authKeyAuxHash = append(authKeyAuxHash, sha1_e[:]...)

	// 至此key已经创建成功
	authKeyId := int64(binary.LittleEndian.Uint64(authKeyAuxHash[len(authKeyMD.NewNonce)+1+12 : len(authKeyMD.NewNonce)+1+12+8]))

	// TODO(@benqi): authKeyId生成后要检查在数据库里是否已经存在，有非常小的概率会碰撞
	// 如果碰撞让客户端重新再来一轮

	authKeyMD.AuthKeyId = authKeyId
	authKeyMD.AuthKey = authKey

	state.Ctx, _ = proto.Marshal(authKeyMD)

	if s.saveAuthKeyInfo(authKeyMD) {
		dhGenOk := &mtproto.TLDhGenOk{Data2: &mtproto.SetClient_DHParamsAnswer_Data{
			Nonce:         authKeyMD.Nonce,
			ServerNonce:   authKeyMD.ServerNonce,
			NewNonceHash1: calcNewNonceHash(authKeyMD.NewNonce, authKey, 0x01),
		}}

		glog.Infof("onSetClient_DHParams - metadata: {%v}, reply: %s", authKeyMD, logger.JsonDebugData(dhGenOk))
		return dhGenOk.To_SetClient_DHParamsAnswer(), nil
	} else {
		// TODO(@benqi): dhGenFail
		dhGenRetry := &mtproto.TLDhGenRetry{Data2: &mtproto.SetClient_DHParamsAnswer_Data{
			Nonce:         authKeyMD.Nonce,
			ServerNonce:   authKeyMD.ServerNonce,
			// NewNonceHash1: authKeyAuxHash[len(authKeyAuxHash)-16 : len(authKeyAuxHash)],
			NewNonceHash2: calcNewNonceHash(authKeyMD.NewNonce, authKey, 0x02),
		}}

		glog.Infof("onSetClient_DHParams - metadata: {%v}, reply: %s", authKeyMD, logger.JsonDebugData(dhGenRetry))
		return dhGenRetry.To_SetClient_DHParamsAnswer(), nil
	}
}

// msgs_ack#62d6b459 msg_ids:Vector<long> = MsgsAck;
func (s *handshake) onMsgsAck(state *zproto.HandshakeState, request *mtproto.TLMsgsAck) error {
	if state.ResState != zproto.RES_STATE_OK {
		return fmt.Errorf("invalid state: %v", state)
	}

	switch state.State {
	case zproto.STATE_pq_res:
		state.State = zproto.STATE_pq_ack
	case zproto.STATE_DH_params_res:
		state.State = zproto.STATE_DH_params_ack
	case zproto.STATE_dh_gen_res:
		state.State = zproto.STATE_dh_gen_ack
	default:
		return fmt.Errorf("invalid state: %v", state)
	}
	return nil
}

func (s *handshake) saveAuthKeyInfo(md *mtproto.AuthKeyMetadata) bool {
	var (
		salt = int64(0)
		serverSalt *mtproto.TLFutureSalt
		now = int32(time.Now().Unix())
	)

	for a := 7; a >= 0; a-- {
		salt <<= 8
		salt |= int64(md.NewNonce[a] ^ md.ServerNonce[a])
	}


	serverSalt = &mtproto.TLFutureSalt{Data2: &mtproto.FutureSalt_Data{
		ValidSince: now,
		ValidUntil: now + 30*60,
		Salt:       salt,
	}}

	authKeyInfo := &mtproto.TLAuthKeyInfo{Data2: &mtproto.AuthKeyInfo_Data{
		AuthKeyId:  md.AuthKeyId,
		AuthKey:    md.AuthKey,
		FutureSalt: serverSalt.To_FutureSalt(),
	}}

	request := &mtproto.TLSessionSetAuthKey{
		AuthKey: authKeyInfo.To_AuthKeyInfo(),
	}
	r, err := s.authSessionRpcClient.SessionSetAuthKey(context.Background(), request)

	if err != nil {
		glog.Error(err)
	}

	if !mtproto.FromBool(r) {
		glog.Error("saveAuthKeyInfo not successful - ", md.AuthKeyId)
	}

	return true
}

func calcNewNonceHash(newNonce, authKey []byte, b byte) []byte {
	authKeyAuxHash := make([]byte, len(newNonce))
	copy(authKeyAuxHash, newNonce)
	authKeyAuxHash = append(authKeyAuxHash, b)
	sha1D := sha1.Sum(authKey)
	authKeyAuxHash = append(authKeyAuxHash, sha1D[:]...)
	sha1E := sha1.Sum(authKeyAuxHash[:len(authKeyAuxHash)-12])
	authKeyAuxHash = append(authKeyAuxHash, sha1E[:]...)
	return authKeyAuxHash[len(authKeyAuxHash)-16 : len(authKeyAuxHash)]
}